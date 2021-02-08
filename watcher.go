// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package rpc

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
)

// ErrWatcherShutdown is returned when the watcher is shut down.
var ErrWatcherShutdown = errors.New("The watcher is shut down")

var eventPool = sync.Pool{New: func() interface{} {
	return &event{}
}}

type event struct {
	Value []byte
	Error error
}

func getEvent() *event {
	return eventPool.Get().(*event)
}

func freeEvent(e *event) {
	*e = event{}
	eventPool.Put(e)
}

// Watcher represents a watcher.
type Watcher interface {
	// Wait will return value when the key is triggered.
	Wait() ([]byte, error)
	// WaitWithContext acts like Wait but takes a context.
	WaitWithContext(context.Context) ([]byte, error)
	// Stop stops the watch.
	Stop() error
}

type watcher struct {
	conn   *Conn
	C      chan *event
	key    string
	mut    sync.Mutex
	events []*event
	done   chan struct{}
	closed uint32
}

func (w *watcher) trigger(e *event) {
	w.mut.Lock()
	if len(w.events) == 0 {
		select {
		case w.C <- e:
			w.mut.Unlock()
			return
		default:
		}
	}
	w.events = append(w.events, e)
	w.mut.Unlock()
}

func (w *watcher) triggerNext() {
	w.mut.Lock()
	for len(w.events) > 0 && cap(w.C) > len(w.C) {
		next := w.events[0]
		w.C <- next
		n := copy(w.events, w.events[1:])
		w.events = w.events[:n]
	}
	w.mut.Unlock()
}

func (w *watcher) Wait() (value []byte, err error) {
	select {
	case e := <-w.C:
		w.triggerNext()
		value = e.Value
		err = e.Error
		freeEvent(e)
	case <-w.done:
		err = ErrWatcherShutdown
	}
	return
}

func (w *watcher) WaitWithContext(ctx context.Context) (value []byte, err error) {
	select {
	case e := <-w.C:
		w.triggerNext()
		value = e.Value
		err = e.Error
		freeEvent(e)
	case <-ctx.Done():
		err = ctx.Err()
	case <-w.done:
		err = ErrWatcherShutdown
	}
	return
}

func (w *watcher) stop() {
	if atomic.CompareAndSwapUint32(&w.closed, 0, 1) {
		close(w.done)
	}
}

func (w *watcher) Stop() error {
	w.stop()
	return w.conn.stopWatch(w.key)
}
