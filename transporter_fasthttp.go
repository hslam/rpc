package rpc
import (
	"github.com/valyala/fasthttp"
	"hslam.com/mgit/Mort/rpc/protocol"
)
type FASTHTTPTransporter struct {
	fastclient 		*fasthttp.Client
	address			string
	CanWork			bool
}

func NewFASTHTTPTransporter(address string)  (Transporter, error)  {
	t:=&FASTHTTPTransporter{
		fastclient:&fasthttp.Client{},
		address:address,
	}
	return t, nil
}

func (t *FASTHTTPTransporter)Handle(readChan chan []byte,writeChan chan []byte, stopChan chan bool){
	go protocol.HandleFASTHTTP(t.fastclient,t.address, readChan,writeChan,stopChan)
}
func (t *FASTHTTPTransporter)TickerFactor()(int){
	return 100
}
func (t *FASTHTTPTransporter)BatchFactor()(int){
	return 64
}
func (t *FASTHTTPTransporter)Close()(error){
	return nil
}