package rpc

import (
	"net"
	"hslam.com/mgit/Mort/rpc/protocol"
	"hslam.com/mgit/Mort/rpc/log"
)

type TCPListener struct {
	server			*Server
	address			string
	netListener		net.Listener
}
func ListenTCP(address string,server *Server) (Listener, error) {
	lis, err := net.Listen("tcp", address)
	if err!=nil{
		log.Errorf("fatal error: %s", err)
		return nil,err
	}
	listener:= &TCPListener{address:address,netListener:lis,server:server}
	return listener,nil
}
func (l *TCPListener)Serve() (error) {
	log.Allf( "%s", "Waiting for clients")
	for {
		conn, err := l.netListener.Accept()
		if err != nil {
			log.Warnf("Accept: %s", err)
			continue
		}
		log.Infof("new client %s comming",conn.RemoteAddr())
		if l.server.useWorkerPool{
			l.server.workerPool.ProcessAsyn( func(obj interface{}, args ...interface{}) interface{} {
				var c = obj.(net.Conn )
				var server = args[0].(*Server)
				return ServeTCPConn(server,c)
			},conn,l.server)
		}else {
			go ServeTCPConn(l.server,conn)
		}
	}
	return nil
}
func (l *TCPListener)Addr() (string) {
	return l.address
}
func ServeTCPConn(server *Server,conn net.Conn)error {
	var RemoteAddr=conn.RemoteAddr().String()
	readChan := make(chan []byte)
	writeChan := make(chan []byte)
	stopChan := make(chan bool)
	go protocol.ReadStream(conn, readChan, stopChan)
	go protocol.WriteStream(conn, writeChan, stopChan)
	for {
		select {
		case data := <-readChan:
			_,res_bytes, _ := server.ServeRPC(data)
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
	log.Infof("client %s exiting",RemoteAddr)
	return ErrConnExit
}