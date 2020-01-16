package rpc

import (
	"sync"
	"time"
)

type batchRequestChan chan *batchRequest

type batchRequest struct {
	id         uint64
	name       string
	argsBytes  []byte
	replyBytes chan []byte
	replyError chan error
	noRequest  bool
	noResponse bool
}
type batch struct {
	mut                sync.Mutex
	reqChan            batchRequestChan
	client             *client
	readyRequests      []*batchRequest
	maxBatchRequest    int
	maxDelayNanoSecond int
	closeChan          chan bool
}

func newBatch(c *client, maxDelayNanoSecond int, maxBatchRequest int) *batch {
	b := &batch{
		reqChan:            make(chan *batchRequest, DefaultMaxCacheRequest),
		client:             c,
		readyRequests:      make([]*batchRequest, 0),
		maxBatchRequest:    maxBatchRequest,
		maxDelayNanoSecond: maxDelayNanoSecond,
		closeChan:          make(chan bool, 1),
	}
	go b.run()
	return b
}

func (b *batch) GetMaxBatchRequest() int {
	return b.maxBatchRequest
}
func (b *batch) SetMaxBatchRequest(maxBatchRequest int) {
	b.maxBatchRequest = maxBatchRequest
}
func (b *batch) run() {
	go func() {
		for cr := range b.reqChan {
			func() {
				b.mut.Lock()
				defer b.mut.Unlock()
				b.readyRequests = append(b.readyRequests, cr)
				if len(b.readyRequests) >= b.maxBatchRequest {
					crs := b.readyRequests[:]
					b.readyRequests = nil
					b.readyRequests = make([]*batchRequest, 0)
					b.Ticker(crs)
				}
			}()
		}
	}()
	tick := time.NewTicker(1 * time.Nanosecond * time.Duration(b.maxDelayNanoSecond))
	for {
		select {
		case <-b.closeChan:
			close(b.closeChan)
			tick.Stop()
			tick = nil
			goto endfor
		case <-tick.C:
			func() {
				b.mut.Lock()
				defer b.mut.Unlock()
				if len(b.readyRequests) > b.maxBatchRequest {
					crs := b.readyRequests[:b.maxBatchRequest]
					b.readyRequests = b.readyRequests[b.maxBatchRequest:]
					b.Ticker(crs)
				} else if len(b.readyRequests) > 0 {
					crs := b.readyRequests[:]
					b.readyRequests = b.readyRequests[len(b.readyRequests):]
					b.Ticker(crs)
				}
			}()
		}
	}
endfor:
}
func (b *batch) Ticker(brs []*batchRequest) {
	NoResponseCnt := 0
	var noResponse bool
	clientCodec := &clientCodec{}
	clientCodec.clientID = b.client.GetID()
	clientCodec.batch = true
	clientCodec.batchingAsync = b.client.batchingAsync
	clientCodec.requests = brs
	clientCodec.compressType = b.client.compressType
	clientCodec.compressLevel = b.client.compressLevel
	clientCodec.funcsCodecType = b.client.funcsCodecType
	for _, v := range brs {
		if v.noResponse == true {
			NoResponseCnt++
		}
	}
	if NoResponseCnt == len(brs) {
		noResponse = true
	} else {
		noResponse = false
	}
	msgBytes, err := clientCodec.Encode()
	if err == nil {
		if noResponse == false {
			data, err := b.client.remoteCall(msgBytes)
			if err == nil {
				err := clientCodec.Decode(data)
				if err != nil {
					return
				}
				if len(clientCodec.responses) == len(brs) {
					for i, v := range brs {
						if brs[i].id == clientCodec.responses[i].id && v.noResponse == false {
							func() {
								defer func() {
									if err := recover(); err != nil {
										logger.Errorln("v.reply err", err)
									}
								}()
								if clientCodec.responses[i].err != nil {
									v.replyError <- clientCodec.responses[i].err
								} else {
									v.replyBytes <- clientCodec.responses[i].data
								}
							}()
						}
					}
				}
			}
		} else {
			_ = b.client.remoteCallNoResponse(msgBytes)
		}

	}
}

func (b *batch) Close() {
	close(b.reqChan)
	b.client = nil
	b.readyRequests = nil
	b.closeChan <- true
}
