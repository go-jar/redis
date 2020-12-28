package redis

import (
	"github.com/go-jar/golog"
	"time"
)

func newRedisTestClient() (*Client, error) {
	return getTestClient(), nil
}

func getTestClient() *Client {
	logger, _ := golog.NewConsoleLogger(golog.LEVEL_INFO)
	config := NewConfig("127.0.0.1", "6379", "passwd")
	config.ConnectTimeout = time.Second * 3

	return NewClient(config, logger)
}
