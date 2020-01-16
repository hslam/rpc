package rpc

import (
	"github.com/gorilla/websocket"
	"github.com/hslam/protocol"
	"net/url"
)

type wsConn struct {
	conn       *protocol.MsgConn
	address    string
	CanWork    bool
	readChan   chan []byte
	writeChan  chan []byte
	stopChan   chan bool
	finishChan chan bool
	closed     bool
}

func dialWS(address string) (Conn, error) {
	u := url.URL{Scheme: "ws", Host: address, Path: "/"}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		logger.Errorf("fatal error: %s", err)
		return nil, err
	}
	t := &wsConn{
		conn:    &protocol.MsgConn{MessageConn: &websocketConn{c}},
		address: address,
	}
	t.CanWork = true
	return t, nil
}
func (t *wsConn) NoDelay(enable bool) {
}
func (t *wsConn) Multiplexing(enable bool) {
}
func (t *wsConn) Handle(readChan chan []byte, writeChan chan []byte, stopChan chan bool, finishChan chan bool) {
	t.readChan = readChan
	t.writeChan = writeChan
	t.stopChan = stopChan
	t.finishChan = finishChan
	t.handle()
}
func (t *wsConn) handle() {
	readChan := make(chan []byte, 1)
	writeChan := make(chan []byte, 1)
	finishChan := make(chan bool, 2)
	stopReadConnChan := make(chan bool, 1)
	stopWriteConnChan := make(chan bool, 1)
	go protocol.ReadConn(t.conn, readChan, stopReadConnChan, finishChan)
	go protocol.WriteConn(t.conn, writeChan, stopWriteConnChan, finishChan)
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
			case stop := <-finishChan:
				if stop {
					stopReadConnChan <- true
					stopWriteConnChan <- true
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
				stopReadConnChan <- true
				stopWriteConnChan <- true
				goto endfor
			}
		}
	endfor:
		close(readChan)
		close(writeChan)
		close(finishChan)
		close(stopReadConnChan)
		close(stopWriteConnChan)
		t.closed = true
	}()
}
func (t *wsConn) TickerFactor() int {
	return 100
}
func (t *wsConn) BatchFactor() int {
	return 512
}
func (t *wsConn) Retry() error {
	u := url.URL{Scheme: "ws", Host: t.address, Path: "/"}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)

	if err != nil {
		logger.Errorf("fatal error: %s", err)
		return err
	}
	t.conn = &protocol.MsgConn{MessageConn: &websocketConn{c}}
	t.handle()
	return nil
}
func (t *wsConn) Close() error {
	return t.conn.Close()
}
func (t *wsConn) Closed() bool {
	return t.closed
}
