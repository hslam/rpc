package rpc

import (
	"github.com/hslam/protocol"
	"net"
)

type UDPListener struct {
	server     *Server
	address    string
	netUDPConn *net.UDPConn
	maxConnNum int
	connNum    int
}

func ListenUDP(address string, server *Server) (Listener, error) {
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		Errorf("fatal error: %s", err)
		return nil, err
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		Errorf("fatal error: %s", err)
		return nil, err
	}
	listener := &UDPListener{address: address, netUDPConn: conn, server: server, maxConnNum: DefaultMaxConnNum * server.asyncMax}
	return listener, nil
}

func (l *UDPListener) Serve() error {
	Allf("%s\n", "waiting for clients")
	workerChan := make(chan bool, l.maxConnNum)
	connChange := make(chan int)
	go func() {
		for conn_change := range connChange {
			l.connNum += conn_change
		}
	}()
	readChan := make(chan *protocol.UDPMsg, l.maxConnNum)
	writeChan := make(chan *protocol.UDPMsg, l.maxConnNum)
	go protocol.ReadUDPConn(l.netUDPConn, readChan)
	go protocol.WriteUDPConn(l.netUDPConn, writeChan)
	for {
		select {
		case udp_msg := <-readChan:
			workerChan <- true
			go func() {
				defer func() {
					if err := recover(); err != nil {
					}
					<-workerChan
				}()
				connChange <- 1
				var RemoteAddr = udp_msg.RemoteAddr.String()
				AllInfof("client %s comming\n", RemoteAddr)
				l.ServeUDPConn(udp_msg, writeChan)
				AllInfof("client %s exiting\n", RemoteAddr)
				connChange <- -1
			}()
		}
	}
	l.netUDPConn.Close()
	close(writeChan)
	close(readChan)
	return nil
}
func (l *UDPListener) Addr() string {
	return l.address
}
func (l *UDPListener) ServeUDPConn(udp_msg *protocol.UDPMsg, writeChan chan *protocol.UDPMsg) error {
	if l.server.multiplexing {
		priority, id, body, err := protocol.UnpackFrame(udp_msg.Data)
		if err != nil {
			return ErrConnExit
		}
		ok, res_bytes := l.server.Serve(body)
		if res_bytes != nil {
			frameBytes := protocol.PacketFrame(priority, id, res_bytes)
			writeChan <- &protocol.UDPMsg{udp_msg.ID, frameBytes, udp_msg.RemoteAddr}
		} else if ok {
			writeChan <- &protocol.UDPMsg{udp_msg.ID, nil, udp_msg.RemoteAddr}
		}
	} else {
		ok, res_bytes := l.server.Serve(udp_msg.Data)
		if res_bytes != nil {
			writeChan <- &protocol.UDPMsg{udp_msg.ID, res_bytes, udp_msg.RemoteAddr}
		} else if ok {
			writeChan <- &protocol.UDPMsg{udp_msg.ID, nil, udp_msg.RemoteAddr}
		}
	}
	return ErrConnExit
}
