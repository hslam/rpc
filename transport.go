package rpc

import (
	"errors"
	"sync"
	"time"
)

const (
	MaxConnsPerHost=8
	keepAlive=time.Minute
)

type Transport struct {
	mut			sync.Mutex
	proxys		[]*Proxy
	MaxConnsPerHost int
}

func NewTransport(maxConnsPerHost int,network,codec string,opts *Options)*Transport {
	t:=&Transport{
		MaxConnsPerHost:maxConnsPerHost,
	}
	for i:=0;i<maxConnsPerHost ;i++  {
		t.proxys=append(t.proxys,NewProxy(network,codec,opts) )
	}
	go t.run()
	return t
}
func (t *Transport)run(){
	ticker:=time.NewTicker(time.Second)
	for range ticker.C {
		for _,r:=range t.proxys{
			r.check()
		}
	}
}
func (t *Transport)GetProxy() (*Proxy) {
	t.mut.Lock()
	defer t.mut.Unlock()
	if len(t.proxys) < 1 {
		return nil
	}
	rpcs := t.proxys[0]
	t.proxys = append(t.proxys[1:], rpcs)
	return rpcs
}
func (t *Transport)GetConn(addr string) (*ProxyClient) {
	proxy:=t.GetProxy()
	if proxy!=nil{
		conn:=proxy.GetConn(addr)
		if conn!=nil{
			return conn
		}
	}
	return nil
}
func (t *Transport)Ping(address string) (bool) {
	proxy:=t.GetProxy()
	if proxy!=nil{
		return proxy.Ping(address)
	}
	return false
}
func (t *Transport)Call(name string,args interface{}, reply interface{},address string) (error) {
	proxy:=t.GetProxy()
	if proxy!=nil{
		return proxy.Call(name,args,reply,address)
	}
	return nil
}
func (t *Transport)Go(name string, args interface{}, reply interface{}, done chan *Call,address string) *Call {
	proxy:=t.GetProxy()
	if proxy!=nil{
		return proxy.Go(name,args,reply,done,address)
	}
	return nil
}


type Proxy struct {
	mu				sync.RWMutex
	clients			map[string]*ProxyClient
	Network			string
	Codec		 	string
	Options			*Options
}
type ProxyClient struct {
	Client
	lastTime 				time.Time
	keepAlive				time.Duration
}
func NewProxy(network,codec string,opts *Options) *Proxy{
	t :=&Proxy{
		clients:make(map[string]*ProxyClient),
		Network:network,
		Codec:codec,
		Options:opts,
	}
	return t
}
func (t *Proxy) getClients()map[string]*ProxyClient{
	if t.clients==nil{
		t.clients=make(map[string]*ProxyClient)
	}
	return t.clients
}
func (t *Proxy) check(){
	for addr,conn:=range t.getClients(){
		if conn.lastTime.Add(conn.keepAlive).Before(time.Now()){
			t.RemoveConn(addr)
		}
	}
}
func (t *Proxy) GetConn(address string) *ProxyClient {
	t.mu.Lock()
	defer t.mu.Unlock()
	if _,ok:=t.getClients()[address];ok{
		return t.getClients()[address]
	}
	conn, err := t.NewConn(address)
	if err==nil{
		c:=&ProxyClient{conn,time.Now(),keepAlive}
		t.getClients()[address] = c
		return t.getClients()[address]
	}
	return nil
}

func (t *Proxy) RemoveConn(address string){
	t.mu.Lock()
	defer t.mu.Unlock()
	if conn,ok:=t.getClients()[address];ok{
		delete(t.getClients(),address)
		go func(conn Client) {
			time.Sleep(keepAlive)
			conn.Close()
		}(conn)
	}
}
func (t *Proxy) NewConn(address string) (Client, error){
	return DialWithOptions(t.Network,address,t.Codec,t.Options)
}
func (t *Proxy) Ping(addr string)bool {
	conn:= t.GetConn(addr)
	if conn!=nil{
		if conn.Ping(){
			return true
		}
	}
	t.RemoveConn(addr)
	return false
}

func (t *Proxy) Call(name string,args interface{}, reply interface{},addr string)error {
	conn:= t.GetConn(addr)
	if conn!=nil{
		err:=conn.Call(name,args,reply)
		if err!=nil{
			t.RemoveConn(addr)
			return err
		}
		conn.lastTime=time.Now()
		return nil
	}
	return errors.New("Proxy.Call can not connect to "+addr)
}

func (t *Proxy) Go(name string, args interface{}, reply interface{}, done chan *Call,addr string) *Call {
	conn:= t.GetConn(addr)
	if conn!=nil{
		call:=conn.Go(name,args,reply,nil)
		if call.Error!=nil{
			t.RemoveConn(addr)
			return call
		}
		conn.lastTime=time.Now()
		return call
	}
	return nil
}
