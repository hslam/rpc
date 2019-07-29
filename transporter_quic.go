package rpc
import (
	"github.com/lucas-clemente/quic-go"
	"hslam.com/mgit/Mort/rpc/protocol"
	"crypto/tls"
)

type QUICTransporter struct {
	stream  		quic.Stream
	address			string
	CanWork			bool
}

func NewQUICTransporter(address string)  (Transporter, error)  {
	session, err := quic.DialAddr(address, &tls.Config{InsecureSkipVerify: true}, nil)
	if err != nil {
		Fatalf("fatal error: %s", err)
	}
	stream, err := session.OpenStreamSync()
	if err != nil {
		Fatalf("fatal error: %s", err)
	}
	t:=&QUICTransporter{
		stream:stream,
		address:address,
	}
	return t, nil
}

func (t *QUICTransporter)Handle(readChan chan []byte,writeChan chan []byte, stopChan chan bool){
	go protocol.ReadStream(t.stream, readChan, stopChan)
	go protocol.WriteStream(t.stream, writeChan, stopChan)
}
func (t *QUICTransporter)TickerFactor()(int){
	return 1000
}
func (t *QUICTransporter)BatchFactor()(int){
	return 1
}
func (t *QUICTransporter)Close()(error){
	return t.stream.Close()
}