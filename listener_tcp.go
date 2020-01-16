package rpc

import (
	"net"
)

type tcpListener struct {
	server      *Server
	address     string
	netListener net.Listener
	maxConnNum  int
	connNum     int
}

func listenTCP(address string, server *Server) (Listener, error) {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		logger.Errorf("fatal error: %s", err)
		return nil, err
	}
	listener := &tcpListener{address: address, netListener: lis, server: server, maxConnNum: DefaultMaxConnNum}
	return listener, nil
}
func (l *tcpListener) Serve() error {
	logger.Noticef("%s\n", "waiting for clients")
	workerChan := make(chan bool, l.maxConnNum)
	connChange := make(chan int)
	go func() {
		for c := range connChange {
			l.connNum += c
		}
	}()
	for {
		conn, err := l.netListener.Accept()
		if err != nil {
			logger.Warnf("Accept: %s\n", err)
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
			defer func() { logger.Infof("client %s exiting\n", conn.RemoteAddr()) }()
			logger.Infof("client %s comming\n", conn.RemoteAddr())
			l.server.ServeConn(conn)
		}()
	}
	return nil
}
func (l *tcpListener) Addr() string {
	return l.address
}
