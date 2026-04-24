package config

import (
	"time"
)

type Config struct {
	LogLevel      string      `yaml:"log_level" required:"true"`
	Tracing       Tracing     `yaml:"tracing" required:"true"`
	Postgres      Postgres    `yaml:"postgres" required:"true"`
	Bot           Bot         `yaml:"bot" required:"true"`
	ButtonRedis   ButtonRedis `yaml:"button_redis"`
	MessageWorker Worker      `yaml:"message_worker" required:"true"`
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

type Postgres struct {
	Connection string `yaml:"connection" required:"true"`
}

type Bot struct {
	Token string `yaml:"token" required:"true"`
}

type ButtonRedis struct {
	Addr string        `yaml:"addr"`
	Pwd  string        `yaml:"pwd"`
	DB   int           `yaml:"db"`
	TTL  time.Duration `yaml:"ttl"`
}

type Worker struct {
	Count         int           `yaml:"count" required:"true"`
	Interval      time.Duration `yaml:"interval" required:"true"`
	BatchSize     int           `yaml:"batch_size" required:"true"`
	MaxRetryCount int           `yaml:"max_retry_count" required:"true"`
}
