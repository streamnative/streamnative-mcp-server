package kafka

import (
	"crypto/tls"
	"fmt"
	"strings"

	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sasl/plain"
	"github.com/twmb/franz-go/pkg/sasl/scram"
	"github.com/twmb/franz-go/pkg/sr"
	"github.com/twmb/tlscfg"
)

type KafkaContext struct {
	BootstrapServers  string
	AuthType          string
	AuthMechanism     string
	AuthUser          string
	AuthPass          string
	UseTLS            bool
	ClientKeyFile     string
	ClientCertFile    string
	CaFile            string
	SchemaRegistryURL string
	ConnectURL        string

	SchemaRegistryAuthUser    string
	SchemaRegistryAuthPass    string
	SchemaRegistryBearerToken string

	ConnectAuthUser string
	ConnectAuthPass string
}

var CurrentKafkaContext KafkaContext
var KafkaAdminClient *kadm.Client
var KafkaClient *kgo.Client
var KafkaSchemaRegistryClient *sr.Client
var KafkaConnectClient Connect
var options []kgo.Opt

func NewCurrentKafkaContext(kc KafkaContext) error {
	CurrentKafkaContext = kc
	return kc.SetKafkaContext()
}

type SASLConfig struct {
	Mechanism string
	Username  string
	Password  string
}

type TLSConfig struct {
	Enabled        bool
	ClientKeyFile  string
	ClientCertFile string
	CaFile         string
}

// Initializes the necessary TLS configuration options
func tlsOpt(config *TLSConfig, opts []kgo.Opt) ([]kgo.Opt, error) {
	if config.Enabled {
		if config.CaFile != "" || config.ClientCertFile != "" || config.ClientKeyFile != "" {
			tc, err := tlscfg.New(
				tlscfg.MaybeWithDiskCA(config.CaFile, tlscfg.ForClient),
				tlscfg.MaybeWithDiskKeyPair(config.ClientCertFile, config.ClientKeyFile),
			)
			if err != nil {
				return nil, fmt.Errorf("unable to create TLS config: %v", err)
			}
			opts = append(opts, kgo.DialTLSConfig(tc))
		} else {
			opts = append(opts, kgo.DialTLSConfig(new(tls.Config)))
		}
	}
	return opts, nil
}

// Initializes the necessary SASL configuration options
func saslOpt(config *SASLConfig, opts []kgo.Opt) ([]kgo.Opt, error) {
	if config.Mechanism != "" || config.Username != "" || config.Password != "" {
		if config.Mechanism == "" || config.Username == "" || config.Password == "" {
			return nil, fmt.Errorf("All of Mechanism, Username, and Password must be specified if any are")
		}
		method := strings.ToLower(config.Mechanism)
		method = strings.ReplaceAll(method, "-", "")
		method = strings.ReplaceAll(method, "_", "")
		switch method {
		case "plain":
			opts = append(opts, kgo.SASL(plain.Auth{
				User: config.Username,
				Pass: config.Password,
			}.AsMechanism()))
		case "scramsha256":
			opts = append(opts, kgo.SASL(scram.Auth{
				User: config.Username,
				Pass: config.Password,
			}.AsSha256Mechanism()))
		case "scramsha512":
			opts = append(opts, kgo.SASL(scram.Auth{
				User: config.Username,
				Pass: config.Password,
			}.AsSha512Mechanism()))
		default:
			return nil, fmt.Errorf("unrecognized SASL method: %s", config.Mechanism)
		}
	}
	return opts, nil
}

func (kc *KafkaContext) SetKafkaContext() error {
	var err error
	options = []kgo.Opt{}
	options = append(options, kgo.SeedBrokers(strings.Split(kc.BootstrapServers, ",")...))
	tlsConfig := &TLSConfig{
		Enabled:        kc.UseTLS,
		ClientKeyFile:  kc.ClientKeyFile,
		ClientCertFile: kc.ClientCertFile,
		CaFile:         kc.CaFile,
	}

	saslConfig := &SASLConfig{
		Mechanism: kc.AuthMechanism,
		Username:  kc.AuthUser,
		Password:  kc.AuthPass,
	}

	options, err = tlsOpt(tlsConfig, options)
	if err != nil {
		return fmt.Errorf("failed to create TLS config: %w", err)
	}
	options, err = saslOpt(saslConfig, options)
	if err != nil {
		return fmt.Errorf("failed to create SASL config: %w", err)
	}

	KafkaClient, err = kgo.NewClient(
		options...,
	)
	if err != nil {
		return fmt.Errorf("failed to create kafka client: %w", err)
	}

	KafkaAdminClient = kadm.NewClient(KafkaClient)
	if kc.SchemaRegistryURL != "" {
		SrOpts := []sr.ClientOpt{}
		SrOpts = append(SrOpts, sr.URLs(kc.SchemaRegistryURL))
		if kc.SchemaRegistryAuthUser != "" && kc.SchemaRegistryAuthPass != "" {
			SrOpts = append(SrOpts, sr.BasicAuth(kc.SchemaRegistryAuthUser, kc.SchemaRegistryAuthPass))
		} else if kc.SchemaRegistryBearerToken != "" {
			SrOpts = append(SrOpts, sr.BearerToken(kc.SchemaRegistryBearerToken))
		}
		SrOpts = append(SrOpts, sr.UserAgent("streamnative-mcp-server"))
		KafkaSchemaRegistryClient, err = sr.NewClient(SrOpts...)
		if err != nil {
			return fmt.Errorf("failed to create kafka schema registry client: %w", err)
		}
	}

	if kc.ConnectURL != "" {
		KafkaConnectClient, err = NewConnect(kc)
		if err != nil {
			return fmt.Errorf("failed to create kafka connect client: %w", err)
		}
	}
	return nil
}

func GetKafkaClient(opts ...kgo.Opt) (*kgo.Client, error) {
	if KafkaClient == nil {
		return nil, fmt.Errorf("kafka client not initialized")
	}
	clientOpts := append(options, opts...)
	cli, err := kgo.NewClient(clientOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka client: %w", err)
	}
	return cli, nil
}

func GetKafkaAdminClient() (*kadm.Client, error) {
	if KafkaAdminClient == nil {
		return nil, fmt.Errorf("kafka admin client not initialized")
	}
	return KafkaAdminClient, nil
}

func GetKafkaSchemaRegistryClient() (*sr.Client, error) {
	if CurrentKafkaContext.SchemaRegistryURL == "" {
		return nil, fmt.Errorf("schema registry not enabled on the current context")
	}
	if KafkaSchemaRegistryClient == nil {
		return nil, fmt.Errorf("kafka schema registry client not initialized")
	}
	return KafkaSchemaRegistryClient, nil
}

func GetKafkaConnectClient() (Connect, error) {
	if CurrentKafkaContext.ConnectURL == "" {
		return nil, fmt.Errorf("kafka connect not enabled on the current context")
	}
	if KafkaConnectClient == nil {
		return nil, fmt.Errorf("kafka connect client not initialized")
	}
	return KafkaConnectClient, nil
}
