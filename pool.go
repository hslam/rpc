package rpc

import (
	"sync"
	"errors"
)

type ConnPool chan Conn

type Pool struct {
	mu 				sync.Mutex
	connPool 		ConnPool
	conns   		[]Conn
	pool_id			int64
}

func (p *Pool)Get()Conn{
	p.mu.Lock()
	defer p.mu.Unlock()
	c:=<-p.connPool
	return c
}
func (p *Pool)Put(c Conn){
	p.mu.Lock()
	defer p.mu.Unlock()
	p.connPool<-c
}
func (p *Pool)All()[]Conn{
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.conns
}
func (p *Pool)EnabledBatch(){
	p.mu.Lock()
	defer p.mu.Unlock()
	for _,c:= range p.conns{
		c.EnabledBatch()
	}
}
func (p *Pool)GetMaxBatchRequest()int {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _,c:= range p.conns{
		return c.GetMaxBatchRequest()
	}
	return -1
}
func (p *Pool)SetMaxBatchRequest(maxBatchRequest int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _,c:= range p.conns{
		c.SetMaxBatchRequest(maxBatchRequest)
	}
}
func (p *Pool)GetMaxConcurrentRequest()(int){
	for _,c:= range p.conns{
		return c.GetMaxConcurrentRequest()
	}
	return -1
}
func (p *Pool)SetID(id int64)error{
	p.mu.Lock()
	defer p.mu.Unlock()
	var i int64=0
	p.pool_id=id
	for _,c:= range p.conns{
		c.SetID(id+i)
		i++
	}
	return nil
}
func (p *Pool)GetID()int64{
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.pool_id
}
func (p *Pool)CodecName()string {
	for _,c:= range p.conns{
		return c.CodecName()
	}
	return ""
}
func (p *Pool)CodecType()CodecType {
	for _,c:= range p.conns{
		return c.CodecType()
	}
	return FUNCS_CODEC_INVALID
}
func (p *Pool)Call(name string, args interface{}, reply interface{}) ( err error) {
	defer func() {
		if err := recover(); err != nil {
			Panicln("Call failed:", err)
		}
	}()
	c:=<-p.connPool
	defer func(c Conn) {
		p.connPool<-c
	}(c)
	return c.Call(name,args,reply)
}
func (p *Pool)CallNoResponse(name string, args interface{}) ( err error) {
	defer func() {
		if err := recover(); err != nil {
			Panicln("Call failed:", err)
		}
	}()
	c:=<-p.connPool
	defer func(c Conn) {
		p.connPool<-c
	}(c)
	return c.CallNoResponse(name,args)
}
func (p *Pool)RemoteCall(b []byte)([]byte,error){
	return nil, errors.New("not suportted")
}
func (p *Pool)RemoteCallNoResponse(b []byte)(error){
	return errors.New("not suportted")
}
func (p *Pool)Close() ( err error) {
	for _,c:= range p.conns{
		err=c.Close()
	}
	return err
}