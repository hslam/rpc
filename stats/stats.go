package stats

import (
	"encoding/json"
	"fmt"
	"sort"
	"sync"
	"os"
	"time"
)
var Log =true

func SetLog(log bool)  {
	Log=log
}

type Stats struct {
	Connections int
	Parallels     int
	AvgDuration float64
	Duration    float64
	Sum         float64
	Times       []int
	ResponseOk	int64
	Errors      int64
}


func StartStats(numParallels int,totalCalls int, clients []Client) []byte {
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
	result := CalcStats(
		responseChannel,
		benchTime.Duration(),
		numConnections,
		numParallels,
	)
	return result
}
func getS(n int,char string) (s string) {
	if n<1{
		return
	}
	for i:=1;i<=n;i++{
		s+=char
	}
	return
}
func CalcStats(responseChannel chan *Response, duration int64,numConnections int,numParallels int) []byte {
	stats := &Stats{
		Connections: numConnections,
		Parallels:     numParallels,
		Times:       make([]int, len(responseChannel)),
		Duration:    float64(duration),
		AvgDuration: float64(duration),
	}

	i := 0
	for res := range responseChannel {
		stats.Sum += float64(res.Duration)
		stats.Times[i] = int(res.Duration)
		i++

		if res.Error {
			stats.Errors++
		}else {
			stats.ResponseOk++
		}

		if len(responseChannel) == 0 {
			break
		}
	}

	sort.Ints(stats.Times)
	PrintStats(stats)
	b, err := json.Marshal(&stats)
	if err != nil {
		fmt.Println(err)
	}
	return b
}

func PrintStats(allStats *Stats) {
	sort.Ints(allStats.Times)
	total := float64(len(allStats.Times))
	totalInt := int64(total)
	fmt.Println("==========================BENCHMARK==========================")
	fmt.Printf("Used Connections:\t\t%d\n", allStats.Connections)
	fmt.Printf("Used Parallel Calls:\t\t%d\n", allStats.Parallels)
	fmt.Printf("Total Number Of Calls:\t\t%d\n\n", totalInt)
	fmt.Println("===========================TIMINGS===========================")
	fmt.Printf("Total time passed:\t\t%.2fs\n", allStats.AvgDuration/1E6)
	fmt.Printf("Avg time per request:\t\t%.2fms\n", allStats.Sum/total/1000)
	fmt.Printf("Requests per second:\t\t%.2f\n", total/(allStats.AvgDuration/1E6))
	fmt.Printf("Median time per request:\t%.2fms\n", float64(allStats.Times[(totalInt-1)/2])/1000)
	fmt.Printf("99th percentile time:\t\t%.2fms\n", float64(allStats.Times[(totalInt/100*99)])/1000)
	fmt.Printf("Slowest time for request:\t%.2fms\n\n", float64(allStats.Times[totalInt-1]/1000))
	fmt.Println("==========================RESPONSES==========================")
	fmt.Printf("ResponseOk:\t\t\t%d (%.2f%%)\n", allStats.ResponseOk, float64(allStats.ResponseOk)/total*1e2)
	fmt.Printf("Errors:\t\t\t\t%d (%.2f%%)\n", allStats.Errors, float64(allStats.Errors)/total*1e2)
}
