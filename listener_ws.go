package rpc

import (
	"github.com/gorilla/websocket"
	"github.com/hslam/protocol"
	"net"
	"net/http"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024 * 64,
		WriteBufferSize: 1024 * 64,
	}
)

type wsListener struct {
	server     *Server
	address    string
	httpServer http.Server
	listener   net.Listener
	maxConnNum int
	connNum    int
}

func listenWS(address string, server *Server) (Listener, error) {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		logger.Errorf("fatal error: %s", err)
		return nil, err
	}
	var httpServer http.Server
	httpServer.Addr = address

	l := &wsListener{address: address, httpServer: httpServer, listener: lis, server: server, maxConnNum: DefaultMaxConnNum}
	return l, nil
}
func (l *wsListener) Serve() error {
	logger.Noticef("%s\n", "waiting for clients")
	workerChan := make(chan bool, l.maxConnNum)
	connChange := make(chan int)
	go func() {
		for c := range connChange {
			l.connNum += c
		}
	}()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		r.Header.Del("Origin")
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			logger.Warnf("upgrade : %s\n", err)
			return
		}
		workerChan <- true
		func() {
			defer func() {
				if err := recover(); err != nil {
				}
				<-workerChan
			}()
			connChange <- 1
			defer func() { connChange <- -1 }()
			defer func() { logger.Infof("client %s exiting\n", conn.RemoteAddr()) }()
			logger.Infof("client %s comming\n", conn.RemoteAddr())
			l.server.ServeMessage(&protocol.MsgConn{MessageConn: &websocketConn{conn}})
		}()
	})
	err := l.httpServer.Serve(l.listener)
	if err != nil {
		logger.Errorf("fatal error: %s", err)
		return err
	}
	return nil
}
func (l *wsListener) Addr() string {
	return l.address
}
