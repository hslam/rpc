package rpc

import "sync"

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
		MaxErrHeartbeat:  DefaultClientMaxErrHearbeat,
		NoDelay:          false,
		noDelay:          true,
	}
}
func NewOptions() *Options {
	return DefaultOptions()
}
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
func (o *Options) SetMaxRequests(max int) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.MaxRequests = max
}
func (o *Options) GetMaxRequests() int {
	o.mu.Lock()
	defer o.mu.Unlock()
	return o.MaxRequests
}
func (o *Options) SetPipelining(enable bool) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.Pipelining = enable
	o.Multiplexing = !enable
}

func (o *Options) SetMultiplexing(enable bool) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.Multiplexing = enable
	o.Pipelining = !enable
}

func (o *Options) SetBatching(enable bool) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.Batching = enable
}

func (o *Options) SetMaxBatchRequest(max int) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.MaxBatchRequest = max
}
func (o *Options) GetMaxBatchRequest() int {
	o.mu.Lock()
	defer o.mu.Unlock()
	return o.MaxBatchRequest
}

func (o *Options) SetCompressType(compress string) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.CompressType = getCompressType(compress)
}
func (o *Options) SetCompressLevel(level string) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.CompressLevel = getCompressLevel(level)
}
func (o *Options) SetID(id uint64) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.ID = id
}
func (o *Options) GetID() uint64 {
	o.mu.Lock()
	defer o.mu.Unlock()
	return o.ID
}
func (o *Options) SetTimeout(timeout int64) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.Timeout = timeout
}
func (o *Options) GetTimeout() int64 {
	o.mu.Lock()
	defer o.mu.Unlock()
	return o.Timeout
}

func (o *Options) SetHeartbeatTimeout(timeout int64) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.HeartbeatTimeout = timeout
}
func (o *Options) GetHeartbeatTimeout() int64 {
	o.mu.Lock()
	defer o.mu.Unlock()
	return o.HeartbeatTimeout
}
func (o *Options) SetMaxErrPerSecond(maxErrPerSecond int) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.MaxErrPerSecond = maxErrPerSecond
}
func (o *Options) GetMaxErrPerSecond() int {
	o.mu.Lock()
	defer o.mu.Unlock()
	return o.MaxErrPerSecond
}
func (o *Options) SetMaxErrHeartbeat(maxErrHeartbeat int) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.MaxErrHeartbeat = maxErrHeartbeat
}
func (o *Options) GetMaxErrHeartbeat() int {
	o.mu.Lock()
	defer o.mu.Unlock()
	return o.MaxErrHeartbeat
}
func (o *Options) SetRetry(enable bool) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.Retry = enable
}
func (o *Options) SetNoDelay(enabled bool) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.NoDelay = enabled
}
