package rpc

import (
	"hslam.com/mgit/Mort/workerpool"
)
var (
	workerPoolSize			=1024
	workerMax				=64
	useWorkerPool			=false
	workerPool *workerpool.Pool
)

func EnabledWorkerPool() {
	useWorkerPool=true
	workerPool= workerpool.New(workerPoolSize,workerMax)
}
func EnabledWorkerPoolWithSize(size ,max int) {
	useWorkerPool=true
	workerPoolSize=size
	workerMax=max
	workerPool= workerpool.New(workerPoolSize,workerMax)
}

