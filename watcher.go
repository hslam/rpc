package rpc

import (
	"sync/atomic"
)

// Watcher represents a watcher.
type Watcher interface {
	// Wait will return value when the key is triggered.
	Wait() ([]byte, error)
	// Stop stops the watch.
	Stop() error
}

type watcher struct {
	client *Client
	C      chan *watcher
	key    string
	Value  []byte
	Error  error
	done   chan struct{}
	closed uint32
}

func (w *watcher) trigger(value []byte, err error) {
	w.Value = value
	w.Error = err
	select {
	case w.C <- w:
	default:
	}
}

func (w *watcher) Wait() ([]byte, error) {
	select {
	case <-w.C:
		return w.Value, w.Error
	case <-w.done:
		return nil, nil
	}
}

func (w *watcher) Stop() error {
	if !atomic.CompareAndSwapUint32(&w.closed, 0, 1) {
		return nil
	}
	close(w.done)
	return w.client.StopWatch(w.key)
}
