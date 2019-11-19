package main
import (
	"hslam.com/git/x/rpc"
	"log"
	"fmt"
)
func main()  {
	conn, err:= rpc.Dial("tcp","127.0.0.1:9999","bytes")//tcp|ws|quic|http
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