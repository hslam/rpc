package rpc
import (
	"hslam.com/mgit/Mort/rpc/protocol"
	"net"
)

type TCPTransporter struct {
	conn 			*net.TCPConn
	address			string
	CanWork			bool
}

func NewTCPTransporter(address string)  (Transporter, error)  {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", address)
	if err != nil {
		Fatalf("fatal error: %s", err)
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		Fatalf("fatal error: %s", err)
	}
	if err != nil {
		return nil, err
	}
	t:=&TCPTransporter{
		conn:conn,
		address:address,
	}
	return t, nil
}

func (t *TCPTransporter)Handle(readChan chan []byte,writeChan chan []byte, stopChan chan bool){
	go protocol.ReadStream(t.conn, readChan, stopChan)
	go protocol.WriteStream(t.conn, writeChan, stopChan)
}
func (t *TCPTransporter)TickerFactor()(int){
	return 300
}
func (t *TCPTransporter)BatchFactor()(int){
	return 64
}
func (t *TCPTransporter)Close()(error){
	return t.conn.Close()
}