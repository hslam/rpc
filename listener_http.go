package rpc

import (
	"hslam.com/git/x/protocol"
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
	Allf( "%s\n", "waiting for clients")
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
		Errorf("fatal error: %s", err)
		return err
	}
	return nil
}


func (l *HTTPListener)Addr() (string) {
	return l.address
}


type Handler struct {
	server			*Server
	workerChan  	chan bool
	connChange 		chan int
}
func (h *Handler)ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	var RemoteAddr=r.RemoteAddr
	if r.Method=="POST"{
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
			defer func() {AllInfof("client %s exiting\n",RemoteAddr)}()
			AllInfof("client %s comming\n",RemoteAddr)
			h.Serve(w,data)
		}()
	}else if r.Method=="CONNECT"{
		conn, _, err := w.(http.Hijacker).Hijack()
		if err != nil {
			return
		}
		io.WriteString(conn, "HTTP/1.1 "+HttpConnected+"\n\n")
		func() {
			defer func() {
				if err := recover(); err != nil {
				}
				<-h.workerChan
			}()
			h.connChange <- 1
			defer func() {h.connChange <- -1}()
			defer func() {Infof("client %s exiting\n",conn.RemoteAddr())}()
			Infof("client %s comming\n",conn.RemoteAddr())
			h.server.ServeConn(conn)
		}()
	}else {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusMethodNotAllowed)
		io.WriteString(w, "405 must CONNECT or POST\n")
	}
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