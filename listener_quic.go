package rpc

import (
	"context"
	"github.com/lucas-clemente/quic-go"
)

type quicListener struct {
	server       *Server
	address      string
	quicListener quic.Listener
	maxConnNum   int
	connNum      int
}

func listenQUIC(address string, server *Server) (Listener, error) {
	l, err := quic.ListenAddr(address, generateQuicTLSConfig(), nil)
	if err != nil {
		logger.Errorf("fatal error: %s", err)
		return nil, err
	}
	listener := &quicListener{address: address, quicListener: l, server: server, maxConnNum: DefaultMaxConnNum}
	return listener, nil
}
func (l *quicListener) Serve() error {
	logger.Noticef("%s\n", "waiting for clients")
	workerChan := make(chan bool, l.maxConnNum)
	connChange := make(chan int)
	go func() {
		for c := range connChange {
			l.connNum += c
		}
	}()
	for {
		sess, err := l.quicListener.Accept(context.Background())
		if err != nil {
			logger.Warnf("Accept: %s\n", err)
			continue
		} else {
			workerChan <- true
			go func() {
				defer func() {
					if err := recover(); err != nil {
					}
					<-workerChan
				}()
				connChange <- 1
				defer func() { connChange <- -1 }()
				defer func() { logger.Infof("client %s exiting\n", sess.RemoteAddr()) }()
				logger.Infof("client %s comming\n", sess.RemoteAddr())
				stream, err := sess.AcceptStream(context.Background())
				if err != nil {
					logger.Errorln(err)
					return
				}
				l.server.ServeConn(stream)
			}()
		}
	}
	return nil
}
func (l *quicListener) Addr() string {
	return l.address
}
