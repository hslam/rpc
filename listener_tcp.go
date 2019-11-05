package rpc

import (
	"net"
	"hslam.com/git/x/rpc/protocol"
	"hslam.com/git/x/rpc/log"
)

type TCPListener struct {
	server			*Server
	address			string
	netListener		net.Listener
	maxConnNum		int
	connNum			int
}
func ListenTCP(address string,server *Server) (Listener, error) {
	lis, err := net.Listen("tcp", address)
	if err!=nil{
		log.Errorf("fatal error: %s", err)
		return nil,err
	}
	listener:= &TCPListener{address:address,netListener:lis,server:server,maxConnNum:DefaultMaxConnNum}
	return listener,nil
}
func (l *TCPListener)Serve() (error) {
	log.Allf( "%s", "Waiting for clients")
	workerChan := make(chan bool,l.maxConnNum)
	connChange := make(chan int)
	go func() {
		for conn_change := range connChange {
			l.connNum += conn_change
		}
	}()
	for {
		conn, err := l.netListener.Accept()
		if err != nil {
			log.Warnf("Accept: %s", err)
			continue
		}
		workerChan<-true
		go func() {
			defer func() {
				if err := recover(); err != nil {
				}
				<-workerChan
			}()
			connChange <- 1
			log.Infof("new client %s comming",conn.RemoteAddr())
			ServeTCPConn(l.server,conn)
			log.Infof("client %s exiting",conn.RemoteAddr())
			connChange <- -1
		}()
	}
	return nil
}
func (l *TCPListener)Addr() (string) {
	return l.address
}
func ServeTCPConn(server *Server,conn net.Conn)error {
	readChan := make(chan []byte,1)
	writeChan := make(chan []byte,1)
	finishChan:= make(chan bool,2)
	stopReadStreamChan := make(chan bool,1)
	stopWriteStreamChan := make(chan bool,1)
	stopChan := make(chan bool,1)
	go protocol.ReadStream(conn, readChan, stopReadStreamChan,finishChan)
	go protocol.WriteStream(conn, writeChan, stopWriteStreamChan,finishChan)
	if server.multiplexing{
		jobChan := make(chan bool,server.asyncMax)
		for {
			select {
			case data := <-readChan:
				go func(data []byte ,writeChan chan []byte) {
					defer func() {
						if err := recover(); err != nil {
						}
						<-jobChan
					}()
					jobChan<-true
					priority,id,body,err:=protocol.UnpackFrame(data)
					if err!=nil{
						return
					}
					_,res_bytes, _ := server.ServeRPC(body)
					if res_bytes!=nil{
						frameBytes:=protocol.PacketFrame(priority,id,res_bytes)
						writeChan <- frameBytes
					}
				}(data,writeChan)
			case stop := <-finishChan:
				if stop {
					stopReadStreamChan<-true
					stopWriteStreamChan<-true
					goto endfor
				}
			}
		}
	}else if server.async{
		syncConn:=newSyncConn(server)
		go protocol.HandleSyncConn(syncConn, writeChan,readChan,stopChan,server.asyncMax)
		select {
		case stop := <-finishChan:
			if stop {
				stopReadStreamChan<-true
				stopWriteStreamChan<-true
				stopChan<-true
				goto endfor
			}
		}
	}else {
		for {
			select {
			case data := <-readChan:
				_,res_bytes, _ := server.ServeRPC(data)
				if res_bytes!=nil{
					writeChan <- res_bytes
				}
			case stop := <-finishChan:
				if stop {
					stopReadStreamChan<-true
					stopWriteStreamChan<-true
					goto endfor
				}
			}
		}
	}
	endfor:
		defer conn.Close()
	close(readChan)
	close(writeChan)
	close(finishChan)
	close(stopReadStreamChan)
	close(stopWriteStreamChan)
	close(stopChan)
	return ErrConnExit
}