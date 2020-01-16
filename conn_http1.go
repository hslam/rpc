package rpc

import (
	"bytes"
	"github.com/hslam/protocol"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"runtime"
	"time"
)

type http1Conn struct {
	conn         *http.Client
	address      string
	url          string
	CanWork      bool
	closed       bool
	multiplexing bool
}

func dialHTTP1(address string) (Conn, error) {
	u := url.URL{Scheme: "http", Host: address, Path: "/"}
	t := &http1Conn{
		conn: &http.Client{
			Transport: &http.Transport{
				DisableKeepAlives: false,
				MaxConnsPerHost:   1,
			},
		},
		address: address,
		url:     u.String(),
	}
	return t, nil
}
func (t *http1Conn) NoDelay(enable bool) {
}
func (t *http1Conn) Multiplexing(enable bool) {
	t.multiplexing = enable
	if t.multiplexing {
		t.conn.Transport = &http.Transport{
			DisableKeepAlives:   false,
			MaxConnsPerHost:     runtime.NumCPU(),
			MaxIdleConnsPerHost: runtime.NumCPU(),
		}
	} else {
		t.conn.Transport = &http.Transport{
			DisableKeepAlives:   false,
			MaxConnsPerHost:     1,
			MaxIdleConnsPerHost: 1,
		}
	}
}
func (t *http1Conn) Handle(readChan chan []byte, writeChan chan []byte, stopChan chan bool, finishChan chan bool) {
	go protocol.HandleSyncConn(t, readChan, writeChan, stopChan, 64)
}
func (t *http1Conn) TickerFactor() int {
	return 100
}
func (t *http1Conn) BatchFactor() int {
	return 512
}
func (t *http1Conn) Retry() error {
	var Transport *http.Transport
	if t.multiplexing {
		Transport = &http.Transport{
			DisableKeepAlives:   false,
			MaxConnsPerHost:     runtime.NumCPU(),
			MaxIdleConnsPerHost: runtime.NumCPU(),
		}
	} else {
		Transport = &http.Transport{
			DisableKeepAlives:   false,
			MaxConnsPerHost:     1,
			MaxIdleConnsPerHost: 1,
		}
	}
	if DefaultClientTimeout > 0 {
		Transport.IdleConnTimeout = time.Duration(DefaultClientTimeout) * time.Millisecond
	}
	t.conn = &http.Client{
		Transport: Transport,
	}
	return nil
}
func (t *http1Conn) Close() error {
	return nil
}

func (t *http1Conn) Do(requestBody []byte) ([]byte, error) {
	var requestBodyReader io.Reader
	if requestBody != nil {
		requestBodyReader = bytes.NewReader(requestBody)
	}
	req, _ := http.NewRequest("POST", t.url, requestBodyReader)
	resp, err := t.conn.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func (t *http1Conn) Closed() bool {
	return t.closed
}
