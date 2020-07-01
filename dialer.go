package rpc

import (
	"errors"
)

func Dial(network, address, codec string) (*Client, error) {
	if tran := NewTransport(network); tran != nil {
		conn, err := tran.Dial(address)
		if err != nil {
			return nil, err
		}
		if c := NewCodec(codec); c != nil {
			return NewClient(conn, c), err
		}
		return nil, errors.New("unsupported codec: " + codec)
	}
	return nil, errors.New("unsupported protocol scheme: " + network)

}

func DialWithOptions(address string, opts *Options) (*Client, error) {
	conn, err := opts.Transport.Dial(address)
	if err != nil {
		return nil, err
	}
	if opts.Encoder != nil {
		return NewClientWithEncoder(conn, opts.Encoder), err
	}
	return NewClient(conn, opts.Codec), nil
}
