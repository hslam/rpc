package rpc
import (
	"hslam.com/mgit/Mort/rpc/protocol"
	"hslam.com/mgit/Mort/rpc/log"
	"net"
)

type TCPConn struct {
	conn 			*net.TCPConn
	address			string
	CanWork			bool
	readChan 		chan []byte
	writeChan 		chan []byte
	stopChan 		chan bool
	finishChan		chan bool
	closed			bool
}

func DialTCP(address string)  (Conn, error)  {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", address)
	if err != nil {
		log.Errorf("fatal error: %s", err)
		return nil,err
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Errorf("fatal error: %s", err)
		return nil,err
	}
	t:=&TCPConn{
		conn:conn,
		address:address,
	}
	return t, nil
}

func (t *TCPConn)Handle(readChan chan []byte,writeChan chan []byte, stopChan chan bool,finishChan chan bool){
	t.readChan=readChan
	t.writeChan=writeChan
	t.stopChan=stopChan
	t.finishChan=finishChan
	t.handle()
}
func (t *TCPConn)handle(){
	readChan:=make(chan []byte)
	writeChan:=make(chan []byte)
	finishChan:= make(chan bool)
	stopReadStreamChan := make(chan bool,1)
	stopWriteStreamChan := make(chan bool,1)
	go protocol.ReadStream(t.conn, readChan, stopReadStreamChan,finishChan)
	go protocol.WriteStream(t.conn, writeChan, stopWriteStreamChan,finishChan)
	go func() {
		t.closed=false
		//log.Traceln("TCPConn.handle start")
		for {
			select {
			case v:=<-readChan:
				t.readChan<-v
			case v:=<-t.writeChan:
				writeChan<-v
			case stop := <-finishChan:
				if stop {
					stopReadStreamChan<-true
					stopWriteStreamChan<-true
					t.finishChan<-true
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
			//log.Traceln("TCPConn.handle end")
		t.closed=true
	}()
}
func (t *TCPConn)TickerFactor()(int){
	return 300
}
func (t *TCPConn)BatchFactor()(int){
	return 64
}
func (t *TCPConn)Retry()(error){
	tcpAddr, err := net.ResolveTCPAddr("tcp4", t.address)
	if err != nil {
		log.Errorf("fatal error: %s", err)
		return err
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Errorf("fatal error: %s", err)
		return err
	}
	t.conn=conn
	t.handle()
	//log.Traceln("TCPConn.Retry")
	return nil
}
func (t *TCPConn)Close()(error){
	return t.conn.Close()
}
func (t *TCPConn)Closed()(bool){
	return t.closed
}