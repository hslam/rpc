package main
import (
	"github.com/hslam/rpc"
	"log"
	"fmt"
)
func main()  {
	conn, err:= rpc.Dial("ipc","/tmp/ipc","bytes")
	if err != nil {
		log.Fatalln("dailing error: ", err)
	}
	defer conn.Close()
	var req =[]byte("Hello World")
	var res []byte
	err = conn.Call("Echo.ToLower", &req, &res)
	if err != nil {
		log.Fatalln("Echo error: ", err)
	}
	fmt.Printf("Echo.ToLower : %s\n", string(res))
}