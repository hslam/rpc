 package rpc

import (
	"net"
	"net/http"
	"github.com/gorilla/websocket"
	"hslam.com/git/x/rpc/protocol"
	"hslam.com/git/x/rpc/log"
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
		log.Errorf("fatal error: %s", err)
		return nil,err
	}
	var httpServer http.Server
	httpServer.Addr = address

	l:= &WSListener{address:address,httpServer:httpServer,listener:lis,server:server,maxConnNum:DefaultMaxConnNum}
	return l,nil
}
func (l *WSListener)Serve() (error) {
	log.Allf( "%s", "Waiting for clients")
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
			log.Warnf("upgrade : %s", err)
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
			log.Infof("new client %s comming",conn.RemoteAddr())
			ServeWSConn(l.server,conn)
			log.Infof("client %s exiting",conn.RemoteAddr())
			connChange <- -1
		}()
	})
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
	readChan := make(chan []byte,1)
	writeChan := make(chan []byte,1)
	finishChan:= make(chan bool,2)
	stopReadConnChan := make(chan bool,1)
	stopWriteConnChan := make(chan bool,1)
	stopChan := make(chan bool,1)
	var wsConn =&protocol.WSConn{conn}
	go protocol.ReadConn(wsConn, readChan, stopReadConnChan,finishChan)
	go protocol.WriteConn(wsConn, writeChan, stopWriteConnChan,finishChan)
	if server.multiplexing{
		jobChan := make(chan bool,server.asyncMax)
		for {
			select {
			case data := <-readChan:
				go func(data []byte ,writeChan chan []byte) {
					jobChan<-true
					priority,id,body,err:=protocol.UnpackFrame(data)
					if err!=nil{
						return
					}
					_,res_bytes, _ := server.ServeRPC(body)
					if res_bytes!=nil{
						frameBytes:=protocol.PacketFrame(priority,id,res_bytes)
						writeChan <- frameBytes
					}
					<-jobChan
				}(data,writeChan)
			case stop := <-finishChan:
				if stop {
					stopReadConnChan<-true
					stopWriteConnChan<-true
					goto endfor
				}
			}
		}
	}else if server.async{
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
	return ErrConnExit
}

