package rpc

import (
	"hslam.com/mgit/Mort/rpc/log"
	"net/http"
	"io/ioutil"
	"io"
)

type HTTPListener struct {
	server			*Server
	address			string
}
func ListenHTTP(address string,server *Server) (Listener, error) {
	listener:=  &HTTPListener{address:address,server:server}
	return listener,nil
}


func (l *HTTPListener)Serve() (error) {
	log.Allf( "%s", "Waiting for clients")
	handler:=new(Handler)
	handler.server=l.server
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
}
func (h *Handler)ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	var RemoteAddr=r.RemoteAddr
	log.AllInfof("new client %s comming",RemoteAddr)
	if r.Method!="POST"{
		return
	}
	if h.server.useWorkerPool{
		h.server.workerPool.Process(func(obj interface{}, args ...interface{}) interface{} {
			var w = obj.(http.ResponseWriter)
			var r = args[0].(*http.Request)
			var server = args[1].(*Server)
			data, _ := ioutil.ReadAll(r.Body)
			return ServeHTTP(server,w,data)
		},w,r,h.server)
	}else {
		data, _ := ioutil.ReadAll(r.Body)
		ServeHTTP(h.server,w,data)
	}
}

func ServeHTTP(server *Server,w io.Writer,data []byte)error {
	_,res_bytes, _ := server.ServeRPC(data)
	if res_bytes!=nil{
		_,err:=w.Write(res_bytes)
		return err
	}else {
		return nil
	}
	return ErrConnExit
}