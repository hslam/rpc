package rpc

import (
	"github.com/hslam/protocol"
	"net"
)

type tcpConn struct {
	conn       *net.TCPConn
	address    string
	CanWork    bool
	readChan   chan []byte
	writeChan  chan []byte
	stopChan   chan bool
	finishChan chan bool
	closed     bool
	noDelay    bool
}

func dialTCP(address string) (Conn, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", address)
	if err != nil {
		logger.Errorf("fatal error: %s", err)
		return nil, err
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		logger.Errorf("fatal error: %s", err)
		return nil, err
	}
	conn.SetNoDelay(true)
	t := &tcpConn{
		conn:    conn,
		address: address,
	}
	return t, nil
}

func (t *tcpConn) Handle(readChan chan []byte, writeChan chan []byte, stopChan chan bool, finishChan chan bool) {
	t.readChan = readChan
	t.writeChan = writeChan
	t.stopChan = stopChan
	t.finishChan = finishChan
	t.handle()
}
func (t *tcpConn) NoDelay(enable bool) {
	t.noDelay = enable
}
func (t *tcpConn) Multiplexing(enable bool) {
}
func (t *tcpConn) handle() {
	readChan := make(chan []byte, 1)
	writeChan := make(chan []byte, 1)
	finishChan := make(chan bool, 2)
	stopReadStreamChan := make(chan bool, 1)
	stopWriteStreamChan := make(chan bool, 1)
	go protocol.ReadStream(t.conn, readChan, stopReadStreamChan, finishChan)
	go protocol.WriteStream(t.conn, writeChan, stopWriteStreamChan, finishChan, t.noDelay)
	go func() {
		t.closed = false
		//logger.Traceln("TCPConn.handle start")
		for {
			select {
			case v := <-readChan:
				func() {
					defer func() {
						if err := recover(); err != nil {
						}
					}()
					t.readChan <- v
				}()
			case v := <-t.writeChan:
				func() {
					defer func() {
						if err := recover(); err != nil {
						}
					}()
					writeChan <- v
				}()
			case stop := <-finishChan:
				if stop {
					stopReadStreamChan <- true
					stopWriteStreamChan <- true
					func() {
						defer func() {
							if err := recover(); err != nil {
							}
						}()
						t.finishChan <- true
					}()
					goto endfor
				}
			case <-t.stopChan:
				stopReadStreamChan <- true
				stopWriteStreamChan <- true
				goto endfor
			}
		}
	endfor:
		close(readChan)
		close(writeChan)
		close(finishChan)
		close(stopReadStreamChan)
		close(stopWriteStreamChan)
		//logger.Traceln("TCPConn.handle end")
		t.closed = true
	}()
}
func (t *tcpConn) TickerFactor() int {
	return 300
}
func (t *tcpConn) BatchFactor() int {
	return 512
}
func (t *tcpConn) Retry() error {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", t.address)
	if err != nil {
		logger.Errorf("fatal error: %s", err)
		return err
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		logger.Errorf("fatal error: %s", err)
		return err
	}
	t.conn = conn
	t.handle()
	//logger.Traceln("TCPConn.Retry")
	return nil
}
func (t *tcpConn) Close() error {
	return t.conn.Close()
}
func (t *tcpConn) Closed() bool {
	return t.closed
}
