package rpc
import (
	"hslam.com/git/x/rpc/protocol"
	"github.com/gorilla/websocket"
	"hslam.com/git/x/rpc/log"
	"net/url"
)

type WSConn struct {
	conn			*protocol.WSConn
	address			string
	CanWork			bool
	readChan 		chan []byte
	writeChan 		chan []byte
	stopChan 		chan bool
	finishChan 		chan bool
	closed			bool
}

func DialWS(address string)  (Conn, error)  {
	u := url.URL{Scheme: "ws", Host: address, Path: "/"}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Errorf("fatal error: %s", err)
		return nil,err
	}
	t:=&WSConn{
		conn:&protocol.WSConn{c},
		address:address,
	}
	t.CanWork=true

	return t, nil
}
func (t *WSConn)Buffer(enable bool){
}
func (t *WSConn)Handle(readChan chan []byte,writeChan chan []byte, stopChan chan bool,finishChan chan bool){
	t.readChan=readChan
	t.writeChan=writeChan
	t.stopChan=stopChan
	t.finishChan=finishChan
	t.handle()
}
func (t *WSConn)handle(){
	readChan:=make(chan []byte,1)
	writeChan:=make(chan []byte,1)
	finishChan:= make(chan bool,2)
	stopReadConnChan := make(chan bool,1)
	stopWriteConnChan := make(chan bool,1)
	go protocol.ReadConn(t.conn, readChan, stopReadConnChan,finishChan)
	go protocol.WriteConn(t.conn, writeChan, stopWriteConnChan,finishChan)
	go func() {
		t.closed=false
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
					stopReadConnChan<-true
					stopWriteConnChan<-true
					func(){
						defer func() {if err := recover(); err != nil {}}()
						t.finishChan<-true
					}()
					goto endfor
				}
			case <-t.stopChan:
				stopReadConnChan<-true
				stopWriteConnChan<-true
				goto endfor
			}
		}
	endfor:
		close(readChan)
		close(writeChan)
		close(finishChan)
		close(stopReadConnChan)
		close(stopWriteConnChan)
		t.closed=true
	}()
}
func (t *WSConn)TickerFactor()(int){
	return 100
}
func (t *WSConn)BatchFactor()(int){
	return 512
}
func (t *WSConn)Retry()(error){
	u := url.URL{Scheme: "ws", Host: t.address, Path: "/"}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)

	if err != nil {
		log.Errorf("fatal error: %s", err)
		return err
	}
	t.conn=&protocol.WSConn{c}
	t.handle()
	return nil
}
func (t *WSConn)Close()(error){
	return 	t.conn.Close()
}
func (t *WSConn)Closed()(bool){
	return t.closed
}