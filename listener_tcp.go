package rpc

import (
	"net"
)

type TCPListener struct {
	server			*Server
	address			string
	netListener		net.Listener
	maxConnNum		int
	connNum			int
}
func ListenTCP(address string,server *Server) (Listener, error) {
	lis, err := net.Listen("tcp", address)
	if err!=nil{
		Errorf("fatal error: %s", err)
		return nil,err
	}
	listener:= &TCPListener{address:address,netListener:lis,server:server,maxConnNum:DefaultMaxConnNum}
	return listener,nil
}
func (l *TCPListener)Serve() (error) {
	Allf( "%s\n", "waiting for clients")
	workerChan := make(chan bool,l.maxConnNum)
	connChange := make(chan int)
	go func() {
		for conn_change := range connChange {
			l.connNum += conn_change
		}
	}()
	for {
		conn, err := l.netListener.Accept()
		if err != nil {
			Warnf("Accept: %s\n", err)
			continue
		}
		workerChan<-true
		go func() {
			defer func() {
				if err := recover(); err != nil {
				}
				<-workerChan
			}()
			connChange <- 1
			defer func() {connChange <- -1}()
			defer func() {Infof("client %s exiting\n",conn.RemoteAddr())}()
			Infof("client %s comming\n",conn.RemoteAddr())
			l.server.ServeConn(conn)
		}()
	}
	return nil
}
func (l *TCPListener)Addr() (string) {
	return l.address
}
