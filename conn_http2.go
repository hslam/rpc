package rpc
import (
	"net/http"
	"net/url"
	"github.com/hslam/protocol"
	"io"
	"bytes"
	"io/ioutil"
	"crypto/tls"
	"time"
	"golang.org/x/net/http2"
	"runtime"
)
type HTTP2Conn struct {
	conn 			*http.Client
	address			string
	url				string
	CanWork			bool
	closed			bool
	multiplexing 	bool
}

func DialHTTP2(address string)  (Conn, error)  {
	u:=url.URL{Scheme: "https", Host: address, Path: "/"}
	var tlsConfig *tls.Config
	tlsConfig=&tls.Config{InsecureSkipVerify: true}
	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
		DisableKeepAlives:false,
		MaxConnsPerHost:1,
	}
	if DefaultClientTimeout>0{
		transport.IdleConnTimeout=time.Duration(DefaultClientTimeout) * time.Millisecond
	}
	http2.ConfigureTransport(transport)
	t:=&HTTP2Conn{
		conn:&http.Client{
			Transport: transport,
		},
		address:address,
		url:u.String(),
	}
	return t, nil
}
func (t *HTTP2Conn)Buffer(enable bool){
}
func (t *HTTP2Conn)Multiplexing(enable bool){
	t.multiplexing=enable
	if t.multiplexing{
		var tlsConfig *tls.Config
		tlsConfig=&tls.Config{InsecureSkipVerify: true}
		transport := &http.Transport{
			TLSClientConfig: tlsConfig,
			DisableKeepAlives:false,
			MaxConnsPerHost:runtime.NumCPU(),
			MaxIdleConnsPerHost:runtime.NumCPU(),
		}
		if DefaultClientTimeout>0{
			transport.IdleConnTimeout=time.Duration(DefaultClientTimeout) * time.Millisecond
		}
		http2.ConfigureTransport(transport)
		t.conn.Transport=transport
	}else {
		var tlsConfig *tls.Config
		tlsConfig=&tls.Config{InsecureSkipVerify: true}
		transport := &http.Transport{
			TLSClientConfig: tlsConfig,
			DisableKeepAlives:false,
			MaxConnsPerHost:1,
			MaxIdleConnsPerHost:1,
		}
		if DefaultClientTimeout>0{
			transport.IdleConnTimeout=time.Duration(DefaultClientTimeout) * time.Millisecond
		}
		http2.ConfigureTransport(transport)
		t.conn.Transport=transport
	}
}
func (t *HTTP2Conn)Handle(readChan chan []byte,writeChan chan []byte, stopChan chan bool,finishChan chan bool){
	go protocol.HandleSyncConn(t, readChan,writeChan,stopChan,64)
}
func (t *HTTP2Conn)TickerFactor()(int){
	return 100
}
func (t *HTTP2Conn)BatchFactor()(int){
	return 512
}
func (t *HTTP2Conn)Retry()(error){
	if t.multiplexing{
		var tlsConfig *tls.Config
		tlsConfig=&tls.Config{InsecureSkipVerify: true}
		transport := &http.Transport{
			TLSClientConfig: tlsConfig,
			DisableKeepAlives:false,
			MaxConnsPerHost:runtime.NumCPU(),
			MaxIdleConnsPerHost:runtime.NumCPU(),
		}
		if DefaultClientTimeout>0{
			transport.IdleConnTimeout=time.Duration(DefaultClientTimeout) * time.Millisecond
		}
		http2.ConfigureTransport(transport)
		t.conn.Transport=transport
	}else {
		var tlsConfig *tls.Config
		tlsConfig=&tls.Config{InsecureSkipVerify: true}
		transport := &http.Transport{
			TLSClientConfig: tlsConfig,
			DisableKeepAlives:false,
			MaxConnsPerHost:1,
			MaxIdleConnsPerHost:1,
		}
		if DefaultClientTimeout>0{
			transport.IdleConnTimeout=time.Duration(DefaultClientTimeout) * time.Millisecond
		}
		http2.ConfigureTransport(transport)
		t.conn.Transport=transport
	}
	return nil
}
func (t *HTTP2Conn)Close()(error){
	return nil
}
func (c *HTTP2Conn)Do(requestBody []byte)([]byte,error) {
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
func (t *HTTP2Conn)Closed()(bool){
	return t.closed
}