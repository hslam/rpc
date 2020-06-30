package rpc

import (
	"github.com/hslam/transport"
	"os"
)

func ListenAndServe(tran transport.Transport, address string, codec *Codec) {
	logger.Noticef("pid %d", os.Getpid())
	logger.Noticef("listening on %s", address)
	lis, err := tran.Listen(address)
	if err != nil {
		logger.Fatalln("fatal error: ", err)
	}
	for {
		conn, err := lis.Accept()
		if err != nil {
			continue
		}
		go ServeConn(conn, codec)
	}
}
