package redis

import (
	"encoding/json"
	"errors"
	"reflect"

	"github.com/go-jar/golog"
)

type SimpleOrm struct {
	traceId []byte
	pool    *Pool

	client *Client
	logger golog.ILogger
}

func NewSimpleOrm(traceId []byte, pool *Pool) *SimpleOrm {
	return &SimpleOrm{
		traceId: traceId,
		pool:    pool,
		logger:  new(golog.NoopLogger),
	}
}

func (so *SimpleOrm) Client() *Client {
	if so.client == nil {
		so.client, _ = so.pool.Get()
		so.client.SetLogger(so.logger).SetTraceId(so.traceId)
	}

	return so.client
}

func (so *SimpleOrm) PutBackClient() {
	if so.client.IsConnected() {
		so.client.SetLogger(new(golog.NoopLogger))
		_ = so.pool.Put(so.client)
	}

	so.client = nil
}

func (so *SimpleOrm) Renew(traceId []byte, pool *Pool) *SimpleOrm {
	if so.client != nil {
		so.PutBackClient()
	}

	so.traceId = traceId
	so.pool = pool

	return so
}

func (so *SimpleOrm) SaveAsJson(key string, value interface{}, expireSeconds int64) error {
	jsonBytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	client := so.Client()
	defer so.PutBackClient()

	if expireSeconds > 0 {
		err = client.Do("set", key, string(jsonBytes), "ex", expireSeconds).Err
	} else {
		err = client.Do("set", key, string(jsonBytes)).Err
	}

	return err
}

func (so *SimpleOrm) GetAsJson(key string, value interface{}) (bool, error) {
	reply := so.Client().Do("get", key)
	defer so.PutBackClient()

	if reply.Err != nil {
		return false, reply.Err
	}

	if reply == nil {
		return false, nil
	}

	jsonBytes, err := reply.Bytes()
	if err != nil {
		return false, err
	}

	err = json.Unmarshal(jsonBytes, value)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (so *SimpleOrm) SaveEntity(key string, entity interface{}, expireSeconds int64) error {
	rev := reflect.ValueOf(entity)
	entityArgs := ReflectSaveEntityArgs(rev)

	args := make([]interface{}, len(entityArgs)+1)
	args[0] = key
	for i, arg := range entityArgs {
		args[i+1] = arg
	}

	client := so.Client()
	defer so.PutBackClient()

	client.Send("hmset", args...)
	if expireSeconds > 0 {
		client.Send("expire", key, expireSeconds)
	}
	replies, errIndexes := client.FlushCmdQueue()
	if len(errIndexes) != 0 {
		msg := "hmset " + key + " to redis error:"
		for _, i := range errIndexes {
			msg += " " + replies[i].Err.Error()
		}
		return errors.New(msg)
	}

	return nil
}

func (so *SimpleOrm) GetEntity(key string, entityPtr interface{}) (bool, error) {
	reply := so.Client().Do("hgetall", key)
	defer so.PutBackClient()

	if reply.Err != nil {
		return false, reply.Err
	}

	if reply.ArrReplyIsNil() {
		return false, nil
	}

	err := reply.Struct(entityPtr)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (so *SimpleOrm) Del(key string) error {
	err := so.Client().Do("del", key).Err
	defer so.PutBackClient()

	return err
}
