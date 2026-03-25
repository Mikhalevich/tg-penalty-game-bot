package config

import (
	"time"
)

type Config struct {
	LogLevel           string             `yaml:"log_level" required:"true"`
	Tracing            Tracing            `yaml:"tracing" required:"true"`
	Bot                Bot                `yaml:"bot" required:"true"`
	Postgres           Postgres           `yaml:"postgres" required:"true"`
	MatchmakingWorker  Worker             `yaml:"matchmaking_worker" required:"true"`
	ExpiredShotsWorker ExpiredShotsWorker `yaml:"expired_shots_worker" required:"true"`
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

type Worker struct {
	Count     int           `yaml:"count" required:"true"`
	Interval  time.Duration `yaml:"interval" required:"true"`
	BatchSize int           `yaml:"batch_size" required:"true"`
}

type ExpiredShotsWorker struct {
	Count              int           `yaml:"count" required:"true"`
	Interval           time.Duration `yaml:"interval" required:"true"`
	BatchSize          int           `yaml:"batch_size" required:"true"`
	ShotExpireDuration time.Duration `yaml:"shot_expire_duration" required:"true"`
}
