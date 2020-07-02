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
			return NewClient(conn, c, nil), err
		}
		return nil, errors.New("unsupported codec: " + codec)
	}
	return nil, errors.New("unsupported protocol scheme: " + network)

}

func DialWithOptions(address string, opts *Options) (*Client, error) {
	if opts.Codec == nil && opts.Encoder == nil {
		return nil, errors.New("need opts.Codec or opts.Encoder")
	}
	conn, err := opts.Transport.Dial(address)
	if err != nil {
		return nil, err
	}
	return NewClient(conn, opts.Codec, opts.Encoder), nil
}
