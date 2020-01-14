package rpc

import (
	"net"
	"os"
)

type IPCListener struct {
	server      *Server
	address     string
	netListener net.Listener
	maxConnNum  int
	connNum     int
}

func ListenIPC(address string, server *Server) (Listener, error) {
	os.Remove(address)
	var addr *net.UnixAddr
	var err error
	if addr, err = net.ResolveUnixAddr("unix", address); err != nil {
		return nil, err
	}
	lis, err := net.ListenUnix("unix", addr)
	if err != nil {
		Errorf("fatal error: %s", err)
		return nil, err
	}
	listener := &IPCListener{address: address, netListener: lis, server: server, maxConnNum: DefaultMaxConnNum}
	return listener, nil
}
func (l *IPCListener) Serve() error {
	Allf("%s\n", "waiting for clients")
	workerChan := make(chan bool, l.maxConnNum)
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
		workerChan <- true
		go func() {
			defer func() {
				if err := recover(); err != nil {
				}
				<-workerChan
			}()
			connChange <- 1
			defer func() { connChange <- -1 }()
			defer func() { Infof("client %s exiting\n", conn.RemoteAddr()) }()
			Infof("client %s comming\n", conn.RemoteAddr())
			l.server.ServeConn(conn)
		}()
	}
	return nil
}
func (l *IPCListener) Addr() string {
	return l.address
}
