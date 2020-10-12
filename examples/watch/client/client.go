package main

import (
	"fmt"
	"github.com/hslam/rpc"
	"time"
)

func main() {
	conn, err := rpc.Dial("tcp", ":9999", "pb")
	if err != nil {
		panic(err)
	}
	watch := conn.Watch("changes", nil)
	for {
		<-watch.Done
		if watch.Error != nil {
			panic(watch.Error)
		}
		fmt.Printf("Watch changes - %v\n", time.Now().Format(time.Stamp))
	}
}
