package redis

import (
	"fmt"
	"io"

	"github.com/garyburd/redigo/redis"
	"github.com/go-jar/golog"
	"github.com/goinbox/gomisc"
)

type Cmd struct {
	cmd  string
	args []interface{}
}

func NewCmd(cmd string, args []interface{}) *Cmd {
	return &Cmd{
		cmd:  cmd,
		args: args,
	}
}

type Client struct {
	config    *Config
	logger    golog.ILogger
	logPrefix []byte
	traceId   []byte

	conn        redis.Conn
	isConnected bool

	cmdQueue []*Cmd
}

func NewClient(config *Config, logger golog.ILogger) *Client {
	if logger == nil {
		logger = new(golog.NoopLogger)
	}

	c := &Client{
		config:  config,
		logger:  logger,
		traceId: []byte("-"),

		cmdQueue: []*Cmd{},
	}

	c.logPrefix = []byte("[Redis " + c.config.Host + ":" + c.config.Port + "]\t")
	return c
}

func (c *Client) SetLogger(logger golog.ILogger) *Client {
	c.logger = logger
	return c
}

func (c *Client) SetTraceId(traceId []byte) *Client {
	c.traceId = traceId
	return c
}

func (c *Client) Connect() error {
	options := []redis.DialOption{
		redis.DialConnectTimeout(c.config.ConnectTimeout),
		redis.DialReadTimeout(c.config.ReadTimeout),
		redis.DialWriteTimeout(c.config.WriteTimeout),
	}

	conn, err := redis.Dial("tcp", c.config.Host+":"+c.config.Port, options...)
	if err != nil {
		return err
	}

	_, err = conn.Do("auth", c.config.Pass)
	if err != nil {
		return err
	}

	c.conn = conn
	c.isConnected = true

	return nil
}

func (c *Client) IsConnected() bool {
	return c.isConnected
}

func (c *Client) Free() {
	if c.conn != nil {
		c.conn.Close()
	}

	c.isConnected = false
}

func (c *Client) Do(cmd string, args ...interface{}) *Reply {
	if !c.isConnected {
		if err := c.Connect(); err != nil {
			return NewReply(nil, err)
		}
	}

	c.log(cmd, args...)

	return c.do(cmd, args...)
}

func (c *Client) DoWithoutLog(cmd string, args ...interface{}) *Reply {
	if !c.isConnected {
		if err := c.Connect(); err != nil {
			return NewReply(nil, err)
		}
	}

	return c.do(cmd, args...)
}

func (c *Client) do(cmd string, args ...interface{}) *Reply {
	if !c.isConnected {
		err := c.Connect()
		if err != nil {
			return NewReply(nil, err)
		}
	}

	for _, cmd := range c.cmdQueue {
		c.conn.Send(cmd.cmd, cmd.args...)
	}

	reply, err := c.conn.Do(cmd, args...)
	if err != nil {
		if err != io.EOF {
			return NewReply(nil, err)
		}

		if !c.config.IsTimeoutAutoConnect {
			return NewReply(nil, err)
		}

		c.reconnect()

		for _, cmd := range c.cmdQueue {
			c.conn.Send(cmd.cmd, cmd.args...)
		}

		reply, err = c.conn.Do(cmd, args...)
		if err != nil {
			return NewReply(nil, err)
		}
	}

	return NewReply(reply, err)
}

func (c *Client) Send(cmd string, args ...interface{}) {
	c.log(cmd, args...)
	c.cmdQueue = append(c.cmdQueue, NewCmd(cmd, args))
}

func (c *Client) FlushCmdQueue() ([]*Reply, []int) {
	if !c.isConnected {
		err := c.Connect()
		if err != nil {
			return []*Reply{NewReply(nil, err)}, []int{0}
		}
	}

	defer func() {
		c.cmdQueue = []*Cmd{}
	}()

	reply, err := c.flushCmdQueue()
	if err != nil {
		return []*Reply{NewReply(nil, err)}, []int{0}
	}

	if err != nil {
		if err != io.EOF {
			return []*Reply{NewReply(nil, err)}, []int{0}
		}

		if !c.config.IsTimeoutAutoConnect {
			return []*Reply{NewReply(nil, err)}, []int{0}
		}

		if err := c.reconnect(); err != nil {
			return []*Reply{NewReply(nil, err)}, []int{0}
		}

		reply, err = c.flushCmdQueue()
		if err != nil {
			return []*Reply{NewReply(nil, err)}, []int{0}
		}
	}

	replies := make([]*Reply, len(c.cmdQueue))
	var errIndexes []int
	replies[0] = NewReply(reply, nil)

	for i := 1; i < len(c.cmdQueue); i++ {
		reply, err := c.conn.Receive()
		replies[i] = NewReply(reply, err)

		if err != nil {
			errIndexes = append(errIndexes, i)
		}
	}

	return replies, errIndexes
}

func (c *Client) flushCmdQueue() (interface{}, error) {
	for _, cmd := range c.cmdQueue {
		if err := c.conn.Send(cmd.cmd, cmd.args...); err != nil {
			return nil, err
		}
	}

	if err := c.conn.Flush(); err != nil {
		return nil, err
	}

	reply, err := c.conn.Receive()
	if err != nil {
		return nil, err
	}

	return reply, nil
}

func (c *Client) BeginTrans() {
	c.Send("multi")
}

func (c *Client) DiscardTrans() error {
	return c.Do("discard").Err
}

func (c *Client) ExecTrans() ([]*Reply, error) {
	reply := c.Do("exec")
	values, err := redis.Values(reply.reply, reply.Err)
	if err != nil {
		return nil, err
	}

	replies := make([]*Reply, len(values))
	for i, value := range values {
		replies[i] = NewReply(value, nil)
	}

	return replies, nil
}

func (c *Client) reconnect() error {
	c.Free()
	return c.Connect()
}

func (c *Client) log(cmd string, args ...interface{}) {
	if cmd == "" {
		return
	}

	for _, arg := range args {
		cmd += " " + fmt.Sprint(arg)
	}

	c.logger.Log(c.config.LogLevel, gomisc.AppendBytes(c.logPrefix, []byte("\t"), c.traceId, []byte("\t"), []byte(cmd)))
}
