package rpc

import (
	"errors"
	"sync"
	"time"
)

const (
	//MaxConnsPerHost defines the max conn per host
	MaxConnsPerHost = 8
	//MaxIdleConnsPerHost defines the max idle conn per host
	MaxIdleConnsPerHost = 2
	keepAlive           = time.Minute
)

//Transport defines the struct of transport
type Transport struct {
	mut                 sync.Mutex
	once                sync.Once
	proxys              []*Proxy
	MaxConnsPerHost     int
	MaxIdleConnsPerHost int
	Network             string
	Codec               string
	Options             *Options
}

//NewTransport creates a new transport.
func NewTransport(MaxConnsPerHost int, MaxIdleConnsPerHost int, network, codec string, opts *Options) *Transport {
	if MaxConnsPerHost < 1 {
		MaxConnsPerHost = 1
	}
	if MaxIdleConnsPerHost < 0 {
		MaxIdleConnsPerHost = 0
	}
	t := &Transport{
		MaxConnsPerHost:     MaxConnsPerHost,
		MaxIdleConnsPerHost: MaxIdleConnsPerHost,
		Network:             network,
		Codec:               codec,
		Options:             opts,
	}
	return t
}
func (t *Transport) run() {
	ticker := time.NewTicker(time.Second)
	for range ticker.C {
		func() {
			defer func() {
				if err := recover(); err != nil {
				}
			}()
			t.mut.Lock()
			defer t.mut.Unlock()
			idles := []int{}
			for i, r := range t.proxys {
				r.check()
				if len(r.clients) == 0 {
					idles = append(idles, i)
				}
			}
			for _, i := range idles {
				if len(idles) <= t.MaxIdleConnsPerHost {
					break
				}
				t.proxys = append(t.proxys[:i], t.proxys[i+1:]...)
			}
		}()

	}
}

//GetProxy returns a proxy
func (t *Transport) GetProxy() *Proxy {
	t.once.Do(func() {
		go t.run()
	})
	t.mut.Lock()
	defer t.mut.Unlock()
	if len(t.proxys) < t.MaxConnsPerHost {
		proxy := NewProxy(t.Network, t.Codec, t.Options)
		t.proxys = append(t.proxys, proxy)
		return proxy
	}
	rpcs := t.proxys[0]
	t.proxys = append(t.proxys[1:], rpcs)
	return rpcs
}

//GetConn returns a proxy client
func (t *Transport) GetConn(addr string) *ProxyClient {
	proxy := t.GetProxy()
	if proxy != nil {
		conn := proxy.GetConn(addr)
		if conn != nil {
			return conn
		}
	}
	return nil
}

// Ping is NOT ICMP ping, this is just used to test whether a connection is still alive.
func (t *Transport) Ping(address string) bool {
	proxy := t.GetProxy()
	if proxy != nil {
		return proxy.Ping(address)
	}
	return false
}

// Call invokes the named function, waits for it to complete, and returns its error status.
func (t *Transport) Call(address string, name string, args interface{}, reply interface{}) error {
	proxy := t.GetProxy()
	if proxy != nil {
		return proxy.Call(address, name, args, reply)
	}
	return nil
}

// Go invokes the function asynchronously. It returns the Call structure representing
// the invocation. The done channel will signal when the call is complete by returning
// the same Call object. If done is nil, Go will allocate a new channel.
// If non-nil, done must be buffered or Go will deliberately crash.
func (t *Transport) Go(address string, name string, args interface{}, reply interface{}, done chan *Call) *Call {
	proxy := t.GetProxy()
	if proxy != nil {
		return proxy.Go(address, name, args, reply, done)
	}
	return nil
}

//Proxy defines the struct of proxy.
type Proxy struct {
	mu      sync.RWMutex
	clients map[string]*ProxyClient
	Network string
	Codec   string
	Options *Options
}

//ProxyClient defines the struct of proxy client.
type ProxyClient struct {
	Client
	lastTime  time.Time
	keepAlive time.Duration
}

//NewProxy creates a new proxy.
func NewProxy(network, codec string, opts *Options) *Proxy {
	t := &Proxy{
		clients: make(map[string]*ProxyClient),
		Network: network,
		Codec:   codec,
		Options: opts,
	}
	return t
}

func (t *Proxy) getClients() map[string]*ProxyClient {
	if t.clients == nil {
		t.clients = make(map[string]*ProxyClient)
	}
	return t.clients
}

func (t *Proxy) check() {
	for addr, conn := range t.getClients() {
		if conn.lastTime.Add(conn.keepAlive).Before(time.Now()) {
			t.RemoveConn(addr)
		}
	}
}

//GetConn returns a proxy client
func (t *Proxy) GetConn(address string) *ProxyClient {
	t.mu.Lock()
	defer t.mu.Unlock()
	if _, ok := t.getClients()[address]; ok {
		return t.getClients()[address]
	}
	conn, err := t.NewConn(address)
	if err == nil {
		c := &ProxyClient{conn, time.Now(), keepAlive}
		t.getClients()[address] = c
		return t.getClients()[address]
	}
	return nil
}

//RemoveConn remove a proxy client by address.
func (t *Proxy) RemoveConn(address string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if conn, ok := t.getClients()[address]; ok {
		delete(t.getClients(), address)
		go func(conn Client) {
			time.Sleep(keepAlive)
			conn.Close()
		}(conn)
	}
}

//NewConn creates a client by address.
func (t *Proxy) NewConn(address string) (Client, error) {
	return DialWithOptions(t.Network, address, t.Codec, t.Options)
}

// Ping is NOT ICMP ping, this is just used to test whether a connection is still alive.
func (t *Proxy) Ping(addr string) bool {
	conn := t.GetConn(addr)
	if conn != nil {
		if conn.Ping() {
			return true
		}
	}
	t.RemoveConn(addr)
	return false
}

// Call invokes the named function, waits for it to complete, and returns its error status.
func (t *Proxy) Call(address string, name string, args interface{}, reply interface{}) error {
	conn := t.GetConn(address)
	if conn != nil {
		err := conn.Call(name, args, reply)
		if err != nil {
			t.RemoveConn(address)
			return err
		}
		conn.lastTime = time.Now()
		return nil
	}
	return errors.New("Proxy.Call can not connect to " + address)
}

// Go invokes the function asynchronously. It returns the Call structure representing
// the invocation. The done channel will signal when the call is complete by returning
// the same Call object. If done is nil, Go will allocate a new channel.
// If non-nil, done must be buffered or Go will deliberately crash.
func (t *Proxy) Go(address string, name string, args interface{}, reply interface{}, done chan *Call) *Call {
	conn := t.GetConn(address)
	if conn != nil {
		call := conn.Go(name, args, reply, nil)
		if call.Error != nil {
			t.RemoveConn(address)
			return call
		}
		conn.lastTime = time.Now()
		return call
	}
	return nil
}
