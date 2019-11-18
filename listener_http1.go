package rpc

import (
	"hslam.com/git/x/rpc/log"
	"hslam.com/git/x/rpc/protocol"
	"net/http"
	"io/ioutil"
	"io"
)

type HTTP1Listener struct {
	server			*Server
	address			string
	maxConnNum		int
	connNum			int
}
func ListenHTTP1(address string,server *Server) (Listener, error) {
	listener:=  &HTTP1Listener{address:address,server:server,maxConnNum:DefaultMaxConnNum*server.asyncMax}
	return listener,nil
}


func (l *HTTP1Listener)Serve() (error) {
	log.Allf( "%s\n", "Waiting for clients")
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


func (l *HTTP1Listener)Addr() (string) {
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
		defer func() {h.connChange <- -1}()
		defer func() {log.AllInfof("client %s exiting\n",RemoteAddr)}()
		log.AllInfof("new client %s comming\n",RemoteAddr)
		h.Serve(w,data)
	}()
}

func (h *Handler)Serve(w io.Writer,data []byte)error {
	if h.server.multiplexing{
		priority,id,body,err:=protocol.UnpackFrame(data)
		if err!=nil{
			return err
		}
		_,res_bytes, _ := h.server.Serve(body)
		if res_bytes!=nil{
			frameBytes:=protocol.PacketFrame(priority,id,res_bytes)
			_,err:=w.Write(frameBytes)
			return err
		}
	}else {
		_,res_bytes, _ := h.server.Serve(data)
		if res_bytes!=nil{
			_,err:=w.Write(res_bytes)
			return err
		}else {
			return nil
		}
	}
	return ErrConnExit
}