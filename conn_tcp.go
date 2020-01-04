package rpc
import (
	"github.com/hslam/protocol"
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
	noDelay 		bool
}

func DialTCP(address string)  (Conn, error)  {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", address)
	if err != nil {
		Errorf("fatal error: %s", err)
		return nil,err
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		Errorf("fatal error: %s", err)
		return nil,err
	}
	conn.SetNoDelay(true)
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
func (t *TCPConn)NoDelay(enable bool){
	t.noDelay=enable
}
func (t *TCPConn)Multiplexing(enable bool){
}
func (t *TCPConn)handle(){
	readChan:=make(chan []byte,1)
	writeChan:=make(chan []byte,1)
	finishChan:= make(chan bool,2)
	stopReadStreamChan := make(chan bool,1)
	stopWriteStreamChan := make(chan bool,1)
	go protocol.ReadStream(t.conn, readChan, stopReadStreamChan,finishChan)
	go protocol.WriteStream(t.conn, writeChan, stopWriteStreamChan,finishChan,t.noDelay)
	go func() {
		t.closed=false
		//Traceln("TCPConn.handle start")
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
			//Traceln("TCPConn.handle end")
		t.closed=true
	}()
}
func (t *TCPConn)TickerFactor()(int){
	return 300
}
func (t *TCPConn)BatchFactor()(int){
	return 512
}
func (t *TCPConn)Retry()(error){
	tcpAddr, err := net.ResolveTCPAddr("tcp4", t.address)
	if err != nil {
		Errorf("fatal error: %s", err)
		return err
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		Errorf("fatal error: %s", err)
		return err
	}
	t.conn=conn
	t.handle()
	//Traceln("TCPConn.Retry")
	return nil
}
func (t *TCPConn)Close()(error){
	return t.conn.Close()
}
func (t *TCPConn)Closed()(bool){
	return t.closed
}