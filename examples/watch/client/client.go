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
	watch := conn.Watch("foo")
	defer watch.Stop()
	for i := 0; i < 3; i++ {
		value, err := watch.Wait()
		if err != nil {
			panic(err)
		}
		fmt.Printf("Watch foo:%s\n", string(value))
	}
}
