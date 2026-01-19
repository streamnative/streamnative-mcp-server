// Copyright 2025 StreamNative
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
	"context"
	"fmt"
	"log/slog"
	"runtime/debug"
	"time"

	sdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/sirupsen/logrus"
	"github.com/streamnative/streamnative-mcp-server/pkg/config"
	"github.com/streamnative/streamnative-mcp-server/pkg/kafka"
	"github.com/streamnative/streamnative-mcp-server/pkg/pulsar"
)

// Server wraps MCP go-sdk server state and StreamNative sessions.
type Server struct {
	MCPServer      *sdk.Server
	KafkaSession   *kafka.Session
	PulsarSession  *pulsar.Session
	SNCloudSession *config.Session
	logger         *logrus.Logger
}

// ServerOption mutates MCP go-sdk server options.
type ServerOption func(*sdk.ServerOptions)

// NewServer creates a new MCP go-sdk server with StreamNative integrations.
func NewServer(name, version string, logger *logrus.Logger, opts ...ServerOption) *Server {
	serverOptions := defaultServerOptions(logger)
	for _, opt := range opts {
		if opt != nil {
			opt(serverOptions)
		}
	}

	impl := &sdk.Implementation{Name: name, Version: version}
	mcpServer := sdk.NewServer(impl, serverOptions)
	addDefaultMiddleware(mcpServer, logger)

	return &Server{
		MCPServer:      mcpServer,
		logger:         logger,
		SNCloudSession: &config.Session{},
		KafkaSession:   &kafka.Session{},
		PulsarSession:  &pulsar.Session{},
	}
}

// WithInstructions sets server instructions returned in initialize responses.
func WithInstructions(instructions string) ServerOption {
	return func(opts *sdk.ServerOptions) {
		opts.Instructions = instructions
	}
}

// WithCapabilities overrides the default server capability configuration.
func WithCapabilities(capabilities *sdk.ServerCapabilities) ServerOption {
	return func(opts *sdk.ServerOptions) {
		opts.Capabilities = capabilities
	}
}

// WithLogger overrides the default slog logger for the go-sdk server.
func WithLogger(logger *slog.Logger) ServerOption {
	return func(opts *sdk.ServerOptions) {
		opts.Logger = logger
	}
}

func defaultServerOptions(logger *logrus.Logger) *sdk.ServerOptions {
	opts := &sdk.ServerOptions{
		Capabilities: &sdk.ServerCapabilities{
			Logging: &sdk.LoggingCapabilities{},
			Resources: &sdk.ResourceCapabilities{
				Subscribe:   true,
				ListChanged: true,
			},
		},
	}

	if logger != nil {
		opts.Logger = slog.New(slog.NewTextHandler(logger.Writer(), &slog.HandlerOptions{
			Level: slogLevelFromLogrus(logger.Level),
		}))
	}

	return opts
}

func slogLevelFromLogrus(level logrus.Level) slog.Level {
	switch level {
	case logrus.TraceLevel, logrus.DebugLevel:
		return slog.LevelDebug
	case logrus.WarnLevel:
		return slog.LevelWarn
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func addDefaultMiddleware(server *sdk.Server, logger *logrus.Logger) {
	if server == nil {
		return
	}

	middlewares := []sdk.Middleware{
		recoveryMiddleware(logger),
	}
	if logger != nil {
		middlewares = append([]sdk.Middleware{loggingMiddleware(logger)}, middlewares...)
	}
	server.AddReceivingMiddleware(middlewares...)
}

func loggingMiddleware(logger *logrus.Logger) sdk.Middleware {
	return func(next sdk.MethodHandler) sdk.MethodHandler {
		return func(ctx context.Context, method string, req sdk.Request) (sdk.Result, error) {
			start := time.Now()
			sessionID := ""
			if req != nil {
				if session := req.GetSession(); session != nil {
					sessionID = session.ID()
				}
			}

			entry := logger.WithFields(logrus.Fields{
				"method":     method,
				"session_id": sessionID,
			})

			if callReq, ok := req.(*sdk.CallToolRequest); ok && callReq != nil && callReq.Params != nil {
				entry = entry.WithField("tool", callReq.Params.Name)
				entry.Debug("MCP tool call started")
			} else {
				entry.Debug("MCP request started")
			}

			result, err := next(ctx, method, req)
			duration := time.Since(start)

			if err != nil {
				entry.WithFields(logrus.Fields{
					"duration_ms": duration.Milliseconds(),
				}).WithError(err).Error("MCP request failed")
				return result, err
			}

			entry.WithFields(logrus.Fields{
				"duration_ms": duration.Milliseconds(),
			}).Debug("MCP request completed")
			return result, nil
		}
	}
}

func recoveryMiddleware(logger *logrus.Logger) sdk.Middleware {
	return func(next sdk.MethodHandler) sdk.MethodHandler {
		return func(ctx context.Context, method string, req sdk.Request) (result sdk.Result, err error) {
			defer func() {
				if recovered := recover(); recovered != nil {
					toolName := ""
					if callReq, ok := req.(*sdk.CallToolRequest); ok && callReq != nil && callReq.Params != nil {
						toolName = callReq.Params.Name
					}

					if logger != nil {
						sessionID := ""
						if req != nil {
							if session := req.GetSession(); session != nil {
								sessionID = session.ID()
							}
						}
						fields := logrus.Fields{
							"method":     method,
							"session_id": sessionID,
							"panic":      recovered,
						}
						if toolName != "" {
							fields["tool"] = toolName
						}
						logger.WithFields(fields).Error("MCP request panic recovered")
						logger.WithField("stack", string(debug.Stack())).Debug("MCP panic stack")
					}

					if toolName != "" {
						err = fmt.Errorf("panic recovered in %s tool handler: %v", toolName, recovered)
					} else {
						err = fmt.Errorf("panic recovered in %s request: %v", method, recovered)
					}
					result = nil
				}
			}()

			return next(ctx, method, req)
		}
	}
}
