package check

import (
	"context"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gardener/gardener/pkg/apis/core/v1beta1"
	command "github.com/pPrecel/cloudagent/cmd"
	"github.com/pPrecel/cloudagent/internal/output"
	"github.com/pPrecel/cloudagent/pkg/agent"
	cloud_agent "github.com/pPrecel/cloudagent/pkg/agent/proto"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	googlerpc "google.golang.org/grpc"
)

func TestNewCmd(t *testing.T) {
	o := NewOptions(&command.Options{})
	c := NewCmd(o)

	t.Run("defaults", func(t *testing.T) {
		assert.Equal(t, *output.NewFlag(&output.Flag{}, "table", "$h/$e/$a", "$E"), o.outFormat)
		assert.Equal(t, 2*time.Second, o.timeout)
		assert.Equal(t, agent.Address, o.socketAddress)
	})

	t.Run("parse flags", func(t *testing.T) {
		c.ParseFlags([]string{
			"--output", "text=$a=$e",
			"--timeout", "5s",
			"--socket-path", "/tmp/tmpsocket.sock",
		})

		assert.Equal(t, *output.NewFlag(&output.Flag{}, "text", "$a", "$e"), o.outFormat)
		assert.Equal(t, 5*time.Second, o.timeout)
		assert.Equal(t, "/tmp/tmpsocket.sock", o.socketAddress)
	})

	t.Run("parse shortcuts", func(t *testing.T) {
		c.ParseFlags([]string{
			"-o", "text=$a$a$a=$e$e$e",
			"-t", "10s",
		})

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
		c := &agent.ServerCache{
			GardenerCache: agent.NewCache[*v1beta1.ShootList](),
		}
		r := c.GardenerCache.Register("test-data")
		stopFn, err := fixServer(l, c)
		assert.NoError(t, err)
		defer stopFn()

		o := &options{
			stdout: io.Discard,
			Options: &command.Options{
				Logger:  l.Logger,
				Context: context.Background(),
			},
		}
		cmd := NewCmd(o)
		o.socketAddress = socketAddress
		o.socketNetwork = socketNetwork
		o.outFormat = *output.NewFlag(&o.outFormat, output.TextType, "$a", "-")

		r.Set(&v1beta1.ShootList{
			Items: []v1beta1.Shoot{
				{}, {}, {},
			},
		}, nil)

		err = cmd.RunE(cmd, []string{})
		assert.NoError(t, err)
	})
	t.Run("client error", func(t *testing.T) {
		c := agent.NewCache[*v1beta1.ShootList]()
		r := c.Register("test-data")

		o := &options{
			stdout: io.Discard,
			Options: &command.Options{
				Logger:  l.Logger,
				Context: context.Background(),
			},
		}
		cmd := NewCmd(o)
		o.socketAddress = socketAddress
		o.socketNetwork = socketNetwork
		o.outFormat = *output.NewFlag(&o.outFormat, output.TextType, "$a", "$e")

		r.Set(&v1beta1.ShootList{
			Items: []v1beta1.Shoot{
				{}, {}, {},
			},
		}, nil)

		err := cmd.RunE(cmd, []string{})
		assert.Error(t, err)
	})

	t.Run("dial error", func(t *testing.T) {
		c := agent.NewCache[*v1beta1.ShootList]()
		r := c.Register("test-data")

		o := &options{
			stdout: io.Discard,
			Options: &command.Options{
				Logger:  l.Logger,
				Context: context.Background(),
			},
		}
		cmd := NewCmd(o)
		o.socketAddress = "\n"
		o.outFormat = *output.NewFlag(&o.outFormat, output.TextType, "$a", "$e")

		r.Set(&v1beta1.ShootList{
			Items: []v1beta1.Shoot{
				{}, {}, {},
			},
		}, nil)

		err := cmd.RunE(cmd, []string{})
		assert.Error(t, err)
	})
	t.Run("general error", func(t *testing.T) {
		c := &agent.ServerCache{
			GardenerCache: agent.NewCache[*v1beta1.ShootList](),
			GeneralError:  errors.New("test error"),
		}
		stopFn, err := fixServer(l, c)
		assert.NoError(t, err)
		defer stopFn()

		o := &options{
			stdout: io.Discard,
			Options: &command.Options{
				Logger:  l.Logger,
				Context: context.Background(),
			},
		}
		cmd := NewCmd(o)
		o.socketAddress = socketAddress
		o.socketNetwork = socketNetwork
		o.outFormat = *output.NewFlag(&o.outFormat, output.TableType, "$r/$h/$u/$a", "-/-/-/-")

		err = cmd.RunE(cmd, []string{})
		assert.NoError(t, err)
	})
}

func fixServer(l *logrus.Entry, c *agent.ServerCache) (stop func(), err error) {
	lis, err := agent.NewSocket(socketNetwork, socketAddress)
	if err != nil {
		return nil, err
	}

	grpcServer := googlerpc.NewServer(googlerpc.EmptyServerOption{})
	agentServer := agent.NewServer(&agent.ServerOption{
		Cache:  c,
		Logger: l,
	})
	cloud_agent.RegisterAgentServer(grpcServer, agentServer)

	go grpcServer.Serve(lis)

	return grpcServer.Stop, nil
}
