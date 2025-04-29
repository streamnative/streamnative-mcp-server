package config

type ExternalPulsar struct {
	WebServiceURL                 string
	AuthPlugin                    string
	AuthParams                    string
	TLSAllowInsecureConnection    bool
	TLSEnableHostnameVerification bool
	TLSTrustCertsFilePath         string
	TLSCertFile                   string
	TLSKeyFile                    string
}
