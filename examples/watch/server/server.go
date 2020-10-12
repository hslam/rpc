package main

import (
	"github.com/hslam/rpc"
	"time"
)

func main() {
	go func() {
		ticker := time.NewTicker(time.Second * 3)
		for range ticker.C {
			rpc.Trigger("changes")
		}
	}()
	rpc.Listen("tcp", ":9999", "pb")
}
