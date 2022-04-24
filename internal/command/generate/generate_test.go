package generate

import (
	"errors"
	"io/ioutil"
	"os"
	"testing"

	"github.com/pPrecel/cloud-agent/internal/command"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestNewCmd(t *testing.T) {
	o := NewOptions(&command.Options{})
	c := NewCmd(o)

	t.Run("defaults", func(t *testing.T) {
		assert.Equal(t, "", o.kubeconfigPath)
		assert.Equal(t, "", o.namespace)
		assert.Equal(t, "@every 60s", o.cronSpec)
		assert.Equal(t, false, o.agentVerbose)
	})

	t.Run("parse flags", func(t *testing.T) {
		c.ParseFlags([]string{
			"--kubeconfigPath", "other-path",
			"--namespace", "other-namespace",
			"--cronSpec", "@every 15m",
			"--agentVerbose", "true",
		})

		assert.Equal(t, "other-path", o.kubeconfigPath)
		assert.Equal(t, "other-namespace", o.namespace)
		assert.Equal(t, "@every 15m", o.cronSpec)
		assert.Equal(t, true, o.agentVerbose)
	})

	t.Run("parse shortcuts", func(t *testing.T) {
		c.ParseFlags([]string{
			"-k", "path",
			"-n", "namespace",
			"-c", "@every 20m",
		})

		assert.Equal(t, "path", o.kubeconfigPath)
		assert.Equal(t, "namespace", o.namespace)
		assert.Equal(t, "@every 20m", o.cronSpec)
		assert.Equal(t, true, o.agentVerbose)
	})
}

func Test_Cmd(t *testing.T) {
	l := logrus.New()
	l.Out = ioutil.Discard

	t.Run("validate and generate plist", func(t *testing.T) {
		o := &options{
			executable: os.Executable,
			stdout:     ioutil.Discard,
			Options: &command.Options{
				Logger: l,
			},
		}
		c := NewCmd(o)
		o.kubeconfigPath = "anything"
		o.namespace = "something"

		err := c.PreRunE(c, []string{})
		assert.NoError(t, err)

		err = c.RunE(c, []string{})
		assert.NoError(t, err)
	})

	t.Run("executable error", func(t *testing.T) {
		err := run(&options{
			kubeconfigPath: "anything",
			namespace:      "something",
			cronSpec:       "2s",
			agentVerbose:   true,
			executable: func() (string, error) {
				return "", errors.New("test error")
			},
			Options: &command.Options{
				Logger: l,
			},
		})

		assert.Error(t, err)
	})
}
