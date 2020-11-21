package redis

import (
	"fmt"
	"testing"
	"time"

	"github.com/go-jar/pool"
)

func TestPool(t *testing.T) {
	config := &pool.Config{
		MaxConns:          100,
		MaxIdleTime:       time.Second * 5,
		KeepAliveInterval: time.Second * 3,
	}

	pool := NewPool(config, newRedisTestClient, true)
	testPool(pool, t)
}

func newRedisTestClient() (*Client, error) {
	return getTestClient(), nil
}

func testPool(p *Pool, t *testing.T) {
	client, _ := p.Get()
	client.Do("set", "redis_pool", "pool_test")
	reply := client.Do("get", "redis_pool")
	fmt.Println(reply.String())
	p.Put(client)

	time.Sleep(time.Second * 4)
	client, _ = p.Get()
	reply = client.Do("get", "redis_pool")
	fmt.Println(reply.String())
	p.Put(client)
}
