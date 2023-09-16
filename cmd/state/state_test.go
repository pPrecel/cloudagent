package state

import (
	"context"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	command "github.com/pPrecel/cloudagent/cmd"
	"github.com/pPrecel/cloudagent/internal/output"
	"github.com/pPrecel/cloudagent/pkg/agent"
	cloud_agent "github.com/pPrecel/cloudagent/pkg/agent/proto"
	"github.com/pPrecel/cloudagent/pkg/cache"
	"github.com/pPrecel/cloudagent/pkg/types"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	googlerpc "google.golang.org/grpc"
)

func TestNewCmd(t *testing.T) {
	o := NewOptions(&command.Options{})
	c := NewCmd(o)

	t.Run("defaults", func(t *testing.T) {
		assert.Equal(t, "", o.createdBy)
		assert.Equal(t, *output.NewFlag(&output.Flag{}, "table", "$r/$h/$x/$a", "-/-/-/-"), o.outFormat)
		assert.Equal(t, 2*time.Second, o.timeout)
		assert.Equal(t, agent.Address, o.socketAddress)
	})

	t.Run("parse flags", func(t *testing.T) {
		c.ParseFlags([]string{
			"--created-by", "owner",
			"--output", "text=$a=$e",
			"--timeout", "5s",
			"--socket-path", "/tmp/tmpsocket.sock",
		})

		assert.Equal(t, "owner", o.createdBy)
		assert.Equal(t, *output.NewFlag(&output.Flag{}, "text", "$a", "$e"), o.outFormat)
		assert.Equal(t, 5*time.Second, o.timeout)
		assert.Equal(t, "/tmp/tmpsocket.sock", o.socketAddress)
	})

	t.Run("parse shortcuts", func(t *testing.T) {
		c.ParseFlags([]string{
			"-c", "other-owner",
			"-o", "text=$a$a$a=$e$e$e",
			"-t", "10s",
		})

		assert.Equal(t, "other-owner", o.createdBy)
		assert.Equal(t, *output.NewFlag(&output.Flag{}, "text", "$a$a$a", "$e$e$e"), o.outFormat)
		assert.Equal(t, 10*time.Second, o.timeout)
	})
}

var (
	socketAddress = filepath.Join(os.TempDir(), "state-test-socket.sock")
	socketNetwork = "unix"
)

func Test_run(t *testing.T) {
	l := &logrus.Entry{
		Logger: logrus.New(),
	}
	l.Logger.Out = ioutil.Discard

	t.Run("run and print text", func(t *testing.T) {
		c := &cache.ServerCache{
			GardenerCache: cache.NewGardenerCache(),
		}
		r := c.GardenerCache.Register("test-data")
		stopFn, err := fixServer(l, c)
		assert.NoError(t, err)
		defer stopFn()

		o := &options{
			writer: io.Discard,
			Options: &command.Options{
				Logger:  l.Logger,
				Context: context.Background(),
			},
		}
		cmd := NewCmd(o)
		o.createdBy = "owner"
		o.socketAddress = socketAddress
		o.socketNetwork = socketNetwork
		o.outFormat = *output.NewFlag(&o.outFormat, output.TextType, "$r/$h/$u/$a", "-/-/-/-")
		o.project = "test"
		o.condition = "test"
		o.createdAfter = "2022-09-11"
		o.createdBefore = "2022-09-11"
		o.updatedAfter = "2022-09-11"
		o.updatedBefore = "2022-09-11"

		r.Set(&types.ShootList{
			Items: []types.Shoot{
				{}, {}, {},
			},
		}, nil)

		err = cmd.RunE(cmd, []string{})
		assert.NoError(t, err)
	})

	t.Run("parse createdAfter error", func(t *testing.T) {
		c := &cache.ServerCache{
			GardenerCache: cache.NewGardenerCache(),
		}
		stopFn, err := fixServer(l, c)
		assert.NoError(t, err)
		defer stopFn()

		o := &options{
			writer: io.Discard,
			Options: &command.Options{
				Logger:  l.Logger,
				Context: context.Background(),
			},
		}
		cmd := NewCmd(o)
		o.createdBy = "owner"
		o.socketAddress = socketAddress
		o.socketNetwork = socketNetwork
		o.outFormat = *output.NewFlag(&o.outFormat, output.TextType, "$r/$h/$u/$a", "-/-/-/-")
		o.project = "test"
		o.condition = "test"
		o.createdAfter = "wrong time"

		err = cmd.RunE(cmd, []string{})
		assert.Error(t, err)
	})

	t.Run("parse createdBefore error", func(t *testing.T) {
		c := &cache.ServerCache{
			GardenerCache: cache.NewGardenerCache(),
		}
		stopFn, err := fixServer(l, c)
		assert.NoError(t, err)
		defer stopFn()

		o := &options{
			writer: io.Discard,
			Options: &command.Options{
				Logger:  l.Logger,
				Context: context.Background(),
			},
		}
		cmd := NewCmd(o)
		o.createdBy = "owner"
		o.socketAddress = socketAddress
		o.socketNetwork = socketNetwork
		o.outFormat = *output.NewFlag(&o.outFormat, output.TextType, "$r/$h/$u/$a", "-/-/-/-")
		o.project = "test"
		o.condition = "test"
		o.createdBefore = "wrong time"

		err = cmd.RunE(cmd, []string{})
		assert.Error(t, err)
	})

	t.Run("parse updatedAfter error", func(t *testing.T) {
		c := &cache.ServerCache{
			GardenerCache: cache.NewGardenerCache(),
		}
		stopFn, err := fixServer(l, c)
		assert.NoError(t, err)
		defer stopFn()

		o := &options{
			writer: io.Discard,
			Options: &command.Options{
				Logger:  l.Logger,
				Context: context.Background(),
			},
		}
		cmd := NewCmd(o)
		o.createdBy = "owner"
		o.socketAddress = socketAddress
		o.socketNetwork = socketNetwork
		o.outFormat = *output.NewFlag(&o.outFormat, output.TextType, "$r/$h/$u/$a", "-/-/-/-")
		o.project = "test"
		o.condition = "test"
		o.updatedAfter = "wrong time"

		err = cmd.RunE(cmd, []string{})
		assert.Error(t, err)
	})

	t.Run("parse updatedBefore error", func(t *testing.T) {
		c := &cache.ServerCache{
			GardenerCache: cache.NewGardenerCache(),
		}
		stopFn, err := fixServer(l, c)
		assert.NoError(t, err)
		defer stopFn()

		o := &options{
			writer: io.Discard,
			Options: &command.Options{
				Logger:  l.Logger,
				Context: context.Background(),
			},
		}
		cmd := NewCmd(o)
		o.createdBy = "owner"
		o.socketAddress = socketAddress
		o.socketNetwork = socketNetwork
		o.outFormat = *output.NewFlag(&o.outFormat, output.TextType, "$r/$h/$u/$a", "-/-/-/-")
		o.project = "test"
		o.condition = "test"
		o.updatedBefore = "wrong time"

		err = cmd.RunE(cmd, []string{})
		assert.Error(t, err)
	})

	t.Run("client error", func(t *testing.T) {
		c := cache.NewGardenerCache()
		r := c.Register("test-data")

		o := &options{
			writer: io.Discard,
			Options: &command.Options{
				Logger:  l.Logger,
				Context: context.Background(),
			},
		}
		cmd := NewCmd(o)
		o.createdBy = "owner"
		o.socketAddress = socketAddress
		o.socketNetwork = socketNetwork
		o.outFormat = *output.NewFlag(&o.outFormat, output.TextType, "$a", "$e")

		r.Set(&types.ShootList{
			Items: []types.Shoot{
				{}, {}, {},
			},
		}, nil)

		err := cmd.RunE(cmd, []string{})
		assert.Error(t, err)
	})

	t.Run("dial error", func(t *testing.T) {
		c := cache.NewGardenerCache()
		r := c.Register("test-data")

		o := &options{
			writer: io.Discard,
			Options: &command.Options{
				Logger:  l.Logger,
				Context: context.Background(),
			},
		}
		cmd := NewCmd(o)
		o.createdBy = "owner"
		o.socketAddress = "\n"
		o.outFormat = *output.NewFlag(&o.outFormat, output.TextType, "$a", "$e")

		r.Set(&types.ShootList{
			Items: []types.Shoot{
				{}, {}, {},
			},
		}, nil)

		err := cmd.RunE(cmd, []string{})
		assert.Error(t, err)
	})
	t.Run("request error", func(t *testing.T) {
		c := &cache.ServerCache{
			GardenerCache: cache.NewGardenerCache(),
		}
		r := c.GardenerCache.Register("test-data")
		stopFn, err := fixServer(l, c)
		assert.NoError(t, err)
		defer stopFn()

		o := &options{
			writer: io.Discard,
			Options: &command.Options{
				Logger:  l.Logger,
				Context: context.Background(),
			},
		}
		cmd := NewCmd(o)
		o.createdBy = "owner"
		o.socketAddress = socketAddress
		o.socketNetwork = socketNetwork
		o.outFormat = *output.NewFlag(&o.outFormat, output.TableType, "$r/$h/$u/$a", "-/-/-/-")

		r.Set(&types.ShootList{
			Items: []types.Shoot{
				{}, {}, {},
			},
		}, errors.New("test error"))

		err = cmd.RunE(cmd, []string{})
		assert.NoError(t, err)
	})
	t.Run("request general error", func(t *testing.T) {
		c := &cache.ServerCache{
			GeneralError:  errors.New("test error"),
			GardenerCache: cache.NewGardenerCache(),
		}
		stopFn, err := fixServer(l, c)
		assert.NoError(t, err)
		defer stopFn()

		o := &options{
			writer: io.Discard,
			Options: &command.Options{
				Logger:  l.Logger,
				Context: context.Background(),
			},
		}
		cmd := NewCmd(o)
		o.createdBy = "owner"
		o.socketAddress = socketAddress
		o.socketNetwork = socketNetwork
		o.outFormat = *output.NewFlag(&o.outFormat, output.TextType, "$r/$h/$u/$a", "-/-/-/-")

		err = cmd.RunE(cmd, []string{})
		assert.NoError(t, err)
	})
}

func fixServer(l *logrus.Entry, c *cache.ServerCache) (stop func(), err error) {
	lis, err := agent.NewSocket(socketNetwork, socketAddress)
	if err != nil {
		return nil, err
	}

	grpcServer := googlerpc.NewServer(googlerpc.EmptyServerOption{})
	agentServer := agent.NewServer(&agent.ServerOption{
		ResourceGetter: c,
		Logger:         l,
	})
	cloud_agent.RegisterAgentServer(grpcServer, agentServer)

	go grpcServer.Serve(lis)

	return grpcServer.Stop, nil
}
