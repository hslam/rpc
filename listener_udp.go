package rpc

import (
	"net"
	"sync"
	"hslam.com/mgit/Mort/rpc/protocol"
)

type UDPListener struct {
	reqMutex 		sync.Mutex
	address			string
	netUDPConn		*net.UDPConn
}
func ListenUDP(address string) (Listener, error) {
	addr, err := net.ResolveUDPAddr("udp", address)
	if err!=nil{
		Fatalf("fatal error: %s", err)
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		Fatalf("fatal error: %s", err)
	}
	listener:=  &UDPListener{address:address,netUDPConn:conn}
	return listener,nil
}

func (l *UDPListener)Serve() (error) {
	Allf( "%s", "Waiting for clients")
	readChan := make(chan *protocol.UDPMsg,10240)
	writeChan := make(chan  *protocol.UDPMsg,10240)
	stopChan := make(chan bool)
	go protocol.ReadUDPConn(l.netUDPConn, readChan, stopChan)
	go protocol.WriteUDPConn(l.netUDPConn, writeChan,stopChan)
	for {
		select {
		case udp_msg := <-readChan:
			var RemoteAddr=udp_msg.RemoteAddr.String()
			AllInfof("new client %s comming",RemoteAddr)
			if useWorkerPool{
				workerPool.ProcessAsyn( func(obj interface{}, args ...interface{}) interface{} {
					var udp_msg = obj.(*protocol.UDPMsg)
					var writeChan=args[0].(chan *protocol.UDPMsg)
					return ServeUDPConn(udp_msg,writeChan)
				},udp_msg,writeChan)
			}else {
				go ServeUDPConn(udp_msg,writeChan)
			}
		case stop := <-stopChan:
			if stop {
				goto endfor
			}
		}
	}
	endfor:
		l.netUDPConn.Close()
		close(writeChan)
		close(readChan)
		close(stopChan)
	return nil
}
func (l *UDPListener)Addr() (string) {
	return l.address
}
func ServeUDPConn(udp_msg *protocol.UDPMsg,writeChan chan *protocol.UDPMsg)error {
	ok,res_bytes, _ := ServeRPC(udp_msg.Data)
	if res_bytes!=nil{
		writeChan <- &protocol.UDPMsg{udp_msg.ID,res_bytes,udp_msg.RemoteAddr}
	}else if ok{
		writeChan <- &protocol.UDPMsg{udp_msg.ID,nil,udp_msg.RemoteAddr}
	}
	AllInfof("client %s exiting",udp_msg.RemoteAddr)
	return ErrConnExit
}
