package redis

import "github.com/go-jar/pool"

type Pool struct {
	pl            *pool.Pool
	NewClientFunc func() (*Client, error)
	LogKeepAlive  bool
}

type NewClientFunc func() (*Client, error)

func NewPool(config *pool.Config, ncf NewClientFunc, LogKeepAlive bool) *Pool {
	p := &Pool{
		pl:            pool.NewPool(config, nil),
		NewClientFunc: ncf,
		LogKeepAlive:  LogKeepAlive,
	}
	p.pl.NewItemFunc = p.newConn

	if config.KeepAliveFunc == nil {
		config.KeepAliveFunc = p.keepAlive
	}

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

	if p.LogKeepAlive == true {
		return client.Do("ping").Err
	}

	return client.DoWithoutLog("ping").Err
}

func (p *Pool) newConn() (pool.IConn, error) {
	return p.NewClientFunc()
}
