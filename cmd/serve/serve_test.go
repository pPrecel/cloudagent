package serve

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	v1beta1_apis "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	command "github.com/pPrecel/cloudagent/cmd"
	"github.com/pPrecel/cloudagent/pkg/agent"
	"github.com/pPrecel/cloudagent/pkg/config"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestNewCmd(t *testing.T) {
	o := NewOptions(&command.Options{})
	c := NewCmd(o)

	t.Run("defaults", func(t *testing.T) {
		assert.Equal(t, config.ConfigPath, o.configPath)
	})

	t.Run("parse flags", func(t *testing.T) {
		c.ParseFlags([]string{
			"--config-path", "path",
		})

		assert.Equal(t, "path", o.configPath)
	})

	t.Run("parse shortcuts", func(t *testing.T) {
		c.ParseFlags([]string{
			"-c", "other-path",
		})

		assert.Equal(t, "other-path", o.configPath)
	})
}

var (
	testNetwork = "unix"
	testAddress = filepath.Join(os.TempDir(), "serve-test-socket.sock")
)

func Test_run(t *testing.T) {
	t.Run("validate and run", func(t *testing.T) {
		l := logrus.New()
		l.Out = io.Discard
		o := &options{
			Options: &command.Options{
				Logger:  l,
				Context: context.Background(),
			},
			configPath:    "/empty/path",
			socketNetwork: testNetwork,
			socketAddress: testAddress,
		}
		c := NewCmd(o)

		err := c.PreRunE(c, []string{})
		assert.NoError(t, err)

		go func() {
			assert.NoError(t, c.RunE(c, []string{}))
		}()

		socketExist := false
		for i := 0; i < 5; i++ {
			time.Sleep(1 * time.Second)

			_, err = os.Stat(testAddress)
			if err == nil {
				socketExist = true
				break
			}
		}

		assert.True(t, socketExist, "socket does not exist")
	})

	t.Run("socket error", func(t *testing.T) {
		l := logrus.New()
		l.Out = io.Discard
		o := &options{
			Options: &command.Options{
				Logger:  l,
				Context: context.Background(),
			},
			socketNetwork: testNetwork,
			socketAddress: "/addr\n\n\n",
		}

		c := NewCmd(o)

		assert.Error(t, c.RunE(c, []string{}))
	})
}

func Test_startWatcher(t *testing.T) {
	l := logrus.New()
	l.Out = io.Discard

	t.Run("handle error", func(t *testing.T) {
		c := agent.NewCache[*v1beta1_apis.ShootList]()
		startWatcher(&options{
			Options: &command.Options{
				Logger:  l,
				Context: context.Background(),
			},
			socketNetwork: testNetwork,
			socketAddress: testAddress,
		}, c)

		assert.Len(t, c.Resources(), 0)
	})
}
