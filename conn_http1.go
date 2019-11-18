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
type HTTP1Conn struct {
	conn 			*http.Client
	address			string
	url				string
	CanWork			bool
	closed			bool
	multiplexing 	bool
}

func DialHTTP1(address string)  (Conn, error)  {
	u:=url.URL{Scheme: "http", Host: address, Path: "/"}
	t:=&HTTP1Conn{
		conn:&http.Client{
			Transport: &http.Transport{
				DisableKeepAlives:false,
				MaxConnsPerHost:1,
			},
		},
		address:address,
		url:u.String(),
	}
	return t, nil
}
func (t *HTTP1Conn)Buffer(enable bool){
}
func (t *HTTP1Conn)Multiplexing(enable bool){
	t.multiplexing=enable
	if t.multiplexing{
		t.conn.Transport=&http.Transport{
			DisableKeepAlives:false,
			MaxConnsPerHost:runtime.NumCPU(),
			MaxIdleConnsPerHost:runtime.NumCPU(),
		}
	}else {
		t.conn.Transport=&http.Transport{
			DisableKeepAlives:false,
			MaxConnsPerHost:1,
			MaxIdleConnsPerHost:1,
		}
	}
}
func (t *HTTP1Conn)Handle(readChan chan []byte,writeChan chan []byte, stopChan chan bool,finishChan chan bool){
	go protocol.HandleSyncConn(t, readChan,writeChan,stopChan,64)
}
func (t *HTTP1Conn)TickerFactor()(int){
	return 100
}
func (t *HTTP1Conn)BatchFactor()(int){
	return 512
}
func (t *HTTP1Conn)Retry()(error){
	var Transport *http.Transport
	if t.multiplexing{
		Transport=&http.Transport{
			DisableKeepAlives:false,
			MaxConnsPerHost:runtime.NumCPU(),
			MaxIdleConnsPerHost:runtime.NumCPU(),
		}
	}else {
		Transport=&http.Transport{
			DisableKeepAlives:false,
			MaxConnsPerHost:1,
			MaxIdleConnsPerHost:1,
		}
	}
	if DefaultClientTimeout>0{
		Transport.IdleConnTimeout=time.Duration(DefaultClientTimeout) * time.Millisecond
	}
	t.conn=&http.Client{
		Transport:Transport,
	}
	return nil
}
func (t *HTTP1Conn)Close()(error){
	return nil
}
func (c *HTTP1Conn)Do(requestBody []byte)([]byte,error) {
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

func (t *HTTP1Conn)Closed()(bool){
	return t.closed
}