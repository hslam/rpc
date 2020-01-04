package rpc
import (
	"github.com/hslam/protocol"
	"net"
)

type IPCConn struct {
	conn 			*net.UnixConn
	address			string
	CanWork			bool
	readChan 		chan []byte
	writeChan 		chan []byte
	stopChan 		chan bool
	finishChan		chan bool
	closed			bool
	noDelay 		bool
}

func DialIPC(address string)  (Conn, error)  {
	var addr *net.UnixAddr
	var err error
	if addr, err = net.ResolveUnixAddr("unix", address); err != nil {
		return nil, err
	}
	conn, err := net.DialUnix("unix", nil, addr)
	if err != nil {
		Errorf("fatal error: %s", err)
		return nil,err
	}
	t:=&IPCConn{
		conn:conn,
		address:address,
	}
	return t, nil
}

func (t *IPCConn)Handle(readChan chan []byte,writeChan chan []byte, stopChan chan bool,finishChan chan bool){
	t.readChan=readChan
	t.writeChan=writeChan
	t.stopChan=stopChan
	t.finishChan=finishChan
	t.handle()
}
func (t *IPCConn)NoDelay(enable bool){
	t.noDelay=enable
}
func (t *IPCConn)Multiplexing(enable bool){
}
func (t *IPCConn)handle(){
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
func (t *IPCConn)TickerFactor()(int){
	return 300
}
func (t *IPCConn)BatchFactor()(int){
	return 512
}
func (t *IPCConn)Retry()(error){
	var addr *net.UnixAddr
	var err error
	if addr, err = net.ResolveUnixAddr("unix", t.address); err != nil {
		return err
	}
	conn, err := net.DialUnix("unix", nil, addr)
	if err != nil {
		Errorf("fatal error: %s", err)
		return err
	}
	t.conn=conn
	t.handle()
	//Traceln("IPCConn.Retry")
	return nil
}
func (t *IPCConn)Close()(error){
	return t.conn.Close()
}
func (t *IPCConn)Closed()(bool){
	return t.closed
}