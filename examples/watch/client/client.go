package main

import (
	"fmt"
	"github.com/hslam/rpc"
)

func main() {
	conn, err := rpc.Dial("tcp", ":9999", "pb")
	if err != nil {
		panic(err)
	}
	watch := conn.Watch("foo", nil)
	for {
		<-watch.Done
		if watch.Error != nil {
			panic(watch.Error)
		}
		fmt.Printf("Watch foo:%s\n", string(watch.Value))
	}
}
