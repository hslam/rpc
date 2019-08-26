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
}

func DialTCP(address string)  (Conn, error)  {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", address)
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}
	if err != nil {
		return nil, err
	}
	t:=&TCPConn{
		conn:conn,
		address:address,
	}
	return t, nil
}

func (t *TCPConn)Handle(readChan chan []byte,writeChan chan []byte, stopChan chan bool){
	t.readChan=readChan
	t.writeChan=writeChan
	t.stopChan=stopChan
	t.handle()
}
func (t *TCPConn)handle(){
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
	return nil
}
func (t *TCPConn)Close()(error){
	return t.conn.Close()
}
