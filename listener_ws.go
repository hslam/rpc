 package rpc

import (
	"net"
	"net/http"
	"github.com/gorilla/websocket"
	"hslam.com/git/x/protocol"
)
var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024 * 64,
		WriteBufferSize: 1024 * 64,
	}
)

type WSListener struct {
	server			*Server
	address			string
	httpServer		http.Server
	listener		net.Listener
	maxConnNum		int
	connNum			int
}
func ListenWS(address string,server *Server) (Listener, error) {
	lis, err := net.Listen("tcp", address)
	if err!=nil{
		Errorf("fatal error: %s", err)
		return nil,err
	}
	var httpServer http.Server
	httpServer.Addr = address

	l:= &WSListener{address:address,httpServer:httpServer,listener:lis,server:server,maxConnNum:DefaultMaxConnNum}
	return l,nil
}
func (l *WSListener)Serve() (error) {
	Allf( "%s\n", "waiting for clients")
	workerChan := make(chan bool,l.maxConnNum)
	connChange := make(chan int)
	go func() {
		for conn_change := range connChange {
			l.connNum += conn_change
		}
	}()
	http.HandleFunc("/",func (w http.ResponseWriter, r *http.Request) {
		r.Header.Del("Origin")
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			Warnf("upgrade : %s\n", err)
			return
		}
		workerChan<-true
		func() {
			defer func() {
				if err := recover(); err != nil {
				}
				<-workerChan
			}()
			connChange <- 1
			defer func() {connChange <- -1}()
			defer func() {Infof("client %s exiting\n",conn.RemoteAddr())}()
			Infof("client %s comming\n",conn.RemoteAddr())
			l.server.ServeMessage(&protocol.WSConn{Conn:conn})
		}()
	})
	err:=l.httpServer.Serve(l.listener)
	if err != nil {
		Errorf("fatal error: %s", err)
		return err
	}
	return nil
}
 func (l *WSListener)Addr() (string) {
	 return l.address
 }
