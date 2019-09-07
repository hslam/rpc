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
type HTTPConn struct {
	conn 			*http.Client
	address			string
	url				string
	CanWork			bool
}

func DialHTTP(address string)  (Conn, error)  {
	u:=url.URL{Scheme: "http", Host: address, Path: "/"}
	t:=&HTTPConn{
		conn:&http.Client{
			Transport: &http.Transport{
				DisableKeepAlives:false,
				MaxConnsPerHost:DefaultMaxPipelineRequest+1,
				MaxIdleConns:DefaultMaxPipelineRequest+1,
				MaxIdleConnsPerHost:DefaultMaxPipelineRequest+1,
				IdleConnTimeout: time.Duration(DefaultClientTimeout) * time.Millisecond,
			},
		},
		address:address,
		url:u.String(),
	}
	return t, nil
}

func (t *HTTPConn)Handle(readChan chan []byte,writeChan chan []byte, stopChan chan bool,finishChan chan bool){
	go protocol.HandleSyncConn(t, readChan,writeChan,stopChan,64)
}
func (t *HTTPConn)TickerFactor()(int){
	return 100
}
func (t *HTTPConn)BatchFactor()(int){
	return 64
}
func (t *HTTPConn)Retry()(error){
	t.conn=&http.Client{
		Transport: &http.Transport{},
	}
	return nil
}
func (t *HTTPConn)Close()(error){
	return nil
}
func (c *HTTPConn)Do(requestBody []byte)([]byte,error) {
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