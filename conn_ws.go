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
}

func DialWS(address string)  (Conn, error)  {
	u := url.URL{Scheme: "ws", Host: address, Path: "/"}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}
	t:=&WSConn{
		conn:&protocol.WSConn{c},
		address:address,
	}
	t.CanWork=true

	return t, nil
}

func (t *WSConn)Handle(readChan chan []byte,writeChan chan []byte, stopChan chan bool){
	t.readChan=readChan
	t.writeChan=writeChan
	t.stopChan=stopChan
	t.handle()
}
func (t *WSConn)handle(){
	readChan:=make(chan []byte)
	writeChan:=make(chan []byte)
	stopChan:=make(chan bool)
	go protocol.ReadConn(t.conn, readChan, stopChan)
	go protocol.WriteConn(t.conn, writeChan, stopChan)
	go func() {
		for {
			select {
			case v:=<-readChan:
				t.readChan<-v
			case v:=<-t.writeChan:
				writeChan<-v
			case <-stopChan:
				t.stopChan<-true
				close(readChan)
				close(writeChan)
				close(stopChan)
				goto endfor
			}
		}
	endfor:
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
