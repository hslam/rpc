package rpc

import (
	"github.com/gorilla/websocket"
	"github.com/hslam/protocol"
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
	httpServer *http.Server
	maxConnNum int
	connNum    int
}

func listenWS(address string, server *Server) (Listener, error) {
	//lis, err := net.Listen("tcp", address)
	//if err != nil {
	//	logger.Errorf("fatal error: %s", err)
	//	return nil, err
	//}
	httpServer := &http.Server{
		Addr: address,
	}
	l := &wsListener{address: address, httpServer: httpServer, server: server, maxConnNum: DefaultMaxConnNum}
	return l, nil
}
func (l *wsListener) Serve() error {
	logger.Noticef("%s\n", "waiting for clients")
	h := new(wsHandler)
	h.server = l.server
	h.workerChan = make(chan bool, l.maxConnNum)
	h.connChange = make(chan int)
	go func() {
		for c := range h.connChange {
			l.connNum += c
		}
	}()
	l.httpServer.Handler = h
	err := l.httpServer.ListenAndServe()
	if err != nil {
		if stringsContains(err.Error(), "http: Server closed") {
			return nil
		}
		logger.Errorf("fatal error: %s", err)
		return err
	}
	return nil
}
func (l *wsListener) Addr() string {
	return l.address
}
func (l *wsListener) Close() error {
	return l.httpServer.Close()
}

type wsHandler struct {
	server     *Server
	workerChan chan bool
	connChange chan int
}

func (h *wsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.Header.Del("Origin")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Warnf("upgrade : %s\n", err)
		return
	}
	h.workerChan <- true
	func() {
		defer func() {
			if err := recover(); err != nil {
			}
			<-h.workerChan
		}()
		h.connChange <- 1
		defer func() { h.connChange <- -1 }()
		defer func() { logger.Infof("client %s exiting\n", conn.RemoteAddr()) }()
		logger.Infof("client %s comming\n", conn.RemoteAddr())
		h.server.ServeMessage(&protocol.MsgConn{MessageConn: &websocketConn{conn}})
	}()
}
