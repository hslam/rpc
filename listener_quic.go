package rpc

import (
	"github.com/lucas-clemente/quic-go"
	"hslam.com/git/x/rpc/log"
)

type QUICListener struct {
	server			*Server
	address			string
	quicListener	quic.Listener
	maxConnNum		int
	connNum			int
}
func ListenQUIC(address string,server *Server) (Listener, error) {
	quic_listener, err := quic.ListenAddr(address, generateTLSConfig(), nil)
	if err!=nil{
		log.Errorf("fatal error: %s", err)
		return nil,err
	}
	listener:= &QUICListener{address:address,quicListener:quic_listener,server:server,maxConnNum:DefaultMaxConnNum}
	return listener,nil
}
func (l *QUICListener)Serve() (error) {
	log.Allf( "%s\n", "Waiting for clients")
	workerChan := make(chan bool,l.maxConnNum)
	connChange := make(chan int)
	go func() {
		for conn_change := range connChange {
			l.connNum += conn_change
		}
	}()
	for{
		sess, err := l.quicListener.Accept()
		if err != nil {
			log.Warnf("Accept: %s\n", err)
			continue
		}else{
			workerChan<-true
			go func() {
				defer func() {
					if err := recover(); err != nil {
					}
					<-workerChan
				}()
				connChange <- 1
				defer func() {connChange <- -1}()
				defer func() {log.Infof("client %s exiting\n",sess.RemoteAddr())}()
				log.Infof("new client %s comming\n",sess.RemoteAddr())
				stream, err := sess.AcceptStream()
				if err != nil {
					log.Errorln(err)
					return
				}
				l.server.ServeConn(stream)
			}()
		}
	}
	return nil
}
func (l *QUICListener)Addr() (string) {
	return l.address
}
