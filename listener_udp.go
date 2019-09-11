package rpc

import (
	"net"
	"hslam.com/mgit/Mort/rpc/protocol"
	"hslam.com/mgit/Mort/rpc/log"
)

type UDPListener struct {
	server			*Server
	address			string
	netUDPConn		*net.UDPConn
}
func ListenUDP(address string,server *Server) (Listener, error) {
	addr, err := net.ResolveUDPAddr("udp", address)
	if err!=nil{
		log.Errorf("fatal error: %s", err)
		return nil,err
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Errorf("fatal error: %s", err)
		return nil,err
	}
	listener:=  &UDPListener{address:address,netUDPConn:conn,server:server}
	return listener,nil
}

func (l *UDPListener)Serve() (error) {
	log.Allf( "%s", "Waiting for clients")
	readChan := make(chan *protocol.UDPMsg,10240)
	writeChan := make(chan  *protocol.UDPMsg,10240)
	go protocol.ReadUDPConn(l.netUDPConn, readChan)
	go protocol.WriteUDPConn(l.netUDPConn, writeChan)
	for {
		select {
		case udp_msg := <-readChan:
			var RemoteAddr=udp_msg.RemoteAddr.String()
			log.AllInfof("new client %s comming",RemoteAddr)
			if l.server.useWorkerPool{
				l.server.workerPool.ProcessAsyn( func(obj interface{}, args ...interface{}) interface{} {
					var udp_msg = obj.(*protocol.UDPMsg)
					var writeChan=args[0].(chan *protocol.UDPMsg)
					var server = args[1].(*Server)
					return ServeUDPConn(server,udp_msg,writeChan)
				},udp_msg,writeChan,l.server)
			}else {
				go ServeUDPConn(l.server,udp_msg,writeChan)
			}
		}
	}
	l.netUDPConn.Close()
	close(writeChan)
	close(readChan)
	return nil
}
func (l *UDPListener)Addr() (string) {
	return l.address
}
func ServeUDPConn(server *Server,udp_msg *protocol.UDPMsg,writeChan chan *protocol.UDPMsg)error {
	if server.multiplexing{
		priority,id,body,err:=UnpackFrame(udp_msg.Data)
		if err!=nil{
			return ErrConnExit
		}
		ok,res_bytes, _ := server.ServeRPC(body)
		if res_bytes!=nil{
			frameBytes:=PacketFrame(priority,id,res_bytes)
			writeChan <- &protocol.UDPMsg{udp_msg.ID,frameBytes,udp_msg.RemoteAddr}
		}else if ok{
			writeChan <- &protocol.UDPMsg{udp_msg.ID,nil,udp_msg.RemoteAddr}
		}
	}else {
		ok,res_bytes, _ := server.ServeRPC(udp_msg.Data)
		if res_bytes!=nil{
			writeChan <- &protocol.UDPMsg{udp_msg.ID,res_bytes,udp_msg.RemoteAddr}
		}else if ok{
			writeChan <- &protocol.UDPMsg{udp_msg.ID,nil,udp_msg.RemoteAddr}
		}
	}

	log.AllInfof("client %s exiting",udp_msg.RemoteAddr)
	return ErrConnExit
}
