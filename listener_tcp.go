package rpc

import (
	"net"
	"hslam.com/mgit/Mort/rpc/protocol"
	"hslam.com/mgit/Mort/rpc/log"
)

type TCPListener struct {
	server			*Server
	address			string
	netListener		net.Listener
}
func ListenTCP(address string,server *Server) (Listener, error) {
	lis, err := net.Listen("tcp", address)
	if err!=nil{
		log.Errorf("fatal error: %s", err)
		return nil,err
	}
	listener:= &TCPListener{address:address,netListener:lis,server:server}
	return listener,nil
}
func (l *TCPListener)Serve() (error) {
	log.Allf( "%s", "Waiting for clients")
	for {
		conn, err := l.netListener.Accept()
		if err != nil {
			log.Warnf("Accept: %s", err)
			continue
		}
		log.Infof("new client %s comming",conn.RemoteAddr())
		if l.server.useWorkerPool{
			l.server.workerPool.ProcessAsyn( func(obj interface{}, args ...interface{}) interface{} {
				var c = obj.(net.Conn )
				var server = args[0].(*Server)
				return ServeTCPConn(server,c)
			},conn,l.server)
		}else {
			go ServeTCPConn(l.server,conn)
		}
	}
	return nil
}
func (l *TCPListener)Addr() (string) {
	return l.address
}
func ServeTCPConn(server *Server,conn net.Conn)error {
	var RemoteAddr=conn.RemoteAddr().String()
	readChan := make(chan []byte)
	writeChan := make(chan []byte)
	finishChan:= make(chan bool,1)
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
					jobChan<-true
					priority,id,body,err:=UnpackFrame(data)
					if err!=nil{
						return
					}
					_,res_bytes, _ := server.ServeRPC(body)
					if res_bytes!=nil{
						frameBytes:=PacketFrame(priority,id,res_bytes)
						writeChan <- frameBytes
					}
					<-jobChan
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
	log.Infof("client %s exiting",RemoteAddr)
	return ErrConnExit
}