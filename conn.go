package rpc

import (
	"errors"
)

//Dialer
type Conn interface {
	Call(name string, args interface{}, reply interface{}) (error)
	CallNoResponse(name string, args interface{}) (error)
	RemoteCall(b []byte)([]byte, error)
	RemoteCallNoResponse(b []byte)(error)
	EnabledBatch()
	GetMaxConcurrentRequest()(int)
	GetMaxBatchRequest()(int)
	SetMaxBatchRequest(maxConcurrentRequest int)error
	SetID(id int64)error
	GetID()int64
	SetTimeout(timeout int64)error
	GetTimeout()int64
	SetHeartbeatTimeout(timeout int64)error
	GetHeartbeatTimeout()int64
	SetMaxErrHeartbeat(maxErrHeartbeat int)error
	GetMaxErrHeartbeat()int
	SetMaxErrPerSecond(maxErrPerSecond int)error
	GetMaxErrPerSecond()int
	CodecName()(string)
	CodecType()(CodecType)
	Close()(error)
	Closed()bool
}

func Dial(network,address,codec string) (Conn, error) {
	var transporter	Transporter
	var err error
	switch network {
	case TCP:
		transporter,err=NewTCPTransporter(address)
	case UDP:
		transporter,err=NewUDPTransporter(address)
	case QUIC:
		transporter,err=NewQUICTransporter(address)
	case WS:
		transporter,err=NewWSTransporter(address)
	case FASTHTTP:
		transporter,err=NewFASTHTTPTransporter(address)
	case HTTP:
		transporter,err=NewHTTPTransporter(address)
	case HTTP2:
		transporter,err=NewHTTP2Transporter(address)
	default:
		return nil, errors.New("this network is not suported")
	}
	if err!=nil{
		return nil, errors.New("init transporter err")
	}
	return NewClient(transporter,codec)
}

func Dials(total int,network,address,codec string)(*Pool,error){
	p :=  &Pool{
		connPool:make(ConnPool,total),
		conns:make([]Conn,total),
	}
	for i := 0;i<total;i++{
		conn,err:=Dial(network,address,codec)
		if err != nil {
			return nil,err
		}
		p.connPool <- conn
		p.conns[i]=conn
	}
	return p,nil
}