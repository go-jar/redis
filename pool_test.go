package redis

import (
	"fmt"
	"testing"
	"time"
)

func TestPool(t *testing.T) {
	config := &PoolConfig{NewClientFunc: newRedisTestClient}

	config.MaxConns = 100
	config.MaxIdleTime = time.Second * 5
	config.KeepAliveInterval = time.Second * 3

	pool := NewPool(config)
	testPool(pool, t)
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
