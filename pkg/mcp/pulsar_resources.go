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
	"encoding/json"
	"fmt"
	"net/url"
	"slices"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/streamnative/pulsarctl/pkg/cmdutils"
	context2 "github.com/streamnative/streamnative-mcp-server/pkg/mcp/internal/context"
	pulsarsession "github.com/streamnative/streamnative-mcp-server/pkg/pulsar"
)

const (
	pulsarResourceContextURI            = "pulsar://context"
	pulsarResourceCatalogURI            = "pulsar://resources"
	pulsarNamespacesResourceTemplateURI = "pulsar://admin/v2/tenants/{tenant}/namespaces"
	pulsarTopicsResourceTemplateURI     = "pulsar://admin/v2/namespaces/{tenant}/{namespace}/topics"
	pulsarResourceJSONMIMEType          = "application/json"
)

type pulsarResourceKind string

const (
	pulsarResourceKindContext    pulsarResourceKind = "context"
	pulsarResourceKindCatalog    pulsarResourceKind = "catalog"
	pulsarResourceKindNamespaces pulsarResourceKind = "namespaces"
	pulsarResourceKindTopics     pulsarResourceKind = "topics"
)

type pulsarResourceURI struct {
	kind      pulsarResourceKind
	tenant    string
	namespace string
}

type pulsarResourceCatalog struct {
	Version   int                     `json:"version"`
	Scheme    string                  `json:"scheme"`
	Resources []pulsarCatalogResource `json:"resources"`
	Templates []pulsarCatalogTemplate `json:"templates"`
	Notes     []string                `json:"notes,omitempty"`
}

type pulsarCatalogResource struct {
	URI         string `json:"uri"`
	Name        string `json:"name"`
	Description string `json:"description"`
	MIMEType    string `json:"mimeType"`
}

type pulsarCatalogTemplate struct {
	URITemplate string `json:"uriTemplate"`
	Name        string `json:"name"`
	Description string `json:"description"`
	MIMEType    string `json:"mimeType"`
}

type pulsarContextResource struct {
	Kind           string                      `json:"kind"`
	URI            string                      `json:"uri"`
	ServiceURL     string                      `json:"serviceUrl,omitempty"`
	WebServiceURL  string                      `json:"webServiceUrl,omitempty"`
	Authentication pulsarAuthenticationSummary `json:"authentication"`
	TLS            pulsarTLSSummary            `json:"tls"`
}

type pulsarAuthenticationSummary struct {
	Configured bool   `json:"configured"`
	Method     string `json:"method"`
	Plugin     string `json:"plugin,omitempty"`
}

type pulsarTLSSummary struct {
	AllowInsecureConnection    bool `json:"allowInsecureConnection"`
	EnableHostnameVerification bool `json:"enableHostnameVerification"`
	TrustCertsFileConfigured   bool `json:"trustCertsFileConfigured"`
	ClientCertFileConfigured   bool `json:"clientCertFileConfigured"`
	ClientKeyFileConfigured    bool `json:"clientKeyFileConfigured"`
}

type pulsarNamespaceCollectionResource struct {
	Kind       string   `json:"kind"`
	URI        string   `json:"uri"`
	Tenant     string   `json:"tenant"`
	Namespaces []string `json:"namespaces"`
	Count      int      `json:"count"`
}

type pulsarTopicCollectionResource struct {
	Kind      string   `json:"kind"`
	URI       string   `json:"uri"`
	Tenant    string   `json:"tenant"`
	Namespace string   `json:"namespace"`
	Topics    []string `json:"topics"`
	Count     int      `json:"count"`
}

// PulsarAddResources registers the read-only Pulsar MCP resource surface.
func PulsarAddResources(s *server.MCPServer, features []string) {
	if !pulsarResourcesEnabled(features) {
		return
	}

	s.AddResources(
		server.ServerResource{
			Resource: mcp.NewResource(pulsarResourceContextURI, "Pulsar Context",
				mcp.WithResourceDescription("Current Pulsar session connection metadata with authentication material redacted."),
				mcp.WithMIMEType(pulsarResourceJSONMIMEType),
			),
			Handler: handlePulsarContextResource,
		},
		server.ServerResource{
			Resource: mcp.NewResource(pulsarResourceCatalogURI, "Pulsar Resource Catalog",
				mcp.WithResourceDescription("Stable catalog of Pulsar MCP resource URIs and URI templates."),
				mcp.WithMIMEType(pulsarResourceJSONMIMEType),
			),
			Handler: handlePulsarCatalogResource,
		},
	)

	s.AddResourceTemplates(
		server.ServerResourceTemplate{
			Template: mcp.NewResourceTemplate(pulsarNamespacesResourceTemplateURI, "Pulsar Namespaces by Tenant",
				mcp.WithTemplateDescription("List namespaces for a Pulsar tenant."),
				mcp.WithTemplateMIMEType(pulsarResourceJSONMIMEType),
			),
			Handler: handlePulsarNamespacesResource,
		},
		server.ServerResourceTemplate{
			Template: mcp.NewResourceTemplate(pulsarTopicsResourceTemplateURI, "Pulsar Topics by Namespace",
				mcp.WithTemplateDescription("List topics for a Pulsar namespace."),
				mcp.WithTemplateMIMEType(pulsarResourceJSONMIMEType),
			),
			Handler: handlePulsarTopicsResource,
		},
	)
}

func pulsarResourcesEnabled(features []string) bool {
	requiredFeatures := []Feature{
		FeatureAll,
		FeatureAllPulsar,
		FeaturePulsarAdmin,
		FeaturePulsarAdminNamespaces,
		FeaturePulsarAdminTopics,
	}
	for _, feature := range requiredFeatures {
		if slices.Contains(features, string(feature)) {
			return true
		}
	}
	return false
}

func handlePulsarContextResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	parsed, err := parsePulsarResourceURI(request.Params.URI)
	if err != nil {
		return nil, err
	}
	if parsed.kind != pulsarResourceKindContext {
		return nil, fmt.Errorf("unsupported Pulsar context resource URI %q", request.Params.URI)
	}
	session, err := requirePulsarResourceSession(ctx)
	if err != nil {
		return nil, err
	}

	resource, err := buildPulsarContextResource(request.Params.URI, session)
	if err != nil {
		return nil, err
	}
	return newPulsarJSONResourceContents(request.Params.URI, resource)
}

func handlePulsarCatalogResource(_ context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	parsed, err := parsePulsarResourceURI(request.Params.URI)
	if err != nil {
		return nil, err
	}
	if parsed.kind != pulsarResourceKindCatalog {
		return nil, fmt.Errorf("unsupported Pulsar resource catalog URI %q", request.Params.URI)
	}
	return newPulsarJSONResourceContents(request.Params.URI, buildPulsarResourceCatalog())
}

func handlePulsarNamespacesResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	parsed, err := parsePulsarResourceURI(request.Params.URI)
	if err != nil {
		return nil, err
	}
	if parsed.kind != pulsarResourceKindNamespaces {
		return nil, fmt.Errorf("unsupported Pulsar namespaces resource URI %q", request.Params.URI)
	}

	session, err := requirePulsarResourceSession(ctx)
	if err != nil {
		return nil, err
	}
	adminClient, err := getPulsarResourceAdminClient(session)
	if err != nil {
		return nil, err
	}

	namespaces, err := adminClient.Namespaces().GetNamespaces(parsed.tenant)
	if err != nil {
		return nil, fmt.Errorf("failed to list namespaces for tenant %q: %w", parsed.tenant, err)
	}

	return newPulsarJSONResourceContents(request.Params.URI, pulsarNamespaceCollectionResource{
		Kind:       string(parsed.kind),
		URI:        request.Params.URI,
		Tenant:     parsed.tenant,
		Namespaces: namespaces,
		Count:      len(namespaces),
	})
}

func handlePulsarTopicsResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	parsed, err := parsePulsarResourceURI(request.Params.URI)
	if err != nil {
		return nil, err
	}
	if parsed.kind != pulsarResourceKindTopics {
		return nil, fmt.Errorf("unsupported Pulsar topics resource URI %q", request.Params.URI)
	}

	session, err := requirePulsarResourceSession(ctx)
	if err != nil {
		return nil, err
	}
	adminClient, err := getPulsarResourceAdminClient(session)
	if err != nil {
		return nil, err
	}

	namespaceName := parsed.tenant + "/" + parsed.namespace
	topics, err := adminClient.Namespaces().GetTopics(namespaceName)
	if err != nil {
		return nil, fmt.Errorf("failed to list topics for namespace %q: %w", namespaceName, err)
	}

	return newPulsarJSONResourceContents(request.Params.URI, pulsarTopicCollectionResource{
		Kind:      string(parsed.kind),
		URI:       request.Params.URI,
		Tenant:    parsed.tenant,
		Namespace: parsed.namespace,
		Topics:    topics,
		Count:     len(topics),
	})
}

func requirePulsarResourceSession(ctx context.Context) (*pulsarsession.Session, error) {
	session := context2.GetPulsarSession(ctx)
	if session == nil {
		return nil, fmt.Errorf("Pulsar session not found in context")
	}
	return session, nil
}

func getPulsarResourceAdminClient(session *pulsarsession.Session) (cmdutils.Client, error) {
	if _, err := session.GetPulsarCtlConfig(); err != nil {
		return nil, fmt.Errorf("failed to get Pulsar admin configuration: %w", err)
	}
	adminClient, err := session.GetAdminClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get Pulsar admin client: %w", err)
	}
	if adminClient == nil {
		return nil, fmt.Errorf("Pulsar admin client not found in session")
	}
	return adminClient, nil
}

func buildPulsarContextResource(uri string, session *pulsarsession.Session) (pulsarContextResource, error) {
	ctx := session.Ctx
	if ctx.WebServiceURL == "" {
		if cfg, err := session.GetPulsarCtlConfig(); err == nil && cfg != nil {
			ctx.WebServiceURL = cfg.WebServiceURL
			ctx.TLSAllowInsecureConnection = cfg.TLSAllowInsecureConnection
			ctx.TLSEnableHostnameVerification = cfg.TLSEnableHostnameVerification
			ctx.TLSTrustCertsFilePath = cfg.TLSTrustCertsFilePath
			ctx.TLSCertFile = cfg.TLSCertFile
			ctx.TLSKeyFile = cfg.TLSKeyFile
		}
	}
	if ctx.ServiceURL == "" && ctx.WebServiceURL == "" {
		return pulsarContextResource{}, fmt.Errorf("Pulsar session is not configured")
	}

	return pulsarContextResource{
		Kind:           string(pulsarResourceKindContext),
		URI:            uri,
		ServiceURL:     ctx.ServiceURL,
		WebServiceURL:  ctx.WebServiceURL,
		Authentication: summarizePulsarAuthentication(ctx),
		TLS: pulsarTLSSummary{
			AllowInsecureConnection:    ctx.TLSAllowInsecureConnection,
			EnableHostnameVerification: ctx.TLSEnableHostnameVerification,
			TrustCertsFileConfigured:   ctx.TLSTrustCertsFilePath != "",
			ClientCertFileConfigured:   ctx.TLSCertFile != "",
			ClientKeyFileConfigured:    ctx.TLSKeyFile != "",
		},
	}, nil
}

func summarizePulsarAuthentication(ctx pulsarsession.PulsarContext) pulsarAuthenticationSummary {
	switch {
	case ctx.Token != "":
		return pulsarAuthenticationSummary{
			Configured: true,
			Method:     "token",
		}
	case ctx.AuthPlugin != "" || ctx.AuthParams != "":
		return pulsarAuthenticationSummary{
			Configured: true,
			Method:     "authPlugin",
			Plugin:     ctx.AuthPlugin,
		}
	default:
		return pulsarAuthenticationSummary{
			Configured: false,
			Method:     "none",
		}
	}
}

func buildPulsarResourceCatalog() pulsarResourceCatalog {
	return pulsarResourceCatalog{
		Version: 1,
		Scheme:  "pulsar",
		Resources: []pulsarCatalogResource{
			{
				URI:         pulsarResourceContextURI,
				Name:        "Pulsar Context",
				Description: "Current Pulsar session connection metadata with authentication material redacted.",
				MIMEType:    pulsarResourceJSONMIMEType,
			},
			{
				URI:         pulsarResourceCatalogURI,
				Name:        "Pulsar Resource Catalog",
				Description: "Stable catalog of Pulsar MCP resource URIs and URI templates.",
				MIMEType:    pulsarResourceJSONMIMEType,
			},
		},
		Templates: []pulsarCatalogTemplate{
			{
				URITemplate: pulsarNamespacesResourceTemplateURI,
				Name:        "Pulsar Namespaces by Tenant",
				Description: "List namespaces for a Pulsar tenant.",
				MIMEType:    pulsarResourceJSONMIMEType,
			},
			{
				URITemplate: pulsarTopicsResourceTemplateURI,
				Name:        "Pulsar Topics by Namespace",
				Description: "List topics for a Pulsar namespace.",
				MIMEType:    pulsarResourceJSONMIMEType,
			},
		},
		Notes: []string{
			"Resource handlers are read-only and do not return tokens, auth params, key files, TLS private keys, or secret values.",
		},
	}
}

func parsePulsarResourceURI(rawURI string) (pulsarResourceURI, error) {
	parsed, err := url.Parse(rawURI)
	if err != nil {
		return pulsarResourceURI{}, fmt.Errorf("malformed Pulsar resource URI %q: %w", rawURI, err)
	}
	if parsed.Scheme != "pulsar" {
		return pulsarResourceURI{}, fmt.Errorf("unsupported Pulsar resource URI scheme %q", parsed.Scheme)
	}
	if parsed.User != nil {
		return pulsarResourceURI{}, fmt.Errorf("Pulsar resource URI %q must not include user info", rawURI)
	}
	if parsed.RawQuery != "" || parsed.Fragment != "" {
		return pulsarResourceURI{}, fmt.Errorf("Pulsar resource URI %q must not include query or fragment", rawURI)
	}

	switch parsed.Host {
	case "context":
		if parsed.Path != "" {
			return pulsarResourceURI{}, fmt.Errorf("unsupported Pulsar context resource URI %q", rawURI)
		}
		return pulsarResourceURI{kind: pulsarResourceKindContext}, nil
	case "resources":
		if parsed.Path != "" {
			return pulsarResourceURI{}, fmt.Errorf("unsupported Pulsar resource catalog URI %q", rawURI)
		}
		return pulsarResourceURI{kind: pulsarResourceKindCatalog}, nil
	case "admin":
		return parsePulsarAdminResourceURI(rawURI, parsed.Path)
	default:
		return pulsarResourceURI{}, fmt.Errorf("unsupported Pulsar resource URI host %q", parsed.Host)
	}
}

func parsePulsarAdminResourceURI(rawURI, path string) (pulsarResourceURI, error) {
	parts := splitPulsarResourcePath(path)
	switch {
	case len(parts) == 4 && parts[0] == "v2" && parts[1] == "tenants" && parts[3] == "namespaces":
		tenant := parts[2]
		if err := validatePulsarResourcePathSegment("tenant", tenant); err != nil {
			return pulsarResourceURI{}, err
		}
		return pulsarResourceURI{
			kind:   pulsarResourceKindNamespaces,
			tenant: tenant,
		}, nil
	case len(parts) == 5 && parts[0] == "v2" && parts[1] == "namespaces" && parts[4] == "topics":
		tenant := parts[2]
		namespace := parts[3]
		if err := validatePulsarResourcePathSegment("tenant", tenant); err != nil {
			return pulsarResourceURI{}, err
		}
		if err := validatePulsarResourcePathSegment("namespace", namespace); err != nil {
			return pulsarResourceURI{}, err
		}
		return pulsarResourceURI{
			kind:      pulsarResourceKindTopics,
			tenant:    tenant,
			namespace: namespace,
		}, nil
	default:
		return pulsarResourceURI{}, fmt.Errorf("unsupported Pulsar admin resource URI %q", rawURI)
	}
}

func splitPulsarResourcePath(path string) []string {
	trimmed := strings.Trim(path, "/")
	if trimmed == "" {
		return nil
	}
	return strings.Split(trimmed, "/")
}

func validatePulsarResourcePathSegment(name, value string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("Pulsar resource URI %s segment must not be empty", name)
	}
	if strings.Contains(value, "/") {
		return fmt.Errorf("Pulsar resource URI %s segment must not contain path separators", name)
	}
	return nil
}

func newPulsarJSONResourceContents(uri string, value any) ([]mcp.ResourceContents, error) {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Pulsar resource %q: %w", uri, err)
	}
	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      uri,
			MIMEType: pulsarResourceJSONMIMEType,
			Text:     string(data),
		},
	}, nil
}
