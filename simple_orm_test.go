package redis

import (
	"github.com/goinbox/gomisc"
	"testing"
	"time"

	"github.com/go-jar/pool"
)

type SqlBaseEntity struct {
	Id       int64  `redis:"id" json:"id"`
	AddTime  string `redis:"add_time" json:"add_time"`
	EditTime string `redis:"edit_time" json:"edit_time"`
}

type DemoEntity struct {
	SqlBaseEntity

	Name   string `redis:"name" json:"name"`
	Status int    `redis:"status" json:"status"`
}

var pl *Pool

func init() {
	config := &pool.Config{
		MaxConns:          100,
		MaxIdleTime:       time.Second * 5,
		KeepAliveInterval: time.Second * 3,
	}

	pl = NewPool(config, newRedisTestClient, true)
}

func TestSetGetJson(t *testing.T) {
	so := NewSimpleOrm([]byte("TestSetGetJson"), pl)

	item := &DemoEntity{
		SqlBaseEntity: SqlBaseEntity{
			Id:       1,
			AddTime:  time.Now().Format(gomisc.TimeGeneralLayout()),
			EditTime: time.Now().Format(gomisc.TimeGeneralLayout()),
		},
		Name:   "tdj",
		Status: 1,
	}

	key := "test_demo_json"
	err := so.SaveAsJson(key, item, 10)
	if err != nil {
		t.Error(err)
	}

	item = &DemoEntity{}
	find, err := so.GetAsJson(key, item)
	if !find {
		t.Error("not found")
	}
	if err != nil {
		t.Error(err)
	}

	t.Log(item)
}

func TestSetGetEntity(t *testing.T) {
	so := NewSimpleOrm([]byte("TestSetGetEntity"), pl)

	item := &DemoEntity{
		SqlBaseEntity: SqlBaseEntity{
			Id:       1,
			AddTime:  time.Now().Format(gomisc.TimeGeneralLayout()),
			EditTime: time.Now().Format(gomisc.TimeGeneralLayout()),
		},
		Name:   "tdj",
		Status: 1,
	}

	key := "test_demo_hash"
	err := so.SaveEntity(key, item, 10)
	if err != nil {
		t.Error(err)
	}

	item = &DemoEntity{}
	find, err := so.GetEntity(key, item)
	if !find {
		t.Error("not found")
	}
	if err != nil {
		t.Error(err)
	}

	t.Log(item)
}
