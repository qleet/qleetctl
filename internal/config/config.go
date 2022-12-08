package config

// QleetConfig is the client's configuration for connecting to QleetOS instances
type QleetConfig struct {
	QleetOSInstances []QleetOSInstance `yaml:"QleetOSInstances"`
	CurrentInstance  string            `yaml:"CurrentInstance"`
}

// QleetOSInstance is an instance of QleetOS the client can use
type QleetOSInstance struct {
	Name      string `yaml:"Name"`
	APIServer string `yaml:"APIServer"`
}
