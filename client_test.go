package redis

import (
	"fmt"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	client := getTestClient()

	reply := client.Do("set", "c", "1")
	fmt.Println(reply.String())
	reply = client.Do("get", "c")
	fmt.Println(reply.Int())

	reply = client.DoWithoutLog("set", "d", "1")
	fmt.Println(reply.String())
	reply = client.DoWithoutLog("get", "d")
	fmt.Println(reply.Int())

	client.Send("set", "a", "a")
	client.Send("set", "b", "b")
	client.Send("get", "a")
	client.Send("get", "b")
	replies, errIndexes := client.FlushCmdQueue()
	fmt.Println(errIndexes)
	for _, reply := range replies {
		fmt.Println(reply.String())
		fmt.Println(reply.Err)
	}

	client.BeginTrans()
	client.Send("set", "a", "1")
	client.Send("set", "b", "2")
	client.Send("get", "a")
	client.Send("get", "b")
	replies, _ = client.ExecTrans()
	for _, reply := range replies {
		fmt.Println(reply.String())
		fmt.Println(reply.Err)
	}

	time.Sleep(time.Second * 5)
	reply = client.Do("get", "c")
	fmt.Println(reply.Int())

	client.Free()
}

func TestAutoReconnect(t *testing.T) {
	client := getTestClient()

	reply := client.Do("set", "a", "1")
	fmt.Println(reply.String())
	time.Sleep(time.Second * 4) //set redis-server timeout = 3
	reply = client.Do("get", "a")
	fmt.Println(reply.Err)
	fmt.Println(reply.Int())

	time.Sleep(time.Second * 4)

	client.Send("set", "a", "a")
	client.Send("set", "b", "b")
	client.Send("get", "a")
	client.Send("get", "b")
	replies, errIndexes := client.FlushCmdQueue()
	fmt.Println(errIndexes)
	for _, reply := range replies {
		fmt.Println(reply.String())
		fmt.Println(reply.Err)
	}

	time.Sleep(time.Second * 4)

	client.BeginTrans()
	client.Send("set", "a", "1")
	client.Send("set", "b", "2")
	client.Send("get", "a")
	client.Send("get", "b")
	replies, _ = client.ExecTrans()
	for _, reply := range replies {
		fmt.Println(reply.String())
		fmt.Println(reply.Err)
	}

	client.Free()
}
