package rpc
import (
	"net/http"
	"net/url"
	"hslam.com/mgit/Mort/rpc/protocol"
	"io"
	"bytes"
	"io/ioutil"
	"time"
)
type HTTPTransporter struct {
	conn 			*http.Client
	address			string
	url				string
	CanWork			bool
}

func NewHTTPTransporter(address string)  (Transporter, error)  {
	u:=url.URL{Scheme: "http", Host: address, Path: "/"}
	t:=&HTTPTransporter{
		conn:&http.Client{
			Transport: &http.Transport{
				DisableKeepAlives:false,
				MaxConnsPerHost:1,
				MaxIdleConns:1,
				MaxIdleConnsPerHost:1,
				IdleConnTimeout: time.Duration(DefaultClientTimeout) * time.Millisecond,
			},
		},
		address:address,
		url:u.String(),
	}
	return t, nil
}

func (t *HTTPTransporter)Handle(readChan chan []byte,writeChan chan []byte, stopChan chan bool){
	go protocol.HandleSyncClient(t, readChan,writeChan,stopChan)
}
func (t *HTTPTransporter)TickerFactor()(int){
	return 100
}
func (t *HTTPTransporter)BatchFactor()(int){
	return 64
}
func (t *HTTPTransporter)Retry()(error){
	t.conn=&http.Client{
		Transport: &http.Transport{},
	}
	return nil
}
func (t *HTTPTransporter)Close()(error){
	return nil
}
func (c *HTTPTransporter)Do(requestBody []byte)([]byte,error) {
	var requestBodyReader io.Reader
	if requestBody!=nil{
		requestBodyReader = bytes.NewReader(requestBody)
	}
	req, _ := http.NewRequest("POST", c.url, requestBodyReader)
	resp, err :=c.conn.Do(req)
	if err!=nil{
		return nil,err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}