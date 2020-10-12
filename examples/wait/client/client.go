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
	for {
		if err := conn.Wait("changes"); err != nil {
			panic(err)
		}
		fmt.Printf("Wait changes - %v\n", time.Now().Format(time.Stamp))
	}
}
