package rpc

import (
	"github.com/lucas-clemente/quic-go"
	"hslam.com/mgit/Mort/rpc/protocol"
	"hslam.com/mgit/Mort/rpc/log"
)

type QUICListener struct {
	server			*Server
	address			string
	quicListener	quic.Listener
	maxConnNum		int
	connNum			int
}
func ListenQUIC(address string,server *Server) (Listener, error) {
	quic_listener, err := quic.ListenAddr(address, generateTLSConfig(), nil)
	if err!=nil{
		log.Errorf("fatal error: %s", err)
		return nil,err
	}
	listener:= &QUICListener{address:address,quicListener:quic_listener,server:server,maxConnNum:DefaultMaxConnNum}
	return listener,nil
}
func (l *QUICListener)Serve() (error) {
	log.Allf( "%s", "Waiting for clients")
	workerChan := make(chan bool,l.maxConnNum)
	connChange := make(chan int)
	go func() {
		for conn_change := range connChange {
			l.connNum += conn_change
		}
	}()
	for{
		sess, err := l.quicListener.Accept()
		if err != nil {
			log.Warnf("Accept: %s", err)
			continue
		}else{
			workerChan<-true
			go func() {
				defer func() {
					if err := recover(); err != nil {
					}
					<-workerChan
				}()
				connChange <- 1
				log.Infof("new client %s comming",sess.RemoteAddr())
				ServeQUICConn(l.server,sess)
				log.Infof("client %s exiting",sess.RemoteAddr())
				connChange <- -1
			}()
		}
	}
	return nil
}
func (l *QUICListener)Addr() (string) {
	return l.address
}
func ServeQUICConn(server *Server,sess quic.Session)error {
	stream, err := sess.AcceptStream()
	if err != nil {
		log.Errorln(err)
		panic(err)
		return err
	}
	readChan := make(chan []byte)
	writeChan := make(chan []byte)
	finishChan:= make(chan bool,1)
	stopReadStreamChan := make(chan bool,1)
	stopWriteStreamChan := make(chan bool,1)
	stopChan := make(chan bool,1)
	go protocol.ReadStream(stream, readChan, stopReadStreamChan,finishChan)
	go protocol.WriteStream(stream, writeChan, stopWriteStreamChan,finishChan)
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
	}else{
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
	defer stream.Close()
	close(readChan)
	close(writeChan)
	close(finishChan)
	close(stopReadStreamChan)
	close(stopWriteStreamChan)
	close(stopChan)
	return ErrConnExit
}