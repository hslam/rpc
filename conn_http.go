package rpc
import (
	"net/http"
	"net/url"
	"hslam.com/git/x/rpc/protocol"
	"io"
	"bytes"
	"io/ioutil"
	"runtime"
	"time"
)
type HTTPConn struct {
	conn 			*http.Client
	address			string
	url				string
	CanWork			bool
	closed			bool
}

func DialHTTP(address string)  (Conn, error)  {
	u:=url.URL{Scheme: "http", Host: address, Path: "/"}
	t:=&HTTPConn{
		conn:&http.Client{
			Transport: &http.Transport{
				DisableKeepAlives:false,
				MaxConnsPerHost:1,
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
	Transport:= &http.Transport{
		DisableKeepAlives:false,
		MaxIdleConnsPerHost: runtime.NumCPU(),
	}
	if DefaultClientTimeout>0{
		Transport.IdleConnTimeout=time.Duration(DefaultClientTimeout) * time.Millisecond
	}
	t.conn=&http.Client{
		Transport:Transport,
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

func (t *HTTPConn)Closed()(bool){
	return t.closed
}