package system

import (
	"errors"
	"io/fs"
	"io/ioutil"
	"os"
	"testing"

	"github.com/fsnotify/fsnotify"
	"github.com/stretchr/testify/assert"
)

func TestNotifyChange(t *testing.T) {
	t.Run("notify change and stop", func(t *testing.T) {
		f, err := os.CreateTemp(os.TempDir(), "cloudagent-notify-test-")
		assert.NoError(t, err)
		defer f.Close()
		defer os.Remove(f.Name())

		n, err := NotifyChange(f.Name())
		assert.NoError(t, err)

		err = ioutil.WriteFile(f.Name(), []byte("any string"), fs.ModePerm)
		assert.NoError(t, err)

		go func() {
			for {
				assert.NotNil(t, <-n.IsMotified)
			}
		}()

		n.Stop()
	})

	t.Run("notifyFn error", func(t *testing.T) {
		w, err := notifyChange("", func() (*fsnotify.Watcher, error) {
			return nil, errors.New("test error")
		})

		assert.Nil(t, w)
		assert.Error(t, err)
	})
}

func Test_handleEvent(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		w := &fsnotify.Watcher{
			Errors: make(chan error),
			Events: make(chan fsnotify.Event),
		}

		n := &Notifier{
			Errors:     make(chan error),
			IsMotified: make(chan interface{}),
			s:          make(chan int),
		}

		go handleEvent(n, w)

		w.Errors <- errors.New("test error")
	})
	t.Run("modify", func(t *testing.T) {
		w := &fsnotify.Watcher{
			Errors: make(chan error),
			Events: make(chan fsnotify.Event),
		}

		n := &Notifier{
			Errors:     make(chan error),
			IsMotified: make(chan interface{}),
			s:          make(chan int),
		}

		go handleEvent(n, w)

		w.Events <- fsnotify.Event{}
	})
}
