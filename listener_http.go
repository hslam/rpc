package rpc

import (
	"hslam.com/git/x/rpc/log"
	"net/http"
	"net"
	"io"
)

type HTTPListener struct {
	server			*Server
	address			string
	maxConnNum		int
	connNum			int
	workerChan  		chan bool
	connChange 		chan int
}
func ListenHTTP(address string,server *Server) (Listener, error) {
	listener:=  &HTTPListener{address:address,server:server,maxConnNum:DefaultMaxConnNum*server.asyncMax}
	return listener,nil
}


func (l *HTTPListener)Serve() (error) {
	log.Allf( "%s\n", "Waiting for clients")
	l.workerChan = make(chan bool,l.maxConnNum)
	l.connChange = make(chan int)
	go func() {
		for conn_change := range l.connChange {
			l.connNum += conn_change
		}
	}()
	http.Handle(HttpPath,l)
	lis, err := net.Listen("tcp", l.address)
	if err != nil {
		log.Errorf("fatal error: %s", err)
		return err
	}
	http.Serve(lis, nil)
	return nil

}


func (l *HTTPListener)Addr() (string) {
	return l.address
}

func (l *HTTPListener) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != "CONNECT" {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusMethodNotAllowed)
		io.WriteString(w, "405 must CONNECT\n")
		return
	}
	conn, _, err := w.(http.Hijacker).Hijack()
	if err != nil {
		return
	}
	io.WriteString(conn, "HTTP/1.1 "+HttpConnected+"\n\n")
	func() {
		defer func() {
			if err := recover(); err != nil {
			}
			<-l.workerChan
		}()
		l.connChange <- 1
		defer func() {l.connChange <- -1}()
		defer func() {log.Infof("client %s exiting\n",conn.RemoteAddr())}()
		log.Infof("new client %s comming\n",conn.RemoteAddr())
		l.server.ServeConn(conn)
	}()
}
