package rpc
import (
	"hslam.com/git/x/protocol"
	"net"
)

type UDPConn struct {
	conn 			net.Conn
	address			string
	CanWork			bool
	readChan 		chan []byte
	writeChan 		chan []byte
	stopChan 		chan bool
	finishChan 		chan bool
	closed			bool
}

func DialUDP(address string)  (Conn, error)  {
	conn, err := net.Dial(UDP, address)
	if err != nil {
		Errorf("fatal error: %s", err)
		return nil,err
	}
	t:=&UDPConn{
		conn:conn,
		address:address,
	}
	return t, nil
}
func (t *UDPConn)Buffer(enable bool){
}
func (t *UDPConn)Multiplexing(enable bool){
}
func (t *UDPConn)Handle(readChan chan []byte,writeChan chan []byte, stopChan chan bool,finishChan chan bool){
	t.readChan=readChan
	t.writeChan=writeChan
	t.stopChan=stopChan
	t.finishChan=finishChan
	t.handle()
}
func (t *UDPConn)handle(){
	readChan:=make(chan []byte,1)
	writeChan:=make(chan []byte,1)
	finishChan:=make(chan bool,1)
	stopChan:=make(chan bool,1)
	go protocol.HandleMessage(t.conn,readChan,writeChan,stopChan,finishChan)
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
			case <-t.stopChan:
				func() {
					defer func() {
						if err := recover(); err != nil {
						}
					}()
					stopChan<-true
				}()
				goto endfor
			case <-finishChan:
				func(){
					defer func() {if err := recover(); err != nil {}}()
					t.finishChan<-true
				}()
				goto endfor
			}
		}
	endfor:
		close(readChan)
		close(writeChan)
		close(finishChan)
		close(stopChan)
		t.closed=true
	}()
}
func (t *UDPConn)TickerFactor()(int){
	return 1
}
func (t *UDPConn)BatchFactor()(int){
	return 1
}
func (t *UDPConn)Retry()(error){
	conn, err := net.Dial(UDP, t.address)
	if err != nil {
		Errorf("fatal error: %s", err)
		return err
	}
	t.conn=conn
	t.handle()
	return nil
}
func (t *UDPConn)Close()(error){
	return t.conn.Close()
}

func (t *UDPConn)Closed()(bool){
	return t.closed
}