package rpc

import (
	"sync"
	"github.com/lucas-clemente/quic-go"
	"hslam.com/mgit/Mort/rpc/protocol"
)

type QUICListener struct {
	reqMutex 		sync.Mutex
	address			string
	quicListener	quic.Listener
}
func ListenQUIC(address string) (Listener, error) {
	quic_listener, err := quic.ListenAddr(address, generateTLSConfig(), nil)
	if err!=nil{
		Fatalf("fatal error: %s", err)
	}
	listener:= &QUICListener{address:address,quicListener:quic_listener}
	return listener,nil
}
func (l *QUICListener)Serve() (error) {
	Allf( "%s", "Waiting for clients")
	for{
		sess, err := l.quicListener.Accept()
		if err != nil {
			Warnf("Accept: %s", err)
			continue
		}else{
			Infof("new client %s comming",sess.RemoteAddr())
			if useWorkerPool{
				workerPool.ProcessAsyn( func(obj interface{}, args ...interface{}) interface{} {
					var s = obj.(quic.Session)
					return ServeQUICConn(s)
				},sess)
			}else {
				go ServeQUICConn(sess)
			}
		}
	}
	return nil
}
func (l *QUICListener)Addr() (string) {
	return l.address
}
func ServeQUICConn(sess quic.Session)error {
	var RemoteAddr=sess.RemoteAddr().String()
	stream, err := sess.AcceptStream()
	if err != nil {
		Errorln(err)
		panic(err)
		return err
	}
	readChan := make(chan []byte)
	writeChan := make(chan []byte)
	stopChan := make(chan bool)
	go protocol.ReadStream(stream, readChan, stopChan)
	go protocol.WriteStream(stream, writeChan, stopChan)
	for {
		select {
		case data := <-readChan:
			_,res_bytes, _ := ServeRPC(data)
			if res_bytes!=nil{
				writeChan <- res_bytes
			}
		case stop := <-stopChan:
			if stop {
				goto endfor
			}
		}
	}
	endfor:
	defer stream.Close()
	close(writeChan)
	close(readChan)
	close(stopChan)
	Infof("client %s exiting",RemoteAddr)
	return ErrConnExit
}