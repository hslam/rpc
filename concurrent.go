package rpc


type ConcurrentRequestChan chan *ConcurrentRequest

type ConcurrentRequest struct {
	id uint32
	data []byte
	noResponse  bool
	cbChan chan []byte
}

type MultiPlexing struct {
	concurrentRequestChan ConcurrentRequestChan
	actionConcurrentRequestChan map[uint32]*ConcurrentRequest
	noResponseConcurrentRequestChan map[uint32]*ConcurrentRequest
	readChan 				chan []byte
	writeChan 				chan []byte
	maxConcurrentRequest	int
	stop					bool
}
