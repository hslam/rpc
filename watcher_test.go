package rpc

import (
	"testing"
)

func TestWatcherTrigger(t *testing.T) {
	watcher := &watcher{C: make(chan *watcher, 1)}
	watcher.trigger([]byte{}, nil)
	watcher.trigger([]byte{}, nil)
}
