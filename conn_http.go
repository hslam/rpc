package rpc
import (
	"net/http"
	"hslam.com/git/x/rpc/protocol"
	"io"
	"bufio"
	"net"
	"errors"
)
type HTTPConn struct {
	conn 			net.Conn
	address			string
	CanWork			bool
	closed			bool
	buffer 			bool
	readChan 		chan []byte
	writeChan 		chan []byte
	stopChan 		chan bool
	finishChan		chan bool
}

func DialHTTP(address string)  (Conn, error)  {
	var err error
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	io.WriteString(conn, "CONNECT "+HttpPath+" HTTP/1.1\n\n")

	// Require successful HTTP response
	// before switching to RPC protocol.
	resp, err := http.ReadResponse(bufio.NewReader(conn), &http.Request{Method: "CONNECT"})
	if err != nil || resp.Status != HttpConnected {
		if err == nil {
			err = errors.New("unexpected HTTP response: " + resp.Status)
		}
		conn.Close()
		return nil, &net.OpError{
			Op:   "dial-http",
			Net:  "tcp" + " " + address,
			Addr: nil,
			Err:  err,
		}
	}
	t:=&HTTPConn{
		address:address,
		conn:conn,
	}
	return t, nil
}
func (t *HTTPConn)Buffer(enable bool){
	t.buffer=enable
}
func (t *HTTPConn)Multiplexing(enable bool){
}
func (t *HTTPConn)Handle(readChan chan []byte,writeChan chan []byte, stopChan chan bool,finishChan chan bool){
	t.readChan=readChan
	t.writeChan=writeChan
	t.stopChan=stopChan
	t.finishChan=finishChan
	t.handle()
}
func (t *HTTPConn)handle(){
	readChan:=make(chan []byte,1)
	writeChan:=make(chan []byte,1)
	finishChan:= make(chan bool,2)
	stopReadStreamChan := make(chan bool,1)
	stopWriteStreamChan := make(chan bool,1)
	go protocol.ReadStream(t.conn, readChan, stopReadStreamChan,finishChan)
	go protocol.WriteStream(t.conn, writeChan, stopWriteStreamChan,finishChan,t.buffer)
	go func() {
		t.closed=false
		//log.Traceln("TCPConn.handle start")
		for {
			select {
			case v:=<-readChan:
				func() {
					defer func() {
						if err := recover(); err != nil {
						}
					}()
					t.readChan<-v
				}()
			case v:=<-t.writeChan:
				func() {
					defer func() {
						if err := recover(); err != nil {
						}
					}()
					writeChan<-v
				}()
			case stop := <-finishChan:
				if stop {
					stopReadStreamChan<-true
					stopWriteStreamChan<-true
					func(){
						defer func() {if err := recover(); err != nil {}}()
						t.finishChan<-true
					}()
					goto endfor
				}
			case <-t.stopChan:
				stopReadStreamChan<-true
				stopWriteStreamChan<-true
				goto endfor
			}
		}
	endfor:
		close(readChan)
		close(writeChan)
		close(finishChan)
		close(stopReadStreamChan)
		close(stopWriteStreamChan)
		//log.Traceln("TCPConn.handle end")
		t.closed=true
	}()
}
func (t *HTTPConn)TickerFactor()(int){
	return 300
}
func (t *HTTPConn)BatchFactor()(int){
	return 512
}
func (t *HTTPConn)Retry()(error){
	var err error
	conn, err := net.Dial("tcp", t.address)
	if err != nil {
		return  err
	}
	path:=""
	io.WriteString(conn, "CONNECT "+path+" HTTP/1.0\n\n")

	// Require successful HTTP response
	// before switching to RPC protocol.
	resp, err := http.ReadResponse(bufio.NewReader(conn), &http.Request{Method: "CONNECT"})
	if err != nil || resp.Status != HttpConnected {
		if err == nil {
			err = errors.New("unexpected HTTP response: " + resp.Status)
		}
		conn.Close()
		return  &net.OpError{
			Op:   "dial-http",
			Net:  "tcp" + " " + t.address,
			Addr: nil,
			Err:  err,
		}
	}
	t.conn=conn
	return nil
}
func (t *HTTPConn)Close()(error){
	return nil
}

func (t *HTTPConn)Closed()(bool){
	return t.closed
}