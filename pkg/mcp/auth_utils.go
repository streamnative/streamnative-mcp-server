package mcp

import (
	"fmt"

	"github.com/streamnative/streamnative-mcp-server/pkg/auth"
	sncloud "github.com/streamnative/streamnative-mcp-server/sdk/sdk-apiserve"
)

// getFlow creates the appropriate OAuth2 flow based on the grant type
func getFlow(issuer *auth.Issuer, grant *auth.AuthorizationGrant) (auth.Flow, error) {
	switch grant.Type {
	case auth.GrantTypeClientCredentials:
		// Use client credentials flow for service accounts
		if grant.ClientCredentials == nil {
			return nil, fmt.Errorf("client credentials grant missing required credentials")
		}
		return auth.NewDefaultClientCredentialsFlowWithKeyFileStruct(*issuer, grant.ClientCredentials)
	default:
		return nil, fmt.Errorf("unsupported grant type: %s", grant.Type)
	}
}

// getBasePath constructs the base path for Pulsar admin API
func getBasePath(proxyLocation, org, clusterID string) string {
	// Ensure proxyLocation doesn't end with a slash to prevent double slashes
	if proxyLocation != "" && proxyLocation[len(proxyLocation)-1] == '/' {
		proxyLocation = proxyLocation[:len(proxyLocation)-1]
	}

	return fmt.Sprintf("%s/pulsar-admin/%s/pulsarcluster-%s", proxyLocation, org, clusterID)
}

func getIssuer(instance *sncloud.ComGithubStreamnativeCloudApiServerPkgApisCloudV1alpha1PulsarInstance, configIssuer auth.Issuer) (*auth.Issuer, error) {
	if instance.Status == nil {
		return nil, fmt.Errorf("PulsarInstance '%s' has no auth configuration", instance.Metadata.Name)
	}

	if instance.Status.Auth.Type != "oauth2" && instance.Status.Auth.Type != "apikey" {
		return nil, fmt.Errorf("PulsarInstance '%s' has unsupported auth type: %s",
			instance.Metadata.Name, instance.Status.Auth.Type)
	}

	if instance.Status.Auth.Oauth2.Audience == "" || instance.Status.Auth.Oauth2.IssuerURL == "" {
		return nil, fmt.Errorf("PulsarInstance '%s' has no OAuth2 configuration", instance.Metadata.Name)
	}

	// Construct issuer information from instance and config
	return &auth.Issuer{
		IssuerEndpoint: instance.Status.Auth.Oauth2.IssuerURL,
		ClientID:       configIssuer.ClientID,
		Audience:       instance.Status.Auth.Oauth2.Audience,
	}, nil
}
