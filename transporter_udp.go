package rpc
import (
	"hslam.com/mgit/Mort/rpc/protocol"
	"hslam.com/mgit/Mort/rpc/log"
	"net"
)

type UDPTransporter struct {
	conn 			net.Conn
	address			string
	CanWork			bool
	readChan 		chan []byte
	writeChan 		chan []byte
	stopChan 		chan bool
}

func NewUDPTransporter(address string)  (Transporter, error)  {
	conn, err := net.Dial(UDP, address)
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}
	t:=&UDPTransporter{
		conn:conn,
		address:address,
	}
	return t, nil
}

func (t *UDPTransporter)Handle(readChan chan []byte,writeChan chan []byte, stopChan chan bool){
	t.readChan=readChan
	t.writeChan=writeChan
	t.stopChan=stopChan
	t.handle()
}
func (t *UDPTransporter)handle(){
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
func (t *UDPTransporter)TickerFactor()(int){
	return 1
}
func (t *UDPTransporter)BatchFactor()(int){
	return 1
}
func (t *UDPTransporter)Retry()(error){
	conn, err := net.Dial(UDP, t.address)
	if err != nil {
		log.Errorf("fatal error: %s", err)
		return err
	}
	t.conn=conn
	t.handle()
	return nil
}
func (t *UDPTransporter)Close()(error){
	return t.conn.Close()
}