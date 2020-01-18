package rpc

import "sync"

//Options defines the struct of options.
type Options struct {
	mu               sync.Mutex
	ID               uint64
	MaxRequests      int
	Pipelining       bool
	Multiplexing     bool
	Batching         bool
	MaxBatchRequest  int
	Retry            bool
	CompressType     CompressType
	CompressLevel    CompressLevel
	Timeout          int64
	HeartbeatTimeout int64
	MaxErrPerSecond  int
	MaxErrHeartbeat  int
	NoDelay          bool
	noDelay          bool
}

//DefaultOptions returns a default options.
func DefaultOptions() *Options {
	return &Options{
		ID:               1,
		MaxRequests:      DefaultMaxRequests,
		Pipelining:       false,
		Multiplexing:     true,
		Batching:         false,
		MaxBatchRequest:  DefaultMaxBatchRequest,
		Retry:            true,
		CompressType:     CompressTypeNo,
		CompressLevel:    NoCompression,
		Timeout:          DefaultServerTimeout,
		HeartbeatTimeout: DefaultClientHearbeatTimeout,
		MaxErrPerSecond:  DefaultClientMaxErrPerSecond,
		MaxErrHeartbeat:  DefaultClientMaxErrHeartbeat,
		NoDelay:          false,
		noDelay:          true,
	}
}

//NewOptions creates a new options.
func NewOptions() *Options {
	return DefaultOptions()
}

//Check checks the options.
func (o *Options) Check() {
	o.mu.Lock()
	defer o.mu.Unlock()
	if o.ID < 1 {
		o.ID = 1
	}
	if o.MaxRequests < DefaultMaxRequests {
		o.MaxRequests = DefaultMaxRequests
	}
	if o.MaxBatchRequest < DefaultMaxBatchRequest {
		o.MaxBatchRequest = DefaultMaxBatchRequest
	}
	if o.Pipelining && o.Multiplexing {
		o.Multiplexing = false
	}
	if o.Pipelining || o.Multiplexing {
		o.noDelay = false
	}
	if o.Batching {
		o.noDelay = true
	}
	if o.NoDelay {
		o.noDelay = true
	}
}

//SetMaxRequests sets the number of max requests.
func (o *Options) SetMaxRequests(max int) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.MaxRequests = max
}

//GetMaxRequests returns the number of max requests.
func (o *Options) GetMaxRequests() int {
	o.mu.Lock()
	defer o.mu.Unlock()
	return o.MaxRequests
}

//SetPipelining enables pipelining.
func (o *Options) SetPipelining(enable bool) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.Pipelining = enable
	o.Multiplexing = !enable
}

//SetMultiplexing enables multiplexing.
func (o *Options) SetMultiplexing(enable bool) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.Multiplexing = enable
	o.Pipelining = !enable
}

//SetBatching enables batching.
func (o *Options) SetBatching(enable bool) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.Batching = enable
}

//SetMaxBatchRequest sets the number of max batch requests.
func (o *Options) SetMaxBatchRequest(max int) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.MaxBatchRequest = max
}

//GetMaxBatchRequest returns the number of max batch requests.
func (o *Options) GetMaxBatchRequest() int {
	o.mu.Lock()
	defer o.mu.Unlock()
	return o.MaxBatchRequest
}

//SetCompressType sets the compress type.
func (o *Options) SetCompressType(compress string) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.CompressType = getCompressType(compress)
}

//SetCompressLevel sets the compress level.
func (o *Options) SetCompressLevel(level string) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.CompressLevel = getCompressLevel(level)
}

//SetID sets the client ID.
func (o *Options) SetID(id uint64) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.ID = id
}

//GetID returns the client ID.
func (o *Options) GetID() uint64 {
	o.mu.Lock()
	defer o.mu.Unlock()
	return o.ID
}

//SetTimeout sets the request timeout.
func (o *Options) SetTimeout(timeout int64) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.Timeout = timeout
}

//GetTimeout returns the request timeout.
func (o *Options) GetTimeout() int64 {
	o.mu.Lock()
	defer o.mu.Unlock()
	return o.Timeout
}

//SetHeartbeatTimeout sets the heartbeat timeout.
func (o *Options) SetHeartbeatTimeout(timeout int64) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.HeartbeatTimeout = timeout
}

//GetHeartbeatTimeout returns the heartbeat timeout.
func (o *Options) GetHeartbeatTimeout() int64 {
	o.mu.Lock()
	defer o.mu.Unlock()
	return o.HeartbeatTimeout
}

//SetMaxErrPerSecond sets the number of max error per second.
func (o *Options) SetMaxErrPerSecond(maxErrPerSecond int) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.MaxErrPerSecond = maxErrPerSecond
}

//GetMaxErrPerSecond returns the number of max error per second.
func (o *Options) GetMaxErrPerSecond() int {
	o.mu.Lock()
	defer o.mu.Unlock()
	return o.MaxErrPerSecond
}

//SetMaxErrHeartbeat sets the number of max hearbeat error.
func (o *Options) SetMaxErrHeartbeat(maxErrHeartbeat int) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.MaxErrHeartbeat = maxErrHeartbeat
}

//GetMaxErrHeartbeat returns the number of max hearbeat error.
func (o *Options) GetMaxErrHeartbeat() int {
	o.mu.Lock()
	defer o.mu.Unlock()
	return o.MaxErrHeartbeat
}

//SetRetry enables retry.
func (o *Options) SetRetry(enable bool) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.Retry = enable
}

// SetNoDelay controls whether the operating system should delay
// packet transmission in hopes of sending fewer packets (Nagle's
// algorithm).  The default is true (no delay), meaning that data is
// sent as soon as possible after a Write.
func (o *Options) SetNoDelay(enabled bool) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.NoDelay = enabled
}
