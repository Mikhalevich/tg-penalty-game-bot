package config

import (
	"time"
)

type Config struct {
	LogLevel            string              `yaml:"log_level" required:"true"`
	Tracing             Tracing             `yaml:"tracing" required:"true"`
	Bot                 Bot                 `yaml:"bot" required:"true"`
	Postgres            Postgres            `yaml:"postgres" required:"true"`
	ButtonRedis         ButtonRedis         `yaml:"button_redis"`
	RandomNameGenerator RandomNameGenerator `yaml:"random_name_generator" required:"true"`
	ChangeName          ChangeName          `yaml:"change_name" required:"true"`
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
	Token        string `yaml:"token" required:"true"`
	WebHookToken string `yaml:"webhook_token"`
}

type Postgres struct {
	Connection string `yaml:"connection" required:"true"`
}

type ButtonRedis struct {
	Addr string        `yaml:"addr"`
	Pwd  string        `yaml:"pwd"`
	DB   int           `yaml:"db"`
	TTL  time.Duration `yaml:"ttl"`
}

type RandomNameGenerator struct {
	Prefix string `yaml:"prefix" required:"true"`
	Length int    `yaml:"length" required:"true"`
}

type ChangeName struct {
	MaxLen        int           `yaml:"max_len" required:"true"`
	RetryInterval time.Duration `yaml:"retry_interval" required:"true"`
}
