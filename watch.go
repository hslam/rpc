package rpc

// Watch represents a watcher.
type Watch struct {
	client *Client
	C      chan *Watch
	key    string
	Value  []byte
	Error  error
}

func (w *Watch) trigger(value []byte, err error) {
	w.Value = value
	w.Error = err
	select {
	case w.C <- w:
	default:
	}
}

// Wait will return value when the key is triggered.
func (w *Watch) Wait() ([]byte, error) {
	<-w.C
	return w.Value, w.Error
}

// Stop stops the watch.
func (w *Watch) Stop() error {
	return w.client.StopWatch(w.key)
}
