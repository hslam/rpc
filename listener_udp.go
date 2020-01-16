package rpc

import (
	"github.com/hslam/protocol"
	"net"
)

type udpListener struct {
	server     *Server
	address    string
	netUDPConn *net.UDPConn
	maxConnNum int
	connNum    int
}

func listenUDP(address string, server *Server) (Listener, error) {
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		logger.Errorf("fatal error: %s", err)
		return nil, err
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		logger.Errorf("fatal error: %s", err)
		return nil, err
	}
	listener := &udpListener{address: address, netUDPConn: conn, server: server, maxConnNum: DefaultMaxConnNum * server.asyncMax}
	return listener, nil
}

func (l *udpListener) Serve() error {
	logger.Noticef("%s\n", "waiting for clients")
	workerChan := make(chan bool, l.maxConnNum)
	connChange := make(chan int)
	go func() {
		for c := range connChange {
			l.connNum += c
		}
	}()
	readChan := make(chan *protocol.UDPMsg, l.maxConnNum)
	writeChan := make(chan *protocol.UDPMsg, l.maxConnNum)
	go protocol.ReadUDPConn(l.netUDPConn, readChan)
	go protocol.WriteUDPConn(l.netUDPConn, writeChan)
	for {
		select {
		case udpMsg := <-readChan:
			workerChan <- true
			go func() {
				defer func() {
					if err := recover(); err != nil {
					}
					<-workerChan
				}()
				connChange <- 1
				var RemoteAddr = udpMsg.RemoteAddr.String()
				logger.Noticef("client %s comming\n", RemoteAddr)
				l.ServeUDPConn(udpMsg, writeChan)
				logger.Noticef("client %s exiting\n", RemoteAddr)
				connChange <- -1
			}()
		}
	}
	l.netUDPConn.Close()
	close(writeChan)
	close(readChan)
	return nil
}
func (l *udpListener) Addr() string {
	return l.address
}
func (l *udpListener) ServeUDPConn(udpMsg *protocol.UDPMsg, writeChan chan *protocol.UDPMsg) error {
	if l.server.multiplexing {
		priority, id, body, err := protocol.UnpackFrame(udpMsg.Data)
		if err != nil {
			return ErrConnExit
		}
		ok, resBytes := l.server.Serve(body)
		if resBytes != nil {
			frameBytes := protocol.PacketFrame(priority, id, resBytes)
			writeChan <- &protocol.UDPMsg{udpMsg.ID, frameBytes, udpMsg.RemoteAddr}
		} else if ok {
			writeChan <- &protocol.UDPMsg{udpMsg.ID, nil, udpMsg.RemoteAddr}
		}
	} else {
		ok, resBytes := l.server.Serve(udpMsg.Data)
		if resBytes != nil {
			writeChan <- &protocol.UDPMsg{udpMsg.ID, resBytes, udpMsg.RemoteAddr}
		} else if ok {
			writeChan <- &protocol.UDPMsg{udpMsg.ID, nil, udpMsg.RemoteAddr}
		}
	}
	return ErrConnExit
}
