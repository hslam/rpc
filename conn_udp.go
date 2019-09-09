package rpc
import (
	"hslam.com/mgit/Mort/rpc/protocol"
	"hslam.com/mgit/Mort/rpc/log"
	"net"
)

type UDPConn struct {
	conn 			net.Conn
	address			string
	CanWork			bool
	readChan 		chan []byte
	writeChan 		chan []byte
	stopChan 		chan bool
	finishChan chan bool
	closed			bool
}

func DialUDP(address string)  (Conn, error)  {
	conn, err := net.Dial(UDP, address)
	if err != nil {
		log.Errorf("fatal error: %s", err)
		return nil,err
	}
	t:=&UDPConn{
		conn:conn,
		address:address,
	}
	return t, nil
}

func (t *UDPConn)Handle(readChan chan []byte,writeChan chan []byte, stopChan chan bool,finishChan chan bool){
	t.readChan=readChan
	t.writeChan=writeChan
	t.stopChan=stopChan
	t.finishChan=finishChan
	t.handle()
}
func (t *UDPConn)handle(){
	readChan:=make(chan []byte)
	writeChan:=make(chan []byte)
	finishChan:=make(chan bool,1)
	stopChan:=make(chan bool,1)
	go protocol.HandleMessage(t.conn,readChan,writeChan,stopChan,finishChan)
	go func() {
		t.closed=false
		log.Traceln("UDPConn.handle start")
		for {
			select {
			case v:=<-readChan:
				t.readChan<-v
			case v:=<-t.writeChan:
				writeChan<-v
			case <-t.stopChan:
				stopChan<-true
				goto endfor
			case <-finishChan:
				t.finishChan<-true
				goto endfor
			}
		}
	endfor:
		close(readChan)
		close(writeChan)
		close(finishChan)
		close(stopChan)
		t.closed=true
		log.Traceln("UDPConn.handle end")
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
		log.Errorf("fatal error: %s", err)
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