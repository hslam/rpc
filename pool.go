package rpc

import (
	"sync"
	"errors"
	"hslam.com/mgit/Mort/rpc/log"
)

func Dials(total int,network,address,codec string)(*Pool,error){
	p :=  &Pool{
		connPool:make(ClientPool,total),
		conns:make([]*Client,total),
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

type ClientPool chan *Client

type Pool struct {
	mu 				sync.Mutex
	connPool 		ClientPool
	conns   		[]*Client
	pool_id			int64
}

func (p *Pool)Get()*Client{
	p.mu.Lock()
	defer p.mu.Unlock()
	c:=<-p.connPool
	return c
}
func (p *Pool)Put(c *Client){
	p.mu.Lock()
	defer p.mu.Unlock()
	p.connPool<-c
}
func (p *Pool)All()[]*Client{
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
func (p *Pool)SetMaxBatchRequest(maxBatchRequest int) error{
	p.mu.Lock()
	defer p.mu.Unlock()
	for _,c:= range p.conns{
		err:=c.SetMaxBatchRequest(maxBatchRequest)
		if err!=nil{
			return err
		}
	}
	return nil
}
func (p *Pool)GetMaxConcurrentRequest()(int){
	p.mu.Lock()
	defer p.mu.Unlock()
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
		err:=c.SetID(id+i)
		if err!=nil{
			return err
		}
		i++
	}
	return nil
}
func (p *Pool)GetID()int64{
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.pool_id
}
func (p *Pool)SetTimeout(timeout int64)error{
	p.mu.Lock()
	defer p.mu.Unlock()
	for _,c:= range p.conns{
		err:=c.SetTimeout(timeout)
		if err!=nil{
			return err
		}
	}
	return nil
}
func (p *Pool)GetTimeout()int64{
	p.mu.Lock()
	defer p.mu.Unlock()
	for _,c:= range p.conns{
		return c.GetTimeout()
	}
	return -1
}

func (p *Pool)SetHeartbeatTimeout(timeout int64)error{
	p.mu.Lock()
	defer p.mu.Unlock()
	for _,c:= range p.conns{
		err:=c.SetHeartbeatTimeout(timeout)
		if err!=nil{
			return err
		}
	}
	return nil
}
func (p *Pool)GetHeartbeatTimeout()int64{
	p.mu.Lock()
	defer p.mu.Unlock()
	for _,c:= range p.conns{
		return c.GetHeartbeatTimeout()
	}
	return -1
}


func (p *Pool)SetMaxErrHeartbeat(maxErrHeartbeat int)error{
	p.mu.Lock()
	defer p.mu.Unlock()
	for _,c:= range p.conns{
		err:=c.SetMaxErrHeartbeat(maxErrHeartbeat)
		if err!=nil{
			return err
		}
	}
	return nil
}
func (p *Pool)GetMaxErrHeartbeat()int{
	p.mu.Lock()
	defer p.mu.Unlock()
	for _,c:= range p.conns{
		return c.GetMaxErrHeartbeat()
	}
	return -1
}

func (p *Pool)SetMaxErrPerSecond(maxErrPerSecond int)error{
	p.mu.Lock()
	defer p.mu.Unlock()
	for _,c:= range p.conns{
		err:=c.SetMaxErrPerSecond(maxErrPerSecond)
		if err!=nil{
			return err
		}
	}
	return nil
}
func (p *Pool)GetMaxErrPerSecond()int{
	p.mu.Lock()
	defer p.mu.Unlock()
	for _,c:= range p.conns{
		return c.GetMaxErrPerSecond()
	}
	return -1
}

func (p *Pool)CodecName()string {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _,c:= range p.conns{
		return c.CodecName()
	}
	return ""
}
func (p *Pool)CodecType()CodecType {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _,c:= range p.conns{
		return c.CodecType()
	}
	return FUNCS_CODEC_INVALID
}
func (p *Pool)Call(name string, args interface{}, reply interface{}) ( err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Errorln("Call failed:", err)
		}
	}()
	c:=<-p.connPool
	defer func(c *Client) {
		p.connPool<-c
	}(c)
	return c.Call(name,args,reply)
}
func (p *Pool)CallNoResponse(name string, args interface{}) ( err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Errorln("Call failed:", err)
		}
	}()
	c:=<-p.connPool
	defer func(c *Client) {
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
func (p *Pool)Closed()bool {
	for _,c:= range p.conns{
		return c.Closed()
	}
	return false
}