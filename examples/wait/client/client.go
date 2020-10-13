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
	for {
		if value, err := conn.Wait("foo"); err != nil {
			panic(err)
		} else {
			fmt.Printf("Wait foo:%s\n", string(value))
		}
	}
}
