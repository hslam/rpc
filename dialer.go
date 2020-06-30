package rpc

import (
	"github.com/hslam/transport"
)

func Dial(tran transport.Transport, address string, codec *Codec) (*Client, error) {
	conn, err := tran.Dial(address)
	if err != nil {
		return nil, err
	}
	return NewClient(conn, codec), err
}
