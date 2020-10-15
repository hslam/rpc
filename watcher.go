package rpc

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
	<-w.C
	return w.Value, w.Error
}

func (w *watcher) Stop() error {
	return w.client.StopWatch(w.key)
}
