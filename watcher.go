// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package rpc

import (
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
	// Stop stops the watch.
	Stop() error
}

type watcher struct {
	conn   *Conn
	mut    sync.Mutex
	cond   sync.Cond
	key    string
	events []*event
	done   bool
	closed uint32
}

func (w *watcher) trigger(e *event) {
	w.mut.Lock()
	w.events = append(w.events, e)
	w.mut.Unlock()
	w.cond.Signal()
}

func (w *watcher) Wait() (value []byte, err error) {
	if atomic.LoadUint32(&w.closed) > 0 {
		err = ErrWatcherShutdown
		return
	}
	w.mut.Lock()
	for {
		if len(w.events) > 0 {
			e := w.events[0]
			w.events = w.events[1:]
			w.mut.Unlock()
			value = e.Value
			err = e.Error
			freeEvent(e)
			return
		}
		w.cond.Wait()
		if atomic.LoadUint32(&w.closed) > 0 {
			w.mut.Unlock()
			err = ErrWatcherShutdown
			return
		}
	}
}

func (w *watcher) stop() {
	atomic.StoreUint32(&w.closed, 1)
	w.cond.Broadcast()
}

func (w *watcher) Stop() error {
	w.stop()
	return w.conn.stopWatch(w.key)
}
