package config

import (
	"time"
)

type Config struct {
	LogLevel            string              `yaml:"log_level" required:"true"`
	Tracing             Tracing             `yaml:"tracing" required:"true"`
	Bot                 Bot                 `yaml:"bot" required:"true"`
	Postgres            Postgres            `yaml:"postgres" required:"true"`
	ButtonRedis         ButtonRedis         `yaml:"button_redis" required:"true"`
	RandomNameGenerator RandomNameGenerator `yaml:"random_name_generator" required:"true"`
	ChangeNameInterval  time.Duration       `yaml:"change_name_interval" required:"true"`
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

type Postgres struct {
	Connection string `yaml:"connection" required:"true"`
}

type ButtonRedis struct {
	Addr string        `yaml:"addr" required:"true"`
	Pwd  string        `yaml:"pwd" required:"true"`
	DB   int           `yaml:"db" required:"true"`
	TTL  time.Duration `yaml:"ttl" required:"true"`
}

type RandomNameGenerator struct {
	Prefix string `yaml:"prefix" required:"true"`
	Length int    `yaml:"length" required:"true"`
}
