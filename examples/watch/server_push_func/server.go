package main

import (
	"fmt"
	"github.com/hslam/rpc"
	"sync"
	"time"
)

func main() {
	var lock sync.Mutex
	var k = "foo"
	var v = []byte(fmt.Sprintf("bar - %s", time.Now().Format(time.Stamp)))
	rpc.PushFunc(func(key string) (value []byte, ok bool) {
		if key == k {
			lock.Lock()
			value = v
			lock.Unlock()
			return value, true
		}
		return nil, false
	})
	go func() {
		ticker := time.NewTicker(time.Second * 3)
		for range ticker.C {
			lock.Lock()
			v = []byte(fmt.Sprintf("bar - %s", time.Now().Format(time.Stamp)))
			value := v
			lock.Unlock()
			rpc.Push(k, value)
		}
	}()
	rpc.Listen("tcp", ":9999", "pb")
}
