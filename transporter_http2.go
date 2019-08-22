package rpc
import (
	"net/http"
	"net/url"
	"hslam.com/mgit/Mort/rpc/protocol"
	"io"
	"bytes"
	"io/ioutil"
	"crypto/tls"
	"time"
	"golang.org/x/net/http2"
)
type HTTP2Transporter struct {
	conn 			*http.Client
	address			string
	url				string
	CanWork			bool
}

func NewHTTP2Transporter(address string)  (Transporter, error)  {
	u:=url.URL{Scheme: "https", Host: address, Path: "/"}
	var tlsConfig *tls.Config
	tlsConfig=&tls.Config{InsecureSkipVerify: true}
	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
		IdleConnTimeout: time.Duration(DefaultClientTimeout) * time.Millisecond,
	}
	http2.ConfigureTransport(transport)
	t:=&HTTP2Transporter{
		conn:&http.Client{
			Transport: transport,
		},
		address:address,
		url:u.String(),
	}
	return t, nil
}

func (t *HTTP2Transporter)Handle(readChan chan []byte,writeChan chan []byte, stopChan chan bool){
	go protocol.HandleSyncClient(t, readChan,writeChan,stopChan)
}
func (t *HTTP2Transporter)TickerFactor()(int){
	return 100
}
func (t *HTTP2Transporter)BatchFactor()(int){
	return 64
}
func (t *HTTP2Transporter)Retry()(error){
	t.conn=&http.Client{
		Transport: &http.Transport{},
	}
	return nil
}
func (t *HTTP2Transporter)Close()(error){
	return nil
}
func (c *HTTP2Transporter)Do(requestBody []byte)([]byte,error) {
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