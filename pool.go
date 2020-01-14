package rpc

import (
	"sync"
)

func Dials(total int, network, address, codec string) (*Pool, error) {
	p := &Pool{
		conns: make([]Client, total),
	}
	for i := 0; i < total; i++ {
		conn, err := Dial(network, address, codec)
		if err != nil {
			return nil, err
		}
		p.conns[i] = conn
	}
	return p, nil
}

func DialsWithOptions(total int, network, address, codec string, opts *Options) (*Pool, error) {
	p := &Pool{
		conns: make([]Client, total),
	}
	for i := 0; i < total; i++ {
		conn, err := DialWithOptions(network, address, codec, opts)
		if err != nil {
			return nil, err
		}
		p.conns[i] = conn
	}
	return p, nil
}

type Pool struct {
	mu      sync.Mutex
	conns   []Client
	pool_id int64
}

func (p *Pool) conn() Client {
	p.mu.Lock()
	defer p.mu.Unlock()
	if len(p.conns) < 1 {
		return nil
	}
	conn := p.conns[0]
	p.conns = append(p.conns[1:], conn)
	return conn
}
func (p *Pool) head() Client {
	p.mu.Lock()
	defer p.mu.Unlock()
	if len(p.conns) < 1 {
		return nil
	}
	return p.conns[0]
}
func (p *Pool) All() []Client {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.conns
}
func (p *Pool) GetMaxRequests() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, c := range p.conns {
		return c.GetMaxRequests()
	}
	return -1
}

func (p *Pool) GetMaxBatchRequest() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, c := range p.conns {
		return c.GetMaxBatchRequest()
	}
	return -1
}

func (p *Pool) GetID() int64 {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.pool_id
}

func (p *Pool) GetTimeout() int64 {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, c := range p.conns {
		return c.GetTimeout()
	}
	return -1
}

func (p *Pool) GetHeartbeatTimeout() int64 {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, c := range p.conns {
		return c.GetHeartbeatTimeout()
	}
	return -1
}

func (p *Pool) GetMaxErrHeartbeat() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, c := range p.conns {
		return c.GetMaxErrHeartbeat()
	}
	return -1
}

func (p *Pool) GetMaxErrPerSecond() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, c := range p.conns {
		return c.GetMaxErrPerSecond()
	}
	return -1
}

func (p *Pool) CodecName() string {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, c := range p.conns {
		return c.CodecName()
	}
	return ""
}
func (p *Pool) CodecType() CodecType {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, c := range p.conns {
		return c.CodecType()
	}
	return FUNCS_CODEC_INVALID
}
func (p *Pool) Go(name string, args interface{}, reply interface{}, done chan *Call) *Call {
	defer func() {
		if err := recover(); err != nil {
			Errorln("Go failed:", err)
		}
	}()
	return p.conn().Go(name, args, reply, done)
}
func (p *Pool) Call(name string, args interface{}, reply interface{}) (err error) {
	defer func() {
		if err := recover(); err != nil {
			Errorln("Call failed:", err)
		}
	}()
	return p.conn().Call(name, args, reply)
}
func (p *Pool) CallNoRequest(name string, reply interface{}) (err error) {
	defer func() {
		if err := recover(); err != nil {
			Errorln("CallNoRequest failed:", err)
		}
	}()
	return p.conn().CallNoRequest(name, reply)
}
func (p *Pool) CallNoResponse(name string, args interface{}) (err error) {
	defer func() {
		if err := recover(); err != nil {
			Errorln("CallNoResponse failed:", err)
		}
	}()
	return p.conn().CallNoResponse(name, args)
}
func (p *Pool) OnlyCall(name string) (err error) {
	defer func() {
		if err := recover(); err != nil {
			Errorln("OnlyCall failed:", err)
		}
	}()
	return p.conn().OnlyCall(name)
}
func (p *Pool) Ping() bool {
	defer func() {
		if err := recover(); err != nil {
			Errorln("Ping failed:", err)
		}
	}()
	return p.head().Ping()
}

func (p *Pool) Close() (err error) {
	for _, c := range p.conns {
		err = c.Close()
	}
	return err
}
func (p *Pool) Closed() bool {
	for _, c := range p.conns {
		return c.Closed()
	}
	return false
}
