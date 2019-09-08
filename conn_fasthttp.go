package rpc
import (
	"github.com/valyala/fasthttp"
	"hslam.com/mgit/Mort/rpc/protocol"
	"net/url"
)
type FASTHTTPConn struct {
	conn 			*fasthttp.Client
	address			string
	url 			string
	CanWork			bool
}

func DialFASTHTTP(address string)  (Conn, error)  {
	u:=url.URL{Scheme: "http", Host: address, Path: "/"}
	t:=&FASTHTTPConn{
		conn:&fasthttp.Client{},
		address:address,
		url:u.String(),
	}
	return t, nil
}

func (t *FASTHTTPConn)Handle(readChan chan []byte,writeChan chan []byte, stopChan chan bool,finishChan chan bool){
	go protocol.HandleSyncConn(t, readChan,writeChan,stopChan,64)
}
func (t *FASTHTTPConn)TickerFactor()(int){
	return 100
}
func (t *FASTHTTPConn)BatchFactor()(int){
	return 64
}
func (t *FASTHTTPConn)Retry()(error){
	t.conn=&fasthttp.Client{}
	return nil
}
func (t *FASTHTTPConn)Close()(error){
	return nil
}
func (c *FASTHTTPConn)Do(requestBody []byte)([]byte,error) {
	req := &fasthttp.Request{}
	req.Header.SetMethod("POST")
	req.SetBody(requestBody)
	req.SetRequestURI(c.url)
	resp := &fasthttp.Response{}
	err := c.conn.Do(req, resp)
	return resp.Body(),err
}