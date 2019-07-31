package rpc
import (
	"github.com/lucas-clemente/quic-go"
	"hslam.com/mgit/Mort/rpc/protocol"
	"crypto/tls"
	"hslam.com/mgit/Mort/rpc/log"
)

type QUICTransporter struct {
	conn  		quic.Stream
	address			string
	CanWork			bool
	readChan 		chan []byte
	writeChan 		chan []byte
	stopChan 		chan bool
}

func NewQUICTransporter(address string)  (Transporter, error)  {
	session, err := quic.DialAddr(address, &tls.Config{InsecureSkipVerify: true}, nil)
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}
	stream, err := session.OpenStreamSync()
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}
	t:=&QUICTransporter{
		conn:stream,
		address:address,
	}
	return t, nil
}

func (t *QUICTransporter)Handle(readChan chan []byte,writeChan chan []byte, stopChan chan bool){
	t.readChan=readChan
	t.writeChan=writeChan
	t.stopChan=stopChan
	t.handle()
}
func (t *QUICTransporter)handle(){
	readChan:=make(chan []byte)
	writeChan:=make(chan []byte)
	stopChan:=make(chan bool)
	go protocol.ReadStream(t.conn, readChan, stopChan)
	go protocol.WriteStream(t.conn, writeChan, stopChan)
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
func (t *QUICTransporter)TickerFactor()(int){
	return 1000
}
func (t *QUICTransporter)BatchFactor()(int){
	return 1
}
func (t *QUICTransporter)Retry()(error){
	session, err := quic.DialAddr(t.address, &tls.Config{InsecureSkipVerify: true}, nil)
	if err != nil {
		log.Errorf("fatal error: %s", err)
		return err
	}
	stream, err := session.OpenStreamSync()
	if err != nil {
		log.Errorf("fatal error: %s", err)
		return err
	}
	t.conn=stream
	t.handle()
	return nil
}
func (t *QUICTransporter)Close()(error){
	return t.conn.Close()
}