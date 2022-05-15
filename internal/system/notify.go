package system

import (
	"github.com/fsnotify/fsnotify"
)

type Notifier struct {
	s chan int

	Errors     chan error
	IsMotified chan interface{}
	Stop       func()
}

func NotifyChange(path string) (*Notifier, error) {
	return notifyChange(path, fsnotify.NewWatcher)
}

func notifyChange(path string, fsNotifyFn func() (*fsnotify.Watcher, error)) (*Notifier, error) {
	n := &Notifier{
		Errors:     make(chan error),
		IsMotified: make(chan interface{}),
		s:          make(chan int),
	}

	n.Stop = func() {
		n.s <- 1
	}

	watcher, err := fsNotifyFn()
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			handleEvent(n, watcher)
		}
	}()

	return n, watcher.Add(path)
}

func handleEvent(n *Notifier, w *fsnotify.Watcher) {
	select {
	case <-n.s:
		defer w.Close()
		return
	case <-w.Events:
		n.IsMotified <- 1
	case err := <-w.Errors:
		n.Errors <- err
	}
}
