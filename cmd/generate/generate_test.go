package generate

import (
	"errors"
	"io/ioutil"
	"os"
	"testing"

	command "github.com/pPrecel/cloud-agent/cmd"
	"github.com/pPrecel/cloud-agent/pkg/config"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestNewCmd(t *testing.T) {
	o := NewOptions(&command.Options{})
	c := NewCmd(o)

	t.Run("defaults", func(t *testing.T) {
		assert.Equal(t, config.ConfigPath, o.configPath)
		assert.Equal(t, false, o.agentVerbose)
	})

	t.Run("parse flags", func(t *testing.T) {
		c.ParseFlags([]string{
			"--configPath", "other-path",
			"--agentVerbose", "true",
		})

		assert.Equal(t, "other-path", o.configPath)
		assert.Equal(t, true, o.agentVerbose)
	})

	t.Run("parse shortcuts", func(t *testing.T) {
		c.ParseFlags([]string{
			"-c", "path",
		})

		assert.Equal(t, "path", o.configPath)
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
		o.configPath = "anything"

		err := c.PreRunE(c, []string{})
		assert.NoError(t, err)

		err = c.RunE(c, []string{})
		assert.NoError(t, err)
	})

	t.Run("executable error", func(t *testing.T) {
		err := run(&options{
			configPath:   "anything",
			agentVerbose: true,
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
