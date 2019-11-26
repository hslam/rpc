package protocol

import (
	"github.com/gorilla/websocket"
	"io"
)
type WSConn struct{
	*websocket.Conn
	buf []byte
}

func (c *WSConn) Write(b []byte) (int, error) {
	err := c.WriteMessage(websocket.BinaryMessage,b)
	return len(b),err
}
func (c *WSConn) Read(p []byte) (n int, err error){
	if len(c.buf)>0{
		if len(p)>=len(c.buf){
			copy(p,c.buf)
			n=len(c.buf)
			c.buf=[]byte{}
			return n,nil
		}
		return 0,io.ErrShortBuffer
	}
	_, data, err := c.ReadMessage()
	if len(p)>=len(data){
		copy(p,data)
		return len(data),err
	}
	c.buf=data
	return 0,io.ErrShortBuffer
}
func (c *WSConn) Close()(error){
	return c.Conn.Close()
}
