package rpc

import (
	"hslam.com/git/x/rpc/log"
	"hslam.com/git/x/rpc/protocol"
	"net/http"
	"io/ioutil"
	"io"
)

type HTTPListener struct {
	server			*Server
	address			string
	maxConnNum		int
	connNum			int
}
func ListenHTTP(address string,server *Server) (Listener, error) {
	listener:=  &HTTPListener{address:address,server:server,maxConnNum:DefaultMaxConnNum*server.asyncMax}
	return listener,nil
}


func (l *HTTPListener)Serve() (error) {
	log.Allf( "%s", "Waiting for clients")
	handler:=new(Handler)
	handler.server=l.server
	handler.workerChan = make(chan bool,l.maxConnNum)
	handler.connChange = make(chan int)
	go func() {
		for conn_change := range handler.connChange {
			l.connNum += conn_change
		}
	}()
	err:=http.ListenAndServe(l.address, handler)

	if err!=nil{
		log.Errorf("fatal error: %s", err)
		return err
	}
	return nil
}


func (l *HTTPListener)Addr() (string) {
	return l.address
}


type Handler struct {
	server			*Server
	workerChan  		chan bool
	connChange 		chan int
}
func (h *Handler)ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	var RemoteAddr=r.RemoteAddr
	if r.Method!="POST"{
		return
	}
	data, err := ioutil.ReadAll(r.Body)
	if err!=nil{
		return
	}
	h.workerChan<-true
	func() {
		defer func() {
			if err := recover(); err != nil {
			}
			<-h.workerChan
		}()
		h.connChange <- 1
		log.AllInfof("new client %s comming",RemoteAddr)
		ServeHTTP(h.server,w,data)
		log.Infof("client %s exiting",RemoteAddr)
		h.connChange <- -1
	}()
}

func ServeHTTP(server *Server,w io.Writer,data []byte)error {
	if server.multiplexing{
		priority,id,body,err:=protocol.UnpackFrame(data)
		if err!=nil{
			return err
		}
		_,res_bytes, _ := server.ServeRPC(body)
		if res_bytes!=nil{
			frameBytes:=protocol.PacketFrame(priority,id,res_bytes)
			_,err:=w.Write(frameBytes)
			return err
		}
	}else {
		_,res_bytes, _ := server.ServeRPC(data)
		if res_bytes!=nil{
			_,err:=w.Write(res_bytes)
			return err
		}else {
			return nil
		}
	}
	return ErrConnExit
}