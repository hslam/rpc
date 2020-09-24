package main

import (
	"github.com/hslam/rpc"
)

func main() {
	rpc.Listen("tcp", ":9999", "pb")
}
