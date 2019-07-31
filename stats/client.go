package stats

import (
	"sync"
	"os"
	"time"
	"fmt"
	"encoding/json"
)
type Client interface {
	Call()(int64,bool)
}
func StartClientStats(numParallels int,totalCalls int, clients []Client) []byte {
	responseChannel := make(chan *Response, totalCalls*2)
	benchTime := NewTimer()
	benchTime.Reset()
	wg := &sync.WaitGroup{}
	numConnections:=len(clients)
	for i := 0; i < numConnections; i++ {
		go startClient(
			responseChannel,
			wg,
			numParallels,
			totalCalls,
			clients[i],
		)
		wg.Add(1)
	}
	var stopLog=false
	if Log{
		go func() {
			for {
				if len(responseChannel) >= totalCalls||stopLog{
					break
				}
				i:=len(responseChannel)*100/totalCalls
				fmt.Fprintf(os.Stdout, "%d%% [%s]\r",i,getS(i,"#") + getS(100-i," "))
				time.Sleep(time.Millisecond * 100)
			}
		}()
	}
	wg.Wait()
	stopLog=true
	if Log{
		fmt.Fprintf(os.Stdout, "%s\r",getS(106," "))
	}
	stats := CalcStats(
		responseChannel,
		benchTime.Duration(),
		numConnections,
		numParallels,
	)
	statsResult:=CalcStatsResult(stats)
	PrintStatsResult(statsResult)
	b, err := json.Marshal(&stats)
	if err != nil {
		fmt.Println(err)
	}
	return b
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
		bodySize,ok:=c.Call()
		respObj.Error=!ok
		respObj.Size=bodySize
		respObj.Duration = timer.Duration()
		if len(responseChannel) >= totalCalls {
			break
		}
		responseChannel <- respObj
	}
}