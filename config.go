package redis

import (
	"github.com/go-jar/golog"
	"time"
)

const (
	DEFAULT_CONNECT_TIMEOUT = 10 * time.Second
	DEFAULT_READ_TIMEOUT    = 10 * time.Second
	DEFAULT_WRITE_TIMEOUT   = 10 * time.Second
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

		ConnectTimeout: DEFAULT_CONNECT_TIMEOUT,
		ReadTimeout:    DEFAULT_READ_TIMEOUT,
		WriteTimeout:   DEFAULT_WRITE_TIMEOUT,
	}
}
