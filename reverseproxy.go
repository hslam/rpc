// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package rpc

// ReverseProxy is an RPC Handler that takes an incoming request and
// sends it to another server, proxying the response back to the
// client.
type ReverseProxy struct {
	TargetAddress string
	Transport     RoundTripper
}

// NewSingleHostReverseProxy returns a new ReverseProxy that routes
// requests to the target.
func NewSingleHostReverseProxy(target string) *ReverseProxy {
	return &ReverseProxy{TargetAddress: target}
}

// RoundTrip executes a single RPC transaction, returning
// a Response for the provided Request.
func (c *ReverseProxy) RoundTrip(call *Call) *Call {
	transport := c.Transport
	if transport == nil {
		transport = DefaultTransport
	}
	return transport.RoundTrip(c.TargetAddress, call)
}

// Call invokes the named function, waits for it to complete, and returns its error status.
func (c *ReverseProxy) Call(serviceMethod string, args interface{}, reply interface{}) error {
	transport := c.Transport
	if transport == nil {
		transport = DefaultTransport
	}
	return transport.Call(c.TargetAddress, serviceMethod, args, reply)
}

// Go invokes the function asynchronously. It returns the Call structure representing
// the invocation. The done channel will signal when the call is complete by returning
// the same Call object. If done is nil, Go will allocate a new channel.
// If non-nil, done must be buffered or Go will deliberately crash.
func (c *ReverseProxy) Go(serviceMethod string, args interface{}, reply interface{}, done chan *Call) *Call {
	transport := c.Transport
	if transport == nil {
		transport = DefaultTransport
	}
	return transport.Go(c.TargetAddress, serviceMethod, args, reply, done)
}

// Ping is NOT ICMP ping, this is just used to test whether a connection is still alive.
func (c *ReverseProxy) Ping() error {
	transport := c.Transport
	if transport == nil {
		transport = DefaultTransport
	}
	return transport.Ping(c.TargetAddress)
}
