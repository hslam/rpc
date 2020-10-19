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
	client *Client
	C      chan *event
	key    string
	mut    sync.Mutex
	events []*event
	done   chan struct{}
	closed uint32
}

func (w *watcher) trigger(e *event) {
	w.mut.Lock()
	defer w.mut.Unlock()
	if len(w.events) == 0 {
		select {
		case w.C <- e:
			return
		default:
		}
	}
	w.events = append(w.events, e)
}

func (w *watcher) triggerNext() {
	w.mut.Lock()
	defer w.mut.Unlock()
	for len(w.events) > 0 && cap(w.C) > len(w.C) {
		next := w.events[0]
		w.C <- next
		n := copy(w.events, w.events[1:])
		w.events = w.events[:n]
	}
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

func (w *watcher) Stop() error {
	if !atomic.CompareAndSwapUint32(&w.closed, 0, 1) {
		return nil
	}
	close(w.done)
	return w.client.StopWatch(w.key)
}
