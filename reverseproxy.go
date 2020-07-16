package rpc

type ReverseProxy struct {
	TargetAddress string
	Transport     RoundTripper
}

func NewSingleHostReverseProxy(addr string) *ReverseProxy {
	return &ReverseProxy{TargetAddress: addr}
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
