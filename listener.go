package rpc

import (
	"errors"
	"os"
)

func Listen(network, address string, codec string) error {
	if tran := NewTransport(network); tran != nil {
		logger.Noticef("pid - %d", os.Getpid())
		logger.Noticef("network - %s", tran.Scheme())
		logger.Noticef("listening on %s", address)
		lis, err := tran.Listen(address)
		if err != nil {
			return err
		}
		if c := NewCodec(codec); c != nil {
			for {
				conn, err := lis.Accept()
				if err != nil {
					continue
				}
				go ServeConn(conn, c, nil)
			}
			return nil
		}
		return errors.New("unsupported codec: " + codec)
	}
	return errors.New("unsupported protocol scheme: " + network)
}

func ListenWithOptions(address string, opts *Options) error {
	if opts.Codec == nil && opts.Encoder == nil {
		return errors.New("need opts.Codec or opts.Encoder")
	}
	logger.Noticef("pid - %d", os.Getpid())
	logger.Noticef("network - %s", opts.Transport.Scheme())
	logger.Noticef("listening on %s", address)
	lis, err := opts.Transport.Listen(address)
	if err != nil {
		return err
	}
	for {
		conn, err := lis.Accept()
		if err != nil {
			continue
		}
		go ServeConn(conn, opts.Codec, opts.Encoder)
	}
}
