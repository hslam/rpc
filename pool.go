package rpc

import (
	"sync"
)

// Dials connects to an RPC server at the specified network address codec
// and returns a pool of clients.
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

// DialsWithOptions connects to an RPC server at the specified network address codec
// and returns a pool of clients.
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

//Pool defines the set of clients.
type Pool struct {
	mu    sync.Mutex
	conns []Client
	ID    uint64
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

//All returns all clients.
func (p *Pool) All() []Client {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.conns
}

//GetMaxRequests returns the number of max requests.
func (p *Pool) GetMaxRequests() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, c := range p.conns {
		return c.GetMaxRequests()
	}
	return -1
}

//GetMaxBatchRequest returns the number of max batch requests.
func (p *Pool) GetMaxBatchRequest() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, c := range p.conns {
		return c.GetMaxBatchRequest()
	}
	return -1
}

//GetID returns a pool id.
func (p *Pool) GetID() uint64 {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.ID
}

//GetTimeout returns the request timeout.
func (p *Pool) GetTimeout() int64 {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, c := range p.conns {
		return c.GetTimeout()
	}
	return -1
}

//GetHeartbeatTimeout returns the heartbeat timeout.
func (p *Pool) GetHeartbeatTimeout() int64 {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, c := range p.conns {
		return c.GetHeartbeatTimeout()
	}
	return -1
}

//GetMaxErrHeartbeat returns the number of max heartbeat errors.
func (p *Pool) GetMaxErrHeartbeat() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, c := range p.conns {
		return c.GetMaxErrHeartbeat()
	}
	return -1
}

//GetMaxErrPerSecond returns the number of max errors per second.
func (p *Pool) GetMaxErrPerSecond() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, c := range p.conns {
		return c.GetMaxErrPerSecond()
	}
	return -1
}

//CodecName returns the codec name.
func (p *Pool) CodecName() string {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, c := range p.conns {
		return c.CodecName()
	}
	return ""
}

//CodecType returns the codec type.
func (p *Pool) CodecType() CodecType {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, c := range p.conns {
		return c.CodecType()
	}
	return FuncsCodecINVALID
}

// Go invokes the function asynchronously. It returns the Call structure representing
// the invocation. The done channel will signal when the call is complete by returning
// the same Call object. If done is nil, Go will allocate a new channel.
// If non-nil, done must be buffered or Go will deliberately crash.
func (p *Pool) Go(name string, args interface{}, reply interface{}, done chan *Call) *Call {
	defer func() {
		if err := recover(); err != nil {
			logger.Errorln("Go failed:", err)
		}
	}()
	return p.conn().Go(name, args, reply, done)
}

// Call invokes the named function, waits for it to complete, and returns its error status.
func (p *Pool) Call(name string, args interface{}, reply interface{}) (err error) {
	defer func() {
		if err := recover(); err != nil {
			logger.Errorln("Call failed:", err)
		}
	}()
	return p.conn().Call(name, args, reply)
}

// CallNoRequest invokes the named function but doesn't use args, waits for it to complete, and returns its error status.
func (p *Pool) CallNoRequest(name string, reply interface{}) (err error) {
	defer func() {
		if err := recover(); err != nil {
			logger.Errorln("CallNoRequest failed:", err)
		}
	}()
	return p.conn().CallNoRequest(name, reply)
}

// CallNoResponse invokes the named function but doesn't return reply, waits for it to complete, and returns its error status.
func (p *Pool) CallNoResponse(name string, args interface{}) (err error) {
	defer func() {
		if err := recover(); err != nil {
			logger.Errorln("CallNoResponse failed:", err)
		}
	}()
	return p.conn().CallNoResponse(name, args)
}

// OnlyCall invokes the named function but doesn't use args and doesn't return reply, waits for it to complete, and returns its error status.
func (p *Pool) OnlyCall(name string) (err error) {
	defer func() {
		if err := recover(); err != nil {
			logger.Errorln("OnlyCall failed:", err)
		}
	}()
	return p.conn().OnlyCall(name)
}

// Ping is NOT ICMP ping, this is just used to test whether a connection is still alive.
func (p *Pool) Ping() bool {
	defer func() {
		if err := recover(); err != nil {
			logger.Errorln("Ping failed:", err)
		}
	}()
	return p.head().Ping()
}

// Close closes the connection
func (p *Pool) Close() (err error) {
	for _, c := range p.conns {
		err = c.Close()
	}
	return err
}

// Closed returns the closed
func (p *Pool) Closed() bool {
	for _, c := range p.conns {
		return c.Closed()
	}
	return false
}
