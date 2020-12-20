package config

// Config represents the configuration values provided by a config YAML, ENV or flags
type Config struct {
	Server                   string `mapstructure:"server"`
	ClientCertificatePath    string `mapstructure:"clientCert"`
	ClientCertificateKeyPath string `mapstructure:"clientCertKey"`
}

// CanProvideClientCertificate checks if the config contains all necessary information to provide a client certificate
func (config *Config) CanProvideClientCertificate() bool {
	return len(config.ClientCertificatePath) > 0 && len(config.ClientCertificateKeyPath) > 0
}
