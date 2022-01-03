// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package rpc

import (
	"errors"
	"github.com/hslam/scheduler"
	"sync"
	"sync/atomic"
)

// ErrStreamShutdown is returned when the stream is shut down.
var ErrStreamShutdown = errors.New("The stream is shut down")

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

// SetStream is used to connect rpc Stream.
type SetStream interface {
	Connect(stream Stream) error
}

// Stream defines the message stream interface.
type Stream interface {
	// WriteMessage writes a message to the stream.
	WriteMessage(m interface{}) error
	// ReadMessage reads a single message from the stream.
	ReadMessage(b []byte, m interface{}) error
	// Close closes the stream.
	Close() error
}

type stream struct {
	close       func() error
	unmarshal   func(data []byte, v interface{}) error
	write       func(m interface{}) (err error)
	writeStream scheduler.Scheduler
	mut         sync.Mutex
	cond        sync.Cond
	seq         uint64
	events      []*event
	done        bool
	closed      int32
}

func (w *stream) trigger(e *event) {
	w.mut.Lock()
	w.events = append(w.events, e)
	w.mut.Unlock()
	w.cond.Signal()
}

func (w *stream) ReadMessage(b []byte, m interface{}) (err error) {
	w.mut.Lock()
	if atomic.LoadInt32(&w.closed) > 0 {
		w.mut.Unlock()
		err = ErrStreamShutdown
		return
	}
	for {
		if len(w.events) > 0 {
			e := w.events[0]
			w.events = w.events[1:]
			w.mut.Unlock()
			if cap(b) > len(e.Value) {
				b = b[:len(e.Value)]
			} else {
				b = make([]byte, len(e.Value))
			}
			copy(b, e.Value)
			PutBuffer(e.Value)
			w.unmarshal(b, m)
			err = e.Error
			freeEvent(e)
			return
		}
		w.cond.Wait()
		if atomic.LoadInt32(&w.closed) > 0 {
			w.mut.Unlock()
			err = ErrStreamShutdown
			return
		}
	}
}

func (w *stream) WriteMessage(m interface{}) (err error) {
	if atomic.LoadInt32(&w.closed) > 0 {
		err = ErrStreamShutdown
		return
	}
	w.writeStream.Schedule(func() {
		w.write(m)
	})
	return
}

func (w *stream) stop() {
	w.writeStream.Close()
	w.mut.Lock()
	atomic.StoreInt32(&w.closed, 1)
	w.mut.Unlock()
	w.cond.Broadcast()
}

func (w *stream) Close() error {
	w.stop()
	return w.close()
}
