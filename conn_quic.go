package rpc
import (
	"github.com/lucas-clemente/quic-go"
	"hslam.com/mgit/Mort/rpc/protocol"
	"crypto/tls"
	"hslam.com/mgit/Mort/rpc/log"
)

type QUICConn struct {
	conn  		quic.Stream
	address			string
	CanWork			bool
	readChan 		chan []byte
	writeChan 		chan []byte
	stopChan 		chan bool
}

func DialQUIC(address string)  (Conn, error)  {
	session, err := quic.DialAddr(address, &tls.Config{InsecureSkipVerify: true}, nil)
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}
	stream, err := session.OpenStreamSync()
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}
	t:=&QUICConn{
		conn:stream,
		address:address,
	}
	return t, nil
}

func (t *QUICConn)Handle(readChan chan []byte,writeChan chan []byte, stopChan chan bool){
	t.readChan=readChan
	t.writeChan=writeChan
	t.stopChan=stopChan
	t.handle()
}
func (t *QUICConn)handle(){
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
func (t *QUICConn)TickerFactor()(int){
	return 1000
}
func (t *QUICConn)BatchFactor()(int){
	return 1
}
func (t *QUICConn)Retry()(error){
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
func (t *QUICConn)Close()(error){
	return t.conn.Close()
}