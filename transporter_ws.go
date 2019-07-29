package rpc
import (
	"hslam.com/mgit/Mort/rpc/protocol"
	"github.com/gorilla/websocket"
	"net/url"
)

type WSTransporter struct {
	wsConn			*protocol.WSConn
	address			string
	CanWork			bool
}

func NewWSTransporter(address string)  (Transporter, error)  {
	u := url.URL{Scheme: "ws", Host: address, Path: "/"}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		Fatalf("fatal error: %s", err)
	}
	t:=&WSTransporter{
		wsConn:&protocol.WSConn{c},
		address:address,
	}
	t.CanWork=true

	return t, nil
}

func (t *WSTransporter)Handle(readChan chan []byte,writeChan chan []byte, stopChan chan bool){
	go protocol.ReadConn(t.wsConn, readChan, stopChan)
	go protocol.WriteConn(t.wsConn, writeChan,stopChan)
}
func (t *WSTransporter)TickerFactor()(int){
	return 100
}
func (t *WSTransporter)BatchFactor()(int){
	return 64
}
func (t *WSTransporter)Close()(error){
	return 	t.wsConn.Close()
}
