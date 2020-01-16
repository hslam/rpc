package rpc

import (
	"github.com/hslam/protocol"
	"net"
)

type udpConn struct {
	conn       net.Conn
	address    string
	CanWork    bool
	readChan   chan []byte
	writeChan  chan []byte
	stopChan   chan bool
	finishChan chan bool
	closed     bool
}

func dialUDP(address string) (Conn, error) {
	conn, err := net.Dial(UDP, address)
	if err != nil {
		logger.Errorf("fatal error: %s", err)
		return nil, err
	}
	t := &udpConn{
		conn:    conn,
		address: address,
	}
	return t, nil
}
func (t *udpConn) NoDelay(enable bool) {
}
func (t *udpConn) Multiplexing(enable bool) {
}
func (t *udpConn) Handle(readChan chan []byte, writeChan chan []byte, stopChan chan bool, finishChan chan bool) {
	t.readChan = readChan
	t.writeChan = writeChan
	t.stopChan = stopChan
	t.finishChan = finishChan
	t.handle()
}
func (t *udpConn) handle() {
	readChan := make(chan []byte, 1)
	writeChan := make(chan []byte, 1)
	finishChan := make(chan bool, 1)
	stopChan := make(chan bool, 1)
	go protocol.HandleMessage(t.conn, readChan, writeChan, stopChan, finishChan)
	go func() {
		t.closed = false
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
			case <-t.stopChan:
				func() {
					defer func() {
						if err := recover(); err != nil {
						}
					}()
					stopChan <- true
				}()
				goto endfor
			case <-finishChan:
				func() {
					defer func() {
						if err := recover(); err != nil {
						}
					}()
					t.finishChan <- true
				}()
				goto endfor
			}
		}
	endfor:
		close(readChan)
		close(writeChan)
		close(finishChan)
		close(stopChan)
		t.closed = true
	}()
}
func (t *udpConn) TickerFactor() int {
	return 1
}
func (t *udpConn) BatchFactor() int {
	return 1
}
func (t *udpConn) Retry() error {
	conn, err := net.Dial(UDP, t.address)
	if err != nil {
		logger.Errorf("fatal error: %s", err)
		return err
	}
	t.conn = conn
	t.handle()
	return nil
}
func (t *udpConn) Close() error {
	return t.conn.Close()
}

func (t *udpConn) Closed() bool {
	return t.closed
}
