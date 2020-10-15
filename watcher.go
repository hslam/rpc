package rpc

// Watcher represents a watcher.
type Watcher interface {
	Wait() ([]byte, error)
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

// Wait will return value when the key is triggered.
func (w *watcher) Wait() ([]byte, error) {
	<-w.C
	return w.Value, w.Error
}

// Stop stops the watch.
func (w *watcher) Stop() error {
	return w.client.StopWatch(w.key)
}
