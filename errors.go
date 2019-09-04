package rpc


import (
	"errors"
)

var (
	ErrConnExit=errors.New("exit")
	RPCConnNoResponse=errors.New("RPC NoResponse")
	ErrHystrix=errors.New("Hystrix")
	ErrSetClientID=errors.New("0<=ClientID<=1023")
	ErrSetMaxErrPerSecond=errors.New("maxErrPerSecond>0")
	ErrSetMaxErrHeartbeat=errors.New("maxErrHeartbeat>0")
	ErrSetMaxBatchRequest=errors.New("maxBatchRequest>0")
	ErrSetTimeout=errors.New("timeout>0")
	ErrRemoteCall=errors.New("RemoteCall cbChan is close")
	ErrTimeOut=errors.New("time out")
	ErrReqId=errors.New("req_id err")
	ErrClientId=errors.New("client_id err")
)
