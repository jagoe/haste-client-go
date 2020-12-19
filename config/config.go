package config

// HasteConfig represents the configuration values provided by a config YAML, ENV or flags
type HasteConfig struct {
	Server                   string `mapstructure:"server"`
	ClientCertificatePath    string `mapstructure:"clientCert"`
	ClientCertificateKeyPath string `mapstructure:"clientCertKey"`
}

// CanProvideClientCertificate checks if the config contains all necessary information to provide a client certificate
func (config *HasteConfig) CanProvideClientCertificate() bool {
	return len(config.ClientCertificatePath) > 0 && len(config.ClientCertificateKeyPath) > 0
}

// GetConfig represents the configuration values provided bt a config YAML, ENV or flags for the get command
type GetConfig struct {
	HasteConfig HasteConfig `mapstructure:",squash"`
	OutputPath  string
}

// ShouldSaveAsFile checks if the config contains an output path
func (config *GetConfig) ShouldSaveAsFile() bool {
	return len(config.OutputPath) > 0
}

// CreateConfig represents the configuration values provided bt a config YAML, ENV or flags for creating a haste
type CreateConfig struct {
	HasteConfig HasteConfig `mapstructure:",squash"`
}
