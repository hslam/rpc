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

func (t *UDPConn)Handle(readChan chan []byte,writeChan chan []byte, stopChan chan bool){
	t.readChan=readChan
	t.writeChan=writeChan
	t.stopChan=stopChan
	t.handle()
}
func (t *UDPConn)handle(){
	readChan:=make(chan []byte)
	writeChan:=make(chan []byte)
	stopChan:=make(chan bool)
	go protocol.HandleMessage(t.conn,readChan,writeChan,stopChan)
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