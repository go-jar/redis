package redis

import (
	"github.com/go-jar/golog"
	"time"
)

const (
	DefaultConnectTimeout = 10 * time.Second
	DefaultReadTimeout    = 10 * time.Second
	DefaultWriteTimeout   = 10 * time.Second
)

type Config struct {
	LogLevel int

	Host string
	Port string
	Pass string

	ConnectTimeout time.Duration
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration

	IsTimeoutAutoConnect bool
}

func NewConfig(host, port, pass string) *Config {
	return &Config{
		LogLevel: golog.LevelInfo,

		Host: host,
		Port: port,
		Pass: pass,

		ConnectTimeout: DefaultConnectTimeout,
		ReadTimeout:    DefaultReadTimeout,
		WriteTimeout:   DefaultWriteTimeout,
	}
}
