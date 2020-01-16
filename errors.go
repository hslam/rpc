package rpc

import (
	"errors"
)

var (
	//ErrConnExit defines the error of conn exit
	ErrConnExit = errors.New("exit")
	//ErrRPCNoResponse defines the error of no response
	ErrRPCNoResponse = errors.New("RPC NoResponse")
	//ErrHystrix defines the error of hystrix
	ErrHystrix = errors.New("Hystrix")
	//ErrSetClientID defines the error of SetClientID
	ErrSetClientID = errors.New("0<=ClientID<=1023")
	//ErrSetMaxErrPerSecond defines the error of SetMaxErrPerSecond
	ErrSetMaxErrPerSecond = errors.New("maxErrPerSecond>0")
	//ErrSetMaxErrHeartbeat defines the error of SetMaxErrHeartbeat
	ErrSetMaxErrHeartbeat = errors.New("maxErrHeartbeat>0")
	//ErrSetMaxBatchRequest defines the error of SetMaxBatchRequest
	ErrSetMaxBatchRequest = errors.New("maxBatchRequest>0")
	//ErrSetTimeout defines the error of SetTimeout
	ErrSetTimeout = errors.New("timeout>0")
	//ErrRemoteCall defines the error of RemoteCall
	ErrRemoteCall = errors.New("RemoteCall cbChan is close")
	//ErrTimeOut defines the error of time out
	ErrTimeOut = errors.New("time out")
	//ErrReqID defines the error of ReqID
	ErrReqID = errors.New("reqId err")
	//ErrClientID defines the error of ClientID
	ErrClientID = errors.New("clientID err")
	//ErrShutdown defines the error of shut down
	ErrShutdown = errors.New("connection is shut down")
)
