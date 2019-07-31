package stats

import (
	"fmt"
	"sort"
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
	Transfered  int64
	ResponseOk	int64
	Errors      int64
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

func CalcStats(responseChannel chan *Response, duration int64,numConnections int,numParallels int) *Stats {
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
		stats.Transfered += res.Size
		if res.Error {
			stats.Errors++
		}else {
			stats.ResponseOk++
		}

		if len(responseChannel) == 0 {
			break
		}
	}
	return stats
}

func CalcStatsResult(allStats *Stats) *StatsResult{
	sort.Ints(allStats.Times)
	total := float64(len(allStats.Times))
	totalInt := int64(total)
	var statsResult =&StatsResult{}
	statsResult.Connections=allStats.Connections
	statsResult.Parallels=allStats.Parallels
	statsResult.TotalCalls=totalInt
	statsResult.TotalTimePassed=allStats.AvgDuration/1E6
	statsResult.AvgTimePerRequest=allStats.Sum/total/1000
	statsResult.RequestsPerSecond=total/(allStats.AvgDuration/1E6)
	statsResult.MedianTimePerRequest=float64(allStats.Times[(totalInt-1)/2])/1000
	statsResult.N99thPercentileTime=float64(allStats.Times[(totalInt/100*99)])/1000
	statsResult.SlowestTimeForRequest=float64(allStats.Times[totalInt-1]/1000)
	statsResult.ResponseOk=allStats.ResponseOk
	statsResult.ResponseOkPercentile=float64(allStats.ResponseOk)/total*1e2
	statsResult.Errors=allStats.Errors
	statsResult.ErrorsPercentile=float64(allStats.Errors)/total*1e2
	if allStats.Transfered>0{
		statsResult.TotalResponseBodySizes=allStats.Transfered
		statsResult.AvgResponseBodyPerRequest=float64(allStats.Transfered)/total
		tr := float64(allStats.Transfered) / (allStats.AvgDuration / 1E6)
		statsResult.TransferRateBytePerSecond=tr
		statsResult.TransferRateMBytePerSecond=tr/1E6
	}
	return statsResult
}
type StatsResult struct {
	Connections 			int
	Parallels     			int
	TotalCalls				int64
	TotalTimePassed			float64
	AvgTimePerRequest		float64
	RequestsPerSecond 		float64
	MedianTimePerRequest	float64
	N99thPercentileTime		float64
	SlowestTimeForRequest 	float64
	TotalResponseBodySizes  int64
	AvgResponseBodyPerRequest  float64
	TransferRateBytePerSecond  float64
	TransferRateMBytePerSecond  float64
	ResponseOk				int64
	ResponseOkPercentile	float64
	Errors      			int64
	ErrorsPercentile		float64

}
func PrintStatsResult(statsResult *StatsResult) {
	fmt.Println("==========================BENCHMARK==========================")
	fmt.Printf("Used Connections:\t\t%d\n", statsResult.Connections)
	fmt.Printf("Used Parallel Calls:\t\t%d\n", statsResult.Parallels)
	fmt.Printf("Total Number Of Calls:\t\t%d\n\n", statsResult.TotalCalls)
	fmt.Println("===========================TIMINGS===========================")
	fmt.Printf("Total time passed:\t\t%.2fs\n", statsResult.TotalTimePassed)
	fmt.Printf("Avg time per request:\t\t%.2fms\n", statsResult.AvgTimePerRequest)
	fmt.Printf("Requests per second:\t\t%.2f\n", statsResult.RequestsPerSecond)
	fmt.Printf("Median time per request:\t%.2fms\n", statsResult.MedianTimePerRequest)
	fmt.Printf("99th percentile time:\t\t%.2fms\n", statsResult.N99thPercentileTime)
	fmt.Printf("Slowest time for request:\t%.2fms\n\n", statsResult.SlowestTimeForRequest)
	if statsResult.TotalResponseBodySizes>0{
		fmt.Println("=============================DATA=============================")
		fmt.Printf("Total response body sizes:\t\t%d\n", statsResult.TotalResponseBodySizes)
		fmt.Printf("Avg response body per request:\t\t%.2f Byte\n", statsResult.AvgResponseBodyPerRequest)
		fmt.Printf("Transfer rate per second:\t\t%.2f Byte/s (%.2f MByte/s)\n", statsResult.TransferRateBytePerSecond,statsResult.TransferRateMBytePerSecond)
	}
	fmt.Println("==========================RESPONSES==========================")
	fmt.Printf("ResponseOk:\t\t\t%d (%.2f%%)\n", statsResult.ResponseOk, statsResult.ResponseOkPercentile)
	fmt.Printf("Errors:\t\t\t\t%d (%.2f%%)\n", statsResult.Errors, statsResult.ErrorsPercentile)
}