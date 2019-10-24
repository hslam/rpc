package rpc

import (
	"sync"
	"errors"
	"hslam.com/mgit/Mort/rpc/log"
)

func Dials(total int,network,address,codec string)(*Pool,error){
	p :=  &Pool{
		conns:make([]Client,total),
	}
	for i := 0;i<total;i++{
		conn,err:=Dial(network,address,codec)
		if err != nil {
			return nil,err
		}
		p.conns[i]=conn
	}
	return p,nil
}

func DialsWithMaxRequests(total int,network,address,codec string,max int)(*Pool,error){
	p :=  &Pool{
		conns:make([]Client,total),
	}
	for i := 0;i<total;i++{
		conn,err:=DialWithMaxRequests(network,address,codec,max)
		if err != nil {
			return nil,err
		}
		p.conns[i]=conn
	}
	return p,nil
}

type Pool struct {
	mu 				sync.Mutex
	conns   		[]Client
	pool_id			int64
}

func (p *Pool)conn() (Client) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if len(p.conns) < 1 {
		return nil
	}
	conn := p.conns[0]
	p.conns = append(p.conns[1:], conn)
	return conn
}
func (p *Pool)head() (Client) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if len(p.conns) < 1 {
		return nil
	}
	return p.conns[0]
}
func (p *Pool)All()[]Client{
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.conns
}
func (p *Pool)SetMaxRequests(max int){
	p.mu.Lock()
	defer p.mu.Unlock()
	for _,c:= range p.conns{
		c.SetMaxRequests(max)
	}
}
func (p *Pool)GetMaxRequests()(int){
	p.mu.Lock()
	defer p.mu.Unlock()
	for _,c:= range p.conns{
		return c.GetMaxRequests()
	}
	return -1
}
func (p *Pool)EnablePipelining(){
	p.mu.Lock()
	defer p.mu.Unlock()
	for _,c:= range p.conns{
		c.EnablePipelining()
	}
}
func (p *Pool)EnableMultiplexing(){
	p.mu.Lock()
	defer p.mu.Unlock()
	for _,c:= range p.conns{
		c.EnableMultiplexing()
	}
}
func (p *Pool)EnableBatch(){
	p.mu.Lock()
	defer p.mu.Unlock()
	for _,c:= range p.conns{
		c.EnableBatch()
	}
}
func (p *Pool)EnableBatchAsync(){
	p.mu.Lock()
	defer p.mu.Unlock()
	for _,c:= range p.conns{
		c.EnableBatchAsync()
	}
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
func (p *Pool)GetMaxBatchRequest()int {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _,c:= range p.conns{
		return c.GetMaxBatchRequest()
	}
	return -1
}

func (p *Pool)SetCompressType(compress string){
	p.mu.Lock()
	defer p.mu.Unlock()
	for _,c:= range p.conns{
		c.SetCompressType(compress)
	}
}
func (p *Pool)SetCompressLevel(compress,level string){
	p.mu.Lock()
	defer p.mu.Unlock()
	for _,c:= range p.conns{
		c.SetCompressLevel(compress,level)
	}
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
	return p.conn().Call(name,args,reply)
}
func (p *Pool)CallNoRequest(name string, reply interface{}) ( err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Errorln("CallNoRequest failed:", err)
		}
	}()
	return p.conn().CallNoRequest(name,reply)
}
func (p *Pool)CallNoResponse(name string, args interface{}) ( err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Errorln("CallNoResponse failed:", err)
		}
	}()
	return p.conn().CallNoResponse(name,args)
}
func (p *Pool)OnlyCall(name string) ( err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Errorln("OnlyCall failed:", err)
		}
	}()
	return p.conn().OnlyCall(name)
}
func (p *Pool)Ping() bool {
	defer func() {
		if err := recover(); err != nil {
			log.Errorln("Ping failed:", err)
		}
	}()
	return p.head().Ping()
}
func (p *Pool)DisableRetry() {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _,c:= range p.conns{
		c.DisableRetry()
	}
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