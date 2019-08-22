 package rpc

import (
	"net"
	"net/http"
	"github.com/gorilla/websocket"
	"hslam.com/mgit/Mort/rpc/protocol"
	"hslam.com/mgit/Mort/rpc/log"
)
var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024 * 64,
		WriteBufferSize: 1024 * 64,
	}
)

type WSListener struct {
	address			string
	httpServer		http.Server
	listener		net.Listener
}
func ListenWS(address string) (Listener, error) {
	lis, err := net.Listen("tcp", address)
	if err!=nil{
		log.Fatalf("fatal error: %s", err)
	}
	var httpServer http.Server
	httpServer.Addr = address
	http.HandleFunc("/",func (w http.ResponseWriter, r *http.Request) {
		r.Header.Del("Origin")
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Warnf("upgrade : %s", err)
			return
		}
		log.Infof("new client %s comming",conn.RemoteAddr())
		if useWorkerPool{
			workerPool.Process(func(obj interface{}, args ...interface{}) interface{} {
				var w = obj.( *websocket.Conn)
				return ServeWSConn(w)
			},conn)
		}else {
			ServeWSConn(conn)
		}
	})
	l:= &WSListener{address:address,httpServer:httpServer,listener:lis}

	return l,nil
}
func (l *WSListener)Serve() (error) {
	log.Allf( "%s", "Waiting for clients")
	err:=l.httpServer.Serve(l.listener)
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}
	return nil
}
 func (l *WSListener)Addr() (string) {
	 return l.address
 }
func ServeWSConn(conn *websocket.Conn)error {
	var RemoteAddr=conn.RemoteAddr().String()
	readChan := make(chan []byte)
	writeChan := make(chan []byte)
	stopChan := make(chan bool)
	var wsConn =&protocol.WSConn{conn}
	go protocol.ReadConn(wsConn, readChan, stopChan)
	go protocol.WriteConn(wsConn, writeChan, stopChan)
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
	log.Infof("client %s exiting",RemoteAddr)
	return ErrConnExit
}

