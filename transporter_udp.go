package rpc
import (
	"hslam.com/mgit/Mort/rpc/protocol"
	"net"
)

type UDPTransporter struct {
	conn 			net.Conn
	address			string
	CanWork			bool
}

func NewUDPTransporter(address string)  (Transporter, error)  {
	conn, err := net.Dial(UDP, address)
	if err != nil {
		Fatalf("fatal error: %s", err)
	}
	t:=&UDPTransporter{
		conn:conn,
		address:address,
	}
	return t, nil
}

func (t *UDPTransporter)Handle(readChan chan []byte,writeChan chan []byte, stopChan chan bool){
	go protocol.HandleMessage(t.conn,readChan,writeChan,stopChan)
}
func (t *UDPTransporter)TickerFactor()(int){
	return 1
}
func (t *UDPTransporter)BatchFactor()(int){
	return 1
}
func (t *UDPTransporter)Close()(error){
	return t.conn.Close()
}