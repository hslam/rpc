package rpc

import (
	"testing"
)

func TestWatcherTrigger(t *testing.T) {
	watcher := &watcher{C: make(chan *watcher, 1)}
	watcher.trigger([]byte{}, nil)
	watcher.trigger([]byte{}, nil)
}

func TestWatcherWait(t *testing.T) {
	watcher := &watcher{C: make(chan *watcher, 1), done: make(chan struct{}, 1)}
	go func() {
		close(watcher.done)
	}()
	watcher.Wait()
}
