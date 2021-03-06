// Copyright (c) 2019 Meng Huang (mhboy@outlook.com)
// This package is licensed under a MIT license that can be found in the LICENSE file.

package rpc

import (
	"context"
	"testing"
	"time"
)

func TestWatcherTrigger(t *testing.T) {
	watcher := &watcher{C: make(chan *event, 10), done: make(chan struct{}, 1)}
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
		close(watcher.done)
	}()
	_, err := watcher.Wait()
	if err != ErrWatcherShutdown {
		t.Error(err)
	}
}

func TestWatcherTriggerTimeout(t *testing.T) {
	watcher := &watcher{C: make(chan *event, 10), done: make(chan struct{}, 1)}
	for i := byte(0); i < 255; i++ {
		e := getEvent()
		e.Value = []byte{i}
		watcher.trigger(e)
	}
	for i := byte(0); i < 255; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		v, err := watcher.WaitWithContext(ctx)
		if err != nil {
			t.Error(err)
		} else if len(v) == 0 {
			t.Error("len == 0")
		} else if v[0] != i {
			t.Error("out of order")
		}
		cancel()
	}
	go func() {
		watcher.stop()
	}()
	_, err := watcher.Wait()
	if err != ErrWatcherShutdown {
		t.Error(err)
	}
}

func TestWatcherTriggerTimeoutErr(t *testing.T) {
	watcher := &watcher{C: make(chan *event, 10), done: make(chan struct{}, 1)}
	{
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*1)
		_, err := watcher.WaitWithContext(ctx)
		if err == nil {
			t.Error()
		}
		cancel()
	}
	watcher.stop()
	{
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*1)
		_, err := watcher.WaitWithContext(ctx)
		if err != ErrWatcherShutdown {
			t.Error(err)
		}
		cancel()
	}

}
