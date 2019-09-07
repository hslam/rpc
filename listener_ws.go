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
	server			*Server
	address			string
	httpServer		http.Server
	listener		net.Listener
}
func ListenWS(address string,server *Server) (Listener, error) {
	lis, err := net.Listen("tcp", address)
	if err!=nil{
		log.Errorf("fatal error: %s", err)
		return nil,err
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
		if server.useWorkerPool{
			server.workerPool.Process(func(obj interface{}, args ...interface{}) interface{} {
				var w = obj.( *websocket.Conn)
				var server = args[0].(*Server)
				return ServeWSConn(server,w)
			},conn,server)
		}else {
			ServeWSConn(server,conn)
		}
	})
	l:= &WSListener{address:address,httpServer:httpServer,listener:lis,server:server}

	return l,nil
}
func (l *WSListener)Serve() (error) {
	log.Allf( "%s", "Waiting for clients")
	err:=l.httpServer.Serve(l.listener)
	if err != nil {
		log.Errorf("fatal error: %s", err)
		return err
	}
	return nil
}
 func (l *WSListener)Addr() (string) {
	 return l.address
 }
func ServeWSConn(server *Server,conn *websocket.Conn)error {
	var RemoteAddr=conn.RemoteAddr().String()
	readChan := make(chan []byte)
	writeChan := make(chan []byte)
	finishChan:= make(chan bool)
	stopReadConnChan := make(chan bool,1)
	stopWriteConnChan := make(chan bool,1)
	stopChan := make(chan bool,1)
	var wsConn =&protocol.WSConn{conn}
	go protocol.ReadConn(wsConn, readChan, stopReadConnChan,finishChan)
	go protocol.WriteConn(wsConn, writeChan, stopWriteConnChan,finishChan)
	if server.async{
		syncConn:=newSyncConn(server)
		go protocol.HandleSyncConn(syncConn, writeChan,readChan,stopChan,server.asyncMax)
		select {
		case stop := <-finishChan:
			if stop {
				stopReadConnChan<-true
				stopWriteConnChan<-true
				stopChan<-true
				goto endfor
			}
		}
	}else {
		for {
			select {
			case data := <-readChan:
				_,res_bytes, _ := server.ServeRPC(data)
				if res_bytes!=nil{
					writeChan <- res_bytes
				}
			case stop := <-finishChan:
				if stop {
					stopReadConnChan<-true
					stopWriteConnChan<-true
					goto endfor
				}
			}
		}
	}

endfor:
	defer conn.Close()
	close(writeChan)
	close(readChan)
	close(finishChan)
	close(stopReadConnChan)
	close(stopWriteConnChan)
	close(stopChan)
	log.Infof("client %s exiting",RemoteAddr)
	return ErrConnExit
}

