package rpc

import (
	"github.com/hslam/protocol"
	"io"
	"io/ioutil"
	"net/http"
)

type httpListener struct {
	server     *Server
	httpServer *http.Server
	address    string
	maxConnNum int
	connNum    int
}

func listenHTTP(address string, server *Server) (Listener, error) {
	httpServer := &http.Server{
		Addr: address,
	}
	listener := &httpListener{address: address, server: server, httpServer: httpServer, maxConnNum: DefaultMaxConnNum * server.asyncMax}
	return listener, nil
}

func (l *httpListener) Serve() error {
	logger.Noticef("%s\n", "waiting for clients")
	h := new(handler)
	h.server = l.server
	h.workerChan = make(chan bool, l.maxConnNum)
	h.connChange = make(chan int)
	go func() {
		for c := range h.connChange {
			l.connNum += c
		}
	}()
	l.httpServer.Handler = h
	err := l.httpServer.ListenAndServe()
	if err != nil {
		if stringsContains(err.Error(), "http: Server closed") {
			return nil
		}
		logger.Errorf("fatal error: %s", err)
		return err
	}
	return nil
}

func (l *httpListener) Addr() string {
	return l.address
}
func (l *httpListener) Close() error {
	return l.httpServer.Close()
}

type handler struct {
	server     *Server
	workerChan chan bool
	connChange chan int
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	var RemoteAddr = r.RemoteAddr
	if r.Method == "POST" {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return
		}
		h.workerChan <- true
		func() {
			defer func() {
				if err := recover(); err != nil {
				}
				<-h.workerChan
			}()
			h.connChange <- 1
			defer func() { h.connChange <- -1 }()
			defer func() { logger.Noticef("client %s exiting\n", RemoteAddr) }()
			logger.Noticef("client %s comming\n", RemoteAddr)
			h.Serve(w, data)
		}()
	} else if r.Method == "CONNECT" {
		conn, _, err := w.(http.Hijacker).Hijack()
		if err != nil {
			return
		}
		io.WriteString(conn, "HTTP/1.1 "+HTTPConnected+"\n\n")
		func() {
			defer func() {
				if err := recover(); err != nil {
				}
				<-h.workerChan
			}()
			h.connChange <- 1
			defer func() { h.connChange <- -1 }()
			defer func() { logger.Infof("client %s exiting\n", conn.RemoteAddr()) }()
			logger.Infof("client %s comming\n", conn.RemoteAddr())
			h.server.ServeConn(conn)
		}()
	} else {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusMethodNotAllowed)
		io.WriteString(w, "405 must CONNECT or POST\n")
	}
}

func (h *handler) Serve(w io.Writer, data []byte) error {
	if h.server.multiplexing {
		priority, id, body, err := protocol.UnpackFrame(data)
		if err != nil {
			return err
		}
		_, resBytes := h.server.Serve(body)
		if resBytes != nil {
			frameBytes := protocol.PacketFrame(priority, id, resBytes)
			_, err := w.Write(frameBytes)
			return err
		}
	} else {
		_, resBytes := h.server.Serve(data)
		if resBytes != nil {
			_, err := w.Write(resBytes)
			return err
		}
		return nil
	}
	return ErrConnExit
}
