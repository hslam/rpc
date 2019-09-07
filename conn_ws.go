package rpc
import (
	"hslam.com/mgit/Mort/rpc/protocol"
	"github.com/gorilla/websocket"
	"hslam.com/mgit/Mort/rpc/log"
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

func (t *WSConn)Handle(readChan chan []byte,writeChan chan []byte, stopChan chan bool,finishChan chan bool){
	t.readChan=readChan
	t.writeChan=writeChan
	t.stopChan=stopChan
	t.finishChan=make(chan bool)
	t.handle()
}
func (t *WSConn)handle(){
	readChan:=make(chan []byte)
	writeChan:=make(chan []byte)
	finishChan:= make(chan bool)
	stopReadConnChan := make(chan bool,1)
	stopWriteConnChan := make(chan bool,1)
	go protocol.ReadConn(t.conn, readChan, stopReadConnChan,finishChan)
	go protocol.WriteConn(t.conn, writeChan, stopWriteConnChan,finishChan)
	go func() {
		for {
			select {
			case v:=<-readChan:
				t.readChan<-v
			case v:=<-t.writeChan:
				writeChan<-v
			case stop := <-finishChan:
				if stop {
					stopReadConnChan<-true
					stopWriteConnChan<-true
					t.finishChan<-true
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
	}()
}
func (t *WSConn)TickerFactor()(int){
	return 100
}
func (t *WSConn)BatchFactor()(int){
	return 64
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
