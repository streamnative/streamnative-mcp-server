package pulsar

import (
	"fmt"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	pulsaradminconfig "github.com/apache/pulsar-client-go/pulsaradmin/pkg/admin/config"
	"github.com/streamnative/pulsarctl/pkg/cmdutils"
)

const (
	DefaultClientTimeout = 30 * time.Second
)

type PulsarContext struct {
	WebServiceURL                 string
	Token                         string
	AuthPlugin                    string
	AuthParams                    string
	TLSAllowInsecureConnection    bool
	TLSEnableHostnameVerification bool
	TLSTrustCertsFilePath         string
	TLSCertFile                   string
	TLSKeyFile                    string
}

var CurrentPulsarContext PulsarContext
var CurrentPulsarClientOptions pulsar.ClientOptions
var AdminClient cmdutils.Client
var AdminV3Client cmdutils.Client
var Client pulsar.Client

func NewCurrentPulsarContext(pc PulsarContext) error {
	CurrentPulsarContext = pc
	return pc.SetPulsarContext()
}

func init() {
	cmdutils.PulsarCtlConfig = &cmdutils.ClusterConfig{}
}

func (pc *PulsarContext) SetPulsarContext() error {
	var err error
	// Configure pulsarctl with the token
	if pc.Token != "" {
		cmdutils.PulsarCtlConfig = &cmdutils.ClusterConfig{
			WebServiceURL: pc.WebServiceURL,
			AuthPlugin:    "org.apache.pulsar.client.impl.auth.AuthenticationToken",
			AuthParams:    fmt.Sprintf("token:%s", pc.Token),
		}

		// Set the global client options
		CurrentPulsarClientOptions = pulsar.ClientOptions{
			URL:               pc.WebServiceURL,
			Authentication:    pulsar.NewAuthenticationToken(pc.Token),
			OperationTimeout:  DefaultClientTimeout,
			ConnectionTimeout: DefaultClientTimeout,
		}
	} else if pc.AuthPlugin != "" && pc.AuthParams != "" {
		cmdutils.PulsarCtlConfig = &cmdutils.ClusterConfig{
			WebServiceURL: pc.WebServiceURL,
			AuthPlugin:    pc.AuthPlugin,
			AuthParams:    pc.AuthParams,
		}

		authProvider, err := pulsar.NewAuthentication(pc.AuthPlugin, pc.AuthParams)
		if err != nil {
			return fmt.Errorf("failed to create authentication provider: %w", err)
		}
		CurrentPulsarClientOptions = pulsar.ClientOptions{
			URL:               pc.WebServiceURL,
			Authentication:    authProvider,
			OperationTimeout:  DefaultClientTimeout,
			ConnectionTimeout: DefaultClientTimeout,
		}
	}

	AdminClient = cmdutils.NewPulsarClient()
	AdminV3Client = cmdutils.NewPulsarClientWithAPIVersion(pulsaradminconfig.V3)

	Client, err = pulsar.NewClient(CurrentPulsarClientOptions)
	if err != nil {
		return fmt.Errorf("failed to create pulsar client: %w", err)
	}

	return nil
}

func GetAdminClient() (cmdutils.Client, error) {
	if cmdutils.PulsarCtlConfig.WebServiceURL == "" {
		return nil, fmt.Errorf("ContextNotSetErr: Please set the cluster context first.")
	}
	return AdminClient, nil
}

func GetAdminV3Client() (cmdutils.Client, error) {
	if cmdutils.PulsarCtlConfig.WebServiceURL == "" {
		return nil, fmt.Errorf("ContextNotSetErr: Please set the cluster context first.")
	}
	return AdminV3Client, nil
}

func GetPulsarClient() (pulsar.Client, error) {
	if CurrentPulsarClientOptions.URL == "" {
		return nil, fmt.Errorf("ContextNotSetErr: Please set the cluster context first.")
	}
	return Client, nil
}
