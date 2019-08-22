package rpc
import (
	"github.com/valyala/fasthttp"
	"hslam.com/mgit/Mort/rpc/protocol"
	"net/url"
)
type FASTHTTPTransporter struct {
	conn 		*fasthttp.Client
	address			string
	url 			string
	CanWork			bool
}

func NewFASTHTTPTransporter(address string)  (Transporter, error)  {
	u:=url.URL{Scheme: "http", Host: address, Path: "/"}
	t:=&FASTHTTPTransporter{
		conn:&fasthttp.Client{},
		address:address,
		url:u.String(),
	}
	return t, nil
}

func (t *FASTHTTPTransporter)Handle(readChan chan []byte,writeChan chan []byte, stopChan chan bool){
	go protocol.HandleSyncClient(t, readChan,writeChan,stopChan)
}
func (t *FASTHTTPTransporter)TickerFactor()(int){
	return 100
}
func (t *FASTHTTPTransporter)BatchFactor()(int){
	return 64
}
func (t *FASTHTTPTransporter)Retry()(error){
	t.conn=&fasthttp.Client{}
	return nil
}
func (t *FASTHTTPTransporter)Close()(error){
	return nil
}
func (c *FASTHTTPTransporter)Do(requestBody []byte)([]byte,error) {
	req := &fasthttp.Request{}
	req.Header.SetMethod("POST")
	req.SetBody(requestBody)
	req.SetRequestURI(c.url)
	resp := &fasthttp.Response{}
	err := c.conn.Do(req, resp)
	return resp.Body(),err
}