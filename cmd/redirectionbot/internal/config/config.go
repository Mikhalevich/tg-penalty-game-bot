package config

type Config struct {
	LogLevel        string  `yaml:"log_level" required:"true"`
	Tracing         Tracing `yaml:"tracing" required:"true"`
	Bot             Bot     `yaml:"bot" required:"true"`
	RedirectionLink string  `yaml:"redirection_link" required:"true"`
}

func (c *Config) Level() string {
	return c.LogLevel
}

func (c *Config) ServiceName() string {
	return c.Tracing.ServiceName
}

func (c *Config) TracingEndpoint() string {
	return c.Tracing.Endpoint
}

type Tracing struct {
	Endpoint    string `yaml:"endpoint" required:"true"`
	ServiceName string `yaml:"service_name" required:"true"`
}

type Bot struct {
	Token string `yaml:"token" required:"true"`
}
