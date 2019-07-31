package rpc

type Transporter interface {
	Handle(readChan chan []byte,writeChan chan []byte, stopChan chan bool)
	TickerFactor()(int)
	BatchFactor()(int)
	Retry()(error)
	Close()(error)
}
