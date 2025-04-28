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

	AdminClient = cmdutils.NewPulsarClient()
	AdminV3Client = cmdutils.NewPulsarClientWithAPIVersion(pulsaradminconfig.V3)

	Client, err = pulsar.NewClient(CurrentPulsarClientOptions)
	if err != nil {
		return err
	}

	return nil
}

func GetAdminClient() (cmdutils.Client, error) {
	if cmdutils.PulsarCtlConfig.WebServiceURL == "" {
		return nil, fmt.Errorf("Please set the cluster context first.")
	}
	return AdminClient, nil
}

func GetAdminV3Client() (cmdutils.Client, error) {
	if cmdutils.PulsarCtlConfig.WebServiceURL == "" {
		return nil, fmt.Errorf("Please set the cluster context first.")
	}
	return AdminV3Client, nil
}

func GetPulsarClient() (pulsar.Client, error) {
	if CurrentPulsarClientOptions.URL == "" {
		return nil, fmt.Errorf("Please set the cluster context first.")
	}
	return Client, nil
}
