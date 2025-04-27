package kafka

type KafkaContext struct {
	BootstrapServers  string
	AuthType          string
	AuthMechanism     string
	AuthUser          string
	AuthPass          string
	TLSCertFile       string
	TLSKeyFile        string
	TLSCAPath         string
	TLSAllowInsecure  bool
	SchemaRegistryURL string
	ConnectURL        string

	SchemaRegistryAuthUser string
	SchemaRegistryAuthPass string

	ConnectAuthUser string
	ConnectAuthPass string
}
