package rpc

import (
	"hslam.com/mgit/Mort/rpc/log"
	"net/http"
	"io/ioutil"
	"io"
)

type HTTPListener struct {
	address			string
}
func ListenHTTP(address string) (Listener, error) {
	http.HandleFunc("/", IndexHandler)
	listener:=  &HTTPListener{address:address}
	return listener,nil
}

func (l *HTTPListener)Serve() (error) {
	log.Allf( "%s", "Waiting for clients")
	err:=http.ListenAndServe(l.address, nil)
	if err!=nil{
		log.Fatalf("fatal error: %s", err)
	}
	return nil
}
func (l *HTTPListener)Addr() (string) {
	return l.address
}
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	var RemoteAddr=r.RemoteAddr
	log.AllInfof("new client %s comming",RemoteAddr)
	if r.Method!="POST"{
		return
	}
	//log.All(r.Proto)
	if useWorkerPool{
		workerPool.Process(func(obj interface{}, args ...interface{}) interface{} {
			var w = obj.(http.ResponseWriter)
			var r = args[0].(*http.Request)
			data, _ := ioutil.ReadAll(r.Body)
			return ServeHTTP(w,data)
		},w,r)
	}else {
		data, _ := ioutil.ReadAll(r.Body)
		ServeHTTP(w,data)
	}
}

func ServeHTTP(w io.Writer,data []byte)error {
	_,res_bytes, _ := ServeRPC(data)
	if res_bytes!=nil{
		log.Tracef("res_bytes %s len %d",res_bytes,len(res_bytes))
		_,err:=w.Write(res_bytes)
		return err
	}else {
		return nil
	}
	return ErrConnExit
}