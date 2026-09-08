// Copyright 2026 StreamNative
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mcp

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path"
	"strings"
	"sync"
	"syscall"
	"time"

	protocol "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/streamnative/streamnative-mcp-server/pkg/common"
	"github.com/streamnative/streamnative-mcp-server/pkg/config"
	app "github.com/streamnative/streamnative-mcp-server/pkg/mcp"
	"github.com/streamnative/streamnative-mcp-server/pkg/mcp/session"
	"github.com/streamnative/streamnative-mcp-server/pkg/pulsar"
)

// httpSessionResolver leases a backend session until release is called after
// ServeHTTP returns. It does not assert that a token is authorized by Pulsar.
type httpSessionResolver func(context.Context, string) (*pulsar.Session, func(), error)

// NewCmdMcpHTTPServer creates the modern, fixed-backend Streamable HTTP command.
// Transport flags deliberately have independent storage from SSE/root flags.
func NewCmdMcpHTTPServer(shared *ServerOptions) *cobra.Command {
	var addr, endpoint string
	var origins []string
	cmd := &cobra.Command{
		Use: "http", Short: "Start modern Streamable HTTP server (external backends only)",
		Args: cobra.NoArgs,
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			// Check both flags and env before Complete touches keyring or credentials.
			if shared.KeyFile != "" || viper.GetString("key-file") != "" ||
				(shared.UseExternalKafka || viper.GetBool("use-external-kafka")) == (shared.UseExternalPulsar || viper.GetBool("use-external-pulsar")) {
				return fmt.Errorf("http transport requires exactly one external Kafka/Pulsar mode; Cloud is unsupported")
			}
			externalPulsar := shared.UseExternalPulsar || viper.GetBool("use-external-pulsar")
			if shared.MultiSessionPulsar && !externalPulsar {
				return fmt.Errorf("http transport multi-session requires external Pulsar")
			}
			host, _, err := net.SplitHostPort(addr)
			if err != nil {
				return fmt.Errorf("http transport invalid listen address: %w", err)
			}
			ip := net.ParseIP(host)
			if (ip == nil || !ip.IsLoopback()) && !shared.MultiSessionPulsar {
				return fmt.Errorf("http transport non-loopback binding requires multi-session Pulsar")
			}
			if !validHTTPEndpoint(endpoint) {
				return fmt.Errorf("http transport invalid http-path")
			}
			for _, origin := range origins {
				if !validHTTPOrigin(origin) {
					return fmt.Errorf("http transport invalid allowed Origin %q", origin)
				}
			}
			if cmd.Flags().Changed("session-cache-size") || cmd.Flags().Changed("session-ttl-minutes") {
				return fmt.Errorf("http transport uses request-scoped sessions, not a shared session cache")
			}
			return shared.Complete()
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			opts := *shared
			opts.HTTPAddr, opts.HTTPPath = addr, endpoint
			return runHTTPServer(cmd.Context(), &opts, origins)
		},
	}
	cmd.Flags().StringVar(&addr, "http-addr", "127.0.0.1:9090", "HTTP listen address (non-loopback requires multi-session Pulsar)")
	cmd.Flags().StringVar(&endpoint, "http-path", "/mcp", "Streamable HTTP endpoint")
	cmd.Flags().StringSliceVar(&origins, "http-allowed-origins", nil, "Exact allowed browser Origins; present Origins are denied by default")
	return cmd
}

func runHTTPServer(parent context.Context, opts *ServerOptions, origins []string) error {
	ctx, stop := signal.NotifyContext(parent, os.Interrupt, syscall.SIGTERM)
	defer stop()
	logger, err := initLogger(opts.LogFile)
	if err != nil {
		return err
	}
	if file, ok := logger.Out.(*os.File); ok && opts.LogFile != "" {
		defer func() { _ = file.Close() }()
	}
	s, err := newMcpServer(ctx, opts, logger)
	if err != nil {
		return fmt.Errorf("create HTTP MCP server: %w", err)
	}
	defer func() {
		if s.KafkaSession != nil {
			s.KafkaSession.ResetKafkaContext()
		}
		if s.PulsarSession != nil {
			s.PulsarSession.ResetPulsarContext()
		}
	}()
	// No Functions-as-tools: the catalog must not depend on protocol sessions.
	var resolve httpSessionResolver
	if opts.MultiSessionPulsar {
		p := opts.Pulsar
		base := pulsar.PulsarContext{
			ServiceURL: p.ServiceURL, WebServiceURL: p.WebServiceURL,
			TLSAllowInsecureConnection:    p.TLSAllowInsecureConnection,
			TLSEnableHostnameVerification: p.TLSEnableHostnameVerification,
			TLSTrustCertsFilePath:         p.TLSTrustCertsFilePath,
		}
		// The existing manager has no leases: LRU/TTL eviction closes in-flight
		// clients. Give each request a private manager and stop it only on release.
		// No fixed token, auth plugin or client certificate is inherited.
		resolve = func(ctx context.Context, token string) (*pulsar.Session, func(), error) {
			manager := session.NewPulsarSessionManager(&session.PulsarSessionManagerConfig{
				MaxSessions: 1, SessionTTL: time.Duration(1<<63 - 1), CleanupInterval: time.Duration(1<<63 - 1), BaseContext: base,
			}, nil, logger)
			ps, err := manager.GetOrCreateSession(ctx, token)
			return ps, sync.OnceFunc(manager.Stop), err
		}
	}
	handler, err := newHTTPHandler(s, opts.Options, opts.HTTPPath, origins, resolve)
	if err != nil {
		return err
	}
	// Track complete handler lifetimes even after forced network shutdown.
	// Backend clients must never close underneath active tool handlers.
	var active sync.WaitGroup
	httpServer := &http.Server{
		Addr: opts.HTTPAddr, ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout: 30 * time.Second, // Bound slow request bodies without timing out streaming writes.
		BaseContext: func(net.Listener) context.Context { return ctx },
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			active.Add(1)
			defer active.Done()
			handler.ServeHTTP(w, r)
		}),
	}
	listener, err := net.Listen("tcp", opts.HTTPAddr)
	if err != nil {
		return fmt.Errorf("listen HTTP: %w", err)
	}
	errCh := make(chan error, 1)
	go func() { errCh <- httpServer.Serve(listener) }()
	fmt.Fprintf(os.Stderr, "Streamable HTTP MCP server listening on http://%s%s\n", listener.Addr(), opts.HTTPPath)
	select {
	case <-ctx.Done():
	case err = <-errCh:
		stop()
	}
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	shutdownErr := httpServer.Shutdown(shutdownCtx)
	if shutdownErr != nil {
		_ = httpServer.Close()
	}
	active.Wait()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("serve HTTP: %w", err)
	}
	if shutdownErr != nil {
		return fmt.Errorf("shutdown HTTP: %w", shutdownErr)
	}
	return nil
}

func newHTTPHandler(s *app.Server, options *config.Options, endpoint string, origins []string, resolve httpSessionResolver) (http.Handler, error) {
	if !validHTTPEndpoint(endpoint) {
		return nil, fmt.Errorf("http-path must be a clean absolute URL path")
	}
	allowed := make(map[string]bool, len(origins))
	for _, origin := range origins {
		if !validHTTPOrigin(origin) {
			return nil, fmt.Errorf("invalid allowed Origin %q", origin)
		}
		allowed[origin] = true
	}
	transport := server.NewStreamableHTTPServer(s.MCPServer,
		server.WithEndpointPath(endpoint),
		server.WithStreamableHTTPProtocolVersions(protocol.ProtocolVersion20260728))
	mux := http.NewServeMux()
	mux.Handle(endpoint, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != endpoint {
			http.NotFound(w, r)
			return
		}
		// Deny every present Origin unless explicitly allowed. Host is not a trust anchor.
		if values, present := r.Header["Origin"]; present && (len(values) != 1 || !validHTTPOrigin(values[0]) || !allowed[values[0]]) {
			http.Error(w, "Origin not allowed", http.StatusForbidden)
			return
		}
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			http.Error(w, "only POST is supported", http.StatusMethodNotAllowed)
			return
		}
		// Bound memory before the SDK's unbounded io.ReadAll. Streaming responses
		// remain unrestricted; this is only the incoming JSON-RPC envelope limit.
		body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 4*1024*1024))
		if err != nil {
			var tooLarge *http.MaxBytesError
			if errors.As(err, &tooLarge) {
				http.Error(w, "request body too large", http.StatusRequestEntityTooLarge)
			} else {
				http.Error(w, "cannot read request body", http.StatusBadRequest)
			}
			return
		}
		r.Body = io.NopCloser(bytes.NewReader(body))
		if !validateHTTPEnvelope(w, r, body) {
			return
		}
		ctx := context.WithValue(r.Context(), common.OptionsKey, options)
		ctx = app.WithKafkaSession(ctx, s.KafkaSession)
		if resolve != nil {
			values := r.Header.Values("Authorization")
			var fields []string
			if len(values) == 1 {
				fields = strings.Fields(values[0])
			}
			if len(fields) != 2 || !strings.EqualFold(fields[0], "Bearer") {
				w.Header().Set("WWW-Authenticate", "Bearer")
				http.Error(w, "Pulsar bearer credentials required", http.StatusUnauthorized)
				return
			}
			ps, release, err := resolve(ctx, fields[1])
			if release != nil {
				defer release()
			}
			if err != nil || ps == nil {
				http.Error(w, "backend session unavailable", http.StatusServiceUnavailable)
				return
			}
			ctx = app.WithPulsarSession(ctx, ps)
		} else {
			ctx = app.WithPulsarSession(ctx, s.PulsarSession)
		}
		transport.ServeHTTP(w, r.WithContext(ctx))
	}))
	mux.HandleFunc(joinHTTPPath(endpoint, "healthz"), healthHandler("ok"))
	mux.HandleFunc(joinHTTPPath(endpoint, "readyz"), healthHandler("ready"))
	return mux, nil
}

func validHTTPEndpoint(endpoint string) bool {
	return endpoint != "" && strings.HasPrefix(endpoint, "/") && path.Clean(endpoint) == endpoint && !strings.ContainsAny(endpoint, "{}?#% \t\r\n")
}

// Preflight only the stateless envelope before acquiring backend credentials.
// The SDK remains responsible for dispatch, schemas, capabilities and results.
func validateHTTPEnvelope(w http.ResponseWriter, r *http.Request, body []byte) bool {
	var request struct {
		ID     json.RawMessage `json:"id"`
		Method protocol.MCPMethod
		Params json.RawMessage
	}
	fail := func(code int, message string) bool {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]any{"jsonrpc": "2.0", "id": request.ID, "error": map[string]any{"code": code, "message": message}})
		return false
	}
	if json.Unmarshal(body, &request) != nil {
		return fail(protocol.PARSE_ERROR, "invalid JSON-RPC request")
	}
	var params struct {
		Meta *protocol.Meta `json:"_meta"`
	}
	if json.Unmarshal(request.Params, &params) != nil || params.Meta == nil {
		return fail(protocol.INVALID_PARAMS, "modern per-request metadata required")
	}
	version := params.Meta.ProtocolVersion()
	if r.Header.Get(protocol.HeaderProtocolVersion) != version {
		return fail(protocol.HEADER_MISMATCH, "protocol header must match request metadata")
	}
	if version != protocol.ProtocolVersion20260728 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		details := (protocol.UnsupportedProtocolVersionError{Version: version, Supported: []string{protocol.ProtocolVersion20260728}}).JSONRPCError().Error
		_ = json.NewEncoder(w).Encode(map[string]any{"jsonrpc": "2.0", "id": request.ID, "error": details})
		return false
	}
	if params.Meta.GetMetaField(protocol.MetaKeyClientCapabilities) == nil {
		return fail(protocol.MISSING_REQUIRED_CLIENT_CAPABILITY, "clientCapabilities are required on every request")
	}
	if params.Meta.ClientCapabilities() == nil {
		return fail(protocol.HEADER_MISMATCH, "invalid clientCapabilities metadata")
	}
	if err := protocol.ValidateStandardHeaders(r.Header.Get, version, request.Method, request.Params); err != nil {
		return fail(protocol.HEADER_MISMATCH, "MCP headers do not match request metadata")
	}
	if r.Header.Get(protocol.HeaderLastEventID) != "" {
		return fail(protocol.HEADER_MISMATCH, "stream resumption is unsupported")
	}
	return true
}

func validHTTPOrigin(origin string) bool {
	u, err := url.Parse(origin)
	return err == nil && (u.Scheme == "http" || u.Scheme == "https") && u.Host != "" &&
		u.User == nil && u.Path == "" && u.RawQuery == "" && !u.ForceQuery && u.Fragment == "" &&
		!strings.ContainsAny(origin, "# \t\r\n,") && u.Hostname() != ""
}
