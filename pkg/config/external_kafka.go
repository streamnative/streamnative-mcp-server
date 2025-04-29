package config

type ExternalKafka struct {
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

	SchemaRegistryAuthUser    string
	SchemaRegistryAuthPass    string
	SchemaRegistryBearerToken string
}
