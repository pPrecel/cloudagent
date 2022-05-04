package schema

import (
	"errors"
	"io/ioutil"
	"testing"

	command "github.com/pPrecel/cloud-agent/cmd"
	"github.com/pPrecel/cloud-agent/pkg/config"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestNewCmd(t *testing.T) {
	o := NewOptions(&command.Options{})
	c := NewCmd(o)

	assert.Equal(t, 0, len(c.Commands()))
}

func Test_Cmd(t *testing.T) {
	l := logrus.New()
	l.Out = ioutil.Discard

	t.Run("validate and generate plist", func(t *testing.T) {
		o := &options{
			stdout:     ioutil.Discard,
			jsonSchema: config.JSONSchema,
			Options:    &command.Options{},
		}
		c := NewCmd(o)

		err := c.RunE(c, []string{})
		assert.NoError(t, err)
	})

	t.Run("reflect error", func(t *testing.T) {
		err := run(&options{
			stdout:     ioutil.Discard,
			jsonSchema: func() ([]byte, error) { return []byte{}, errors.New("test error") },
			Options:    &command.Options{},
		})

		assert.Error(t, err)
	})

	t.Run("reflect error", func(t *testing.T) {
		err := run(&options{
			stdout:     ioutil.Discard,
			jsonSchema: func() ([]byte, error) { return []byte("test test"), nil },
			Options:    &command.Options{},
		})

		assert.Error(t, err)
	})
}
