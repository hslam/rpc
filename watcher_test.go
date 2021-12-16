// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package rpc

import (
	"testing"
)

func TestWatcherTrigger(t *testing.T) {
	watcher := &watcher{}
	watcher.cond.L = &watcher.mut
	for i := byte(0); i < 255; i++ {
		e := getEvent()
		e.Value = []byte{i}
		watcher.trigger(e)
	}
	for i := byte(0); i < 255; i++ {
		v, err := watcher.Wait()
		if err != nil {
			t.Error(err)
		} else if len(v) == 0 {
			t.Error("len == 0")
		} else if v[0] != i {
			t.Error("out of order")
		}
	}
	go func() {
		watcher.stop()
	}()
	_, err := watcher.Wait()
	if err != ErrWatcherShutdown {
		t.Error(err)
	}
}
