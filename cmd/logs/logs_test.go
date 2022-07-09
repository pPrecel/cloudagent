package logs

import (
	"errors"
	"testing"

	"github.com/hpcloud/tail"
	command "github.com/pPrecel/cloudagent/cmd"
	"github.com/pPrecel/cloudagent/pkg/brew"
	"github.com/stretchr/testify/assert"
)

func TestNewCmd(t *testing.T) {
	o := NewOptions(&command.Options{})
	c := NewCmd(o)

	t.Run("defaults", func(t *testing.T) {
		assert.Equal(t, brew.StdoutPath, o.filePath)
		assert.Equal(t, false, o.followLogs)
	})

	t.Run("parse shortcuts", func(t *testing.T) {
		c.ParseFlags([]string{
			"-f",
		})

		assert.Equal(t, true, o.followLogs)
	})

	t.Run("parse flags", func(t *testing.T) {
		c.ParseFlags([]string{
			"--file", "/test/path",
			"--follow", "true",
		})

		assert.Equal(t, "/test/path", o.filePath)
		assert.Equal(t, true, o.followLogs)
	})
}

type mockWriter struct {
	val    []byte
	rInt   int
	rError error
}

func (w *mockWriter) Write(b []byte) (int, error) {
	w.val = append(w.val, b...)
	return w.rInt, w.rError
}

func Test_run(t *testing.T) {
	t.Run("run and print", func(t *testing.T) {
		w := &mockWriter{}
		o := &options{
			writer: w,
			tailFile: func(_ string, _ tail.Config) (*tail.Tail, error) {
				lines := make(chan *tail.Line)

				t := &tail.Tail{
					Lines: lines,
				}

				go func() {
					defer t.Done()
					defer close(lines)
					lines <- &tail.Line{Text: "line1"}
					lines <- &tail.Line{Text: "line2"}
				}()

				return t, nil
			},
		}
		cmd := NewCmd(o)

		assert.NoError(t, cmd.RunE(cmd, []string{}))
		assert.Equal(t, []byte("line1\nline2\n"), w.val)
	})

	t.Run("tailFile error", func(t *testing.T) {
		o := &options{
			tailFile: func(_ string, _ tail.Config) (*tail.Tail, error) {
				return nil, errors.New("test error")
			},
		}
		cmd := NewCmd(o)

		assert.Error(t, cmd.RunE(cmd, []string{}))
	})

	t.Run("run and print text", func(t *testing.T) {
		o := &options{
			writer: &mockWriter{
				rError: errors.New("test error"),
			},
			tailFile: func(_ string, _ tail.Config) (*tail.Tail, error) {
				lines := make(chan *tail.Line)

				t := &tail.Tail{
					Lines: lines,
				}

				go func() {
					defer t.Done()
					defer close(lines)
					lines <- &tail.Line{Text: "line1"}
				}()

				return t, nil
			},
		}
		cmd := NewCmd(o)

		assert.Error(t, cmd.RunE(cmd, []string{}))
	})
}
