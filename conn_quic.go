package rpc
import (
	"github.com/lucas-clemente/quic-go"
	"hslam.com/git/x/rpc/protocol"
	"crypto/tls"
	"hslam.com/git/x/rpc/log"
)

type QUICConn struct {
	conn  		quic.Stream
	address			string
	CanWork			bool
	readChan 		chan []byte
	writeChan 		chan []byte
	stopChan 		chan bool
	finishChan		chan bool
	closed			bool

}

func DialQUIC(address string)  (Conn, error)  {
	session, err := quic.DialAddr(address, &tls.Config{InsecureSkipVerify: true}, nil)
	if err != nil {
		log.Errorf("fatal error: %s", err)
		return nil,err
	}
	stream, err := session.OpenStreamSync()
	if err != nil {
		log.Errorf("fatal error: %s", err)
		return nil,err
	}
	t:=&QUICConn{
		conn:stream,
		address:address,
	}
	return t, nil
}

func (t *QUICConn)Handle(readChan chan []byte,writeChan chan []byte, stopChan chan bool,finishChan chan bool){
	t.readChan=readChan
	t.writeChan=writeChan
	t.stopChan=stopChan
	t.finishChan=finishChan
	t.handle()
}
func (t *QUICConn)handle(){
	readChan:=make(chan []byte,1)
	writeChan:=make(chan []byte,1)
	finishChan:= make(chan bool,2)
	stopReadStreamChan := make(chan bool,1)
	stopWriteStreamChan := make(chan bool,1)
	go protocol.ReadStream(t.conn, readChan, stopReadStreamChan,finishChan)
	go protocol.WriteStream(t.conn, writeChan, stopWriteStreamChan,finishChan)
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
		t.closed=true
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
func (t *QUICConn)Closed()(bool){
	return t.closed
}