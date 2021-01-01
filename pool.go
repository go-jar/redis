package redis

import "github.com/go-jar/pool"

type PoolConfig struct {
	pool.Config

	NewClientFunc func() (*Client, error)

	LogKeepAlive bool
}

type Pool struct {
	pl *pool.Pool

	config *PoolConfig
}

type NewClientFunc func() (*Client, error)

func NewPool(config *PoolConfig) *Pool {
	p := &Pool{
		config: config,
	}

	if config.NewConnFunc == nil {
		config.NewConnFunc = p.newConn
	}

	if config.KeepAliveFunc == nil {
		config.KeepAliveFunc = p.keepAlive
	}

	p.pl = pool.NewPool(&p.config.Config)

	return p
}

func (p *Pool) Get() (*Client, error) {
	conn, err := p.pl.Get()
	if err != nil {
		return nil, err
	}
	return conn.(*Client), err
}

func (p *Pool) Put(client *Client) error {
	if client.IsConnected() {
		return p.pl.Put(client)
	}
	return nil
}

func (p *Pool) keepAlive(conn pool.IConn) error {
	client := conn.(*Client)

	if p.config.LogKeepAlive == true {
		return client.Do("ping").Err
	}

	return client.DoWithoutLog("ping").Err
}

func (p *Pool) newConn() (pool.IConn, error) {
	return p.config.NewClientFunc()
}
