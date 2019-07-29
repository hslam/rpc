package rpc

import (
	"net"
	"sync"
	"hslam.com/mgit/Mort/rpc/protocol"
)

type TCPListener struct {
	reqMutex 		sync.Mutex
	address			string
	netListener		net.Listener
}
func ListenTCP(address string) (Listener, error) {
	lis, err := net.Listen("tcp", address)
	if err!=nil{
		Fatalf("fatal error: %s", err)
	}
	listener:= &TCPListener{address:address,netListener:lis}
	return listener,nil
}
func (l *TCPListener)Serve() (error) {
	Allf( "%s", "Waiting for clients")
	for {
		conn, err := l.netListener.Accept()
		if err != nil {
			Warnf("Accept: %s", err)
			continue
		}
		Infof("new client %s comming",conn.RemoteAddr())
		if useWorkerPool{
			workerPool.ProcessAsyn( func(obj interface{}, args ...interface{}) interface{} {
				var c = obj.(net.Conn )
				return ServeTCPConn(c)
			},conn)
		}else {
			go ServeTCPConn(conn)
		}
	}
	return nil
}
func (l *TCPListener)Addr() (string) {
	return l.address
}
func ServeTCPConn(conn net.Conn)error {
	var RemoteAddr=conn.RemoteAddr().String()
	readChan := make(chan []byte)
	writeChan := make(chan []byte)
	stopChan := make(chan bool)
	go protocol.ReadStream(conn, readChan, stopChan)
	go protocol.WriteStream(conn, writeChan, stopChan)
	for {
		select {
		case data := <-readChan:
			_,res_bytes, _ := ServeRPC(data)
			if res_bytes!=nil{
				writeChan <- res_bytes
			}
		case stop := <-stopChan:
			if stop {
				goto endfor
			}
		}
	}
	endfor:
		defer conn.Close()
		close(writeChan)
		close(readChan)
		close(stopChan)
	Infof("client %s exiting",RemoteAddr)
	return ErrConnExit
}