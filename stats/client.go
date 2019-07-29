package stats

import (
	"sync"
)
type Client interface {
	Call()bool
}

func startClient(responseChannel chan *Response, waitGroup *sync.WaitGroup, numParallels int,totalCalls int,c Client) {
	defer waitGroup.Done()
	wg := &sync.WaitGroup{}
	for i:=0;i<numParallels;i++{
		go func() {
			run(responseChannel,wg,totalCalls,c)
		}()
		wg.Add(1)
	}
	wg.Wait()
}

func run (responseChannel chan *Response,waitGroup *sync.WaitGroup,totalCalls int, c Client){
	defer waitGroup.Done()
	timer := NewTimer()
	for {
		timer.Reset()
		respObj := &Response{}
		ok:=c.Call()
		respObj.Error=!ok
		respObj.Duration = timer.Duration()
		if len(responseChannel) >= totalCalls {
			break
		}
		responseChannel <- respObj
	}
}