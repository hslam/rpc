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
	start := time.Now()
	if err := conn.Ping(); err != nil {
		panic(err)
	}
	fmt.Printf("time=%v", time.Now().Sub(start))
}
