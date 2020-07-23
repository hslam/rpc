package rpc

type ReverseProxy struct {
	TargetAddress string
	Transport     RoundTripper
}

func NewSingleHostReverseProxy(target string) *ReverseProxy {
	return &ReverseProxy{TargetAddress: target}
}

func (c *ReverseProxy) RoundTrip(call *Call) *Call {
	transport := c.Transport
	if transport == nil {
		transport = DefaultTransport
	}
	return transport.RoundTrip(c.TargetAddress, call)
}

func (c *ReverseProxy) Call(serviceMethod string, args interface{}, reply interface{}) error {
	transport := c.Transport
	if transport == nil {
		transport = DefaultTransport
	}
	return transport.Call(c.TargetAddress, serviceMethod, args, reply)
}

func (c *ReverseProxy) Go(serviceMethod string, args interface{}, reply interface{}, done chan *Call) *Call {
	transport := c.Transport
	if transport == nil {
		transport = DefaultTransport
	}
	return transport.Go(c.TargetAddress, serviceMethod, args, reply, done)
}

func (c *ReverseProxy) Ping() error {
	transport := c.Transport
	if transport == nil {
		transport = DefaultTransport
	}
	return transport.Ping(c.TargetAddress)
}
