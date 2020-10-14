package main

import (
	"fmt"
	"github.com/hslam/rpc"
	"time"
)

func main() {
	go func() {
		ticker := time.NewTicker(time.Second * 3)
		for range ticker.C {
			rpc.Push("foo", []byte(fmt.Sprintf("bar - %s", time.Now().Format(time.Stamp))))
		}
	}()
	rpc.Listen("tcp", ":9999", "pb")
}
