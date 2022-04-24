package state

import (
	"context"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/pPrecel/cloud-agent/internal/agent"
	cloud_agent "github.com/pPrecel/cloud-agent/internal/agent/proto"
	"github.com/pPrecel/cloud-agent/internal/command"
	"github.com/pPrecel/cloud-agent/internal/gardener"
	"github.com/pPrecel/cloud-agent/internal/output"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	googlerpc "google.golang.org/grpc"
)

func TestNewCmd(t *testing.T) {
	o := NewOptions(&command.Options{})
	c := NewCmd(o)

	t.Run("defaults", func(t *testing.T) {
		assert.Equal(t, "", o.createdBy)
		assert.Equal(t, *output.New(&output.Output{}, "table", "Shoots: %r/%h/%u/%a", "Error: %e"), o.outFormat)
		assert.Equal(t, 2*time.Second, o.timeout)
	})

	t.Run("parse flags", func(t *testing.T) {
		c.ParseFlags([]string{
			"--createdBy", "owner",
			"--output", "text=%a=%e",
			"--timeout", "5s",
		})

		assert.Equal(t, "owner", o.createdBy)
		assert.Equal(t, *output.New(&output.Output{}, "text", "%a", "%e"), o.outFormat)
		assert.Equal(t, 5*time.Second, o.timeout)
	})

	t.Run("parse shortcuts", func(t *testing.T) {
		c.ParseFlags([]string{
			"-c", "other-owner",
			"-o", "text=%a%a%a=%e%e%e",
			"-t", "10s",
		})

		assert.Equal(t, "other-owner", o.createdBy)
		assert.Equal(t, *output.New(&output.Output{}, "text", "%a%a%a", "%e%e%e"), o.outFormat)
		assert.Equal(t, 10*time.Second, o.timeout)
	})
}

var (
	socketAddress = filepath.Join(os.TempDir(), "state-test-socket.sock")
	socketNetwork = "unix"
)

func Test_run(t *testing.T) {
	l := logrus.New()
	l.Out = ioutil.Discard

	t.Run("run and print text", func(t *testing.T) {
		s := &gardener.LastState{}
		stopFn, err := fixServer(l, s)
		assert.NoError(t, err)
		defer stopFn()

		o := &options{
			socketAddress: socketAddress,
			socketNetwork: socketNetwork,
			writer:        io.Discard,
			Options: &command.Options{
				Logger:  l,
				Context: context.Background(),
			},
		}
		cmd := NewCmd(o)
		o.createdBy = "owner"

		s.Set(&v1beta1.ShootList{
			Items: []v1beta1.Shoot{
				{}, {}, {},
			},
		})

		err = cmd.RunE(cmd, []string{})
		assert.NoError(t, err)
	})

	t.Run("run and print json", func(t *testing.T) {
		s := &gardener.LastState{}
		stopFn, err := fixServer(l, s)
		assert.NoError(t, err)
		defer stopFn()

		o := &options{
			socketAddress: socketAddress,
			socketNetwork: socketNetwork,
			writer:        io.Discard,
			Options: &command.Options{
				Logger:  l,
				Context: context.Background(),
			},
		}
		cmd := NewCmd(o)
		o.createdBy = "owner"
		o.outFormat = *output.New(&o.outFormat, output.JsonType, "", "")

		s.Set(&v1beta1.ShootList{
			Items: []v1beta1.Shoot{
				{}, {}, {},
			},
		})

		err = cmd.RunE(cmd, []string{})
		assert.NoError(t, err)
	})

	t.Run("run and print table", func(t *testing.T) {
		s := &gardener.LastState{}
		stopFn, err := fixServer(l, s)
		assert.NoError(t, err)
		defer stopFn()

		o := &options{
			socketAddress: socketAddress,
			socketNetwork: socketNetwork,
			writer:        io.Discard,
			Options: &command.Options{
				Logger:  l,
				Context: context.Background(),
			},
		}
		cmd := NewCmd(o)
		o.createdBy = "owner"
		o.outFormat = *output.New(&o.outFormat, output.TextType, "%a", "%e")

		s.Set(&v1beta1.ShootList{
			Items: []v1beta1.Shoot{
				{}, {}, {},
			},
		})

		err = cmd.RunE(cmd, []string{})
		assert.NoError(t, err)
	})

	t.Run("empty output format", func(t *testing.T) {
		s := &gardener.LastState{}
		stopFn, err := fixServer(l, s)
		assert.NoError(t, err)
		defer stopFn()

		o := &options{
			socketAddress: socketAddress,
			socketNetwork: socketNetwork,
			writer:        io.Discard,
			Options: &command.Options{
				Logger:  l,
				Context: context.Background(),
			},
		}
		cmd := NewCmd(o)
		o.createdBy = "owner"
		o.outFormat = *output.New(&o.outFormat, "", "%a", "%e")

		s.Set(&v1beta1.ShootList{
			Items: []v1beta1.Shoot{
				{}, {}, {},
			},
		})

		err = cmd.RunE(cmd, []string{})
		assert.NoError(t, err)
	})

	t.Run("client error", func(t *testing.T) {
		s := &gardener.LastState{}

		o := &options{
			socketAddress: socketAddress,
			socketNetwork: socketNetwork,
			writer:        io.Discard,
			Options: &command.Options{
				Logger:  l,
				Context: context.Background(),
			},
		}
		cmd := NewCmd(o)
		o.createdBy = "owner"
		o.outFormat = *output.New(&o.outFormat, output.TextType, "%a", "%e")

		s.Set(&v1beta1.ShootList{
			Items: []v1beta1.Shoot{
				{}, {}, {},
			},
		})

		err := cmd.RunE(cmd, []string{})
		assert.NoError(t, err)
	})

	t.Run("dial error", func(t *testing.T) {
		s := &gardener.LastState{}

		o := &options{
			socketAddress: "\n",
			writer:        io.Discard,
			Options: &command.Options{
				Logger:  l,
				Context: context.Background(),
			},
		}
		cmd := NewCmd(o)
		o.createdBy = "owner"
		o.outFormat = *output.New(&o.outFormat, output.TextType, "%a", "%e")

		s.Set(&v1beta1.ShootList{
			Items: []v1beta1.Shoot{
				{}, {}, {},
			},
		})

		err := cmd.RunE(cmd, []string{})
		assert.NoError(t, err)
	})

	t.Run("set nil and print error", func(t *testing.T) {
		s := &gardener.LastState{}

		o := &options{
			socketAddress: socketAddress,
			socketNetwork: socketNetwork,
			writer:        io.Discard,
			Options: &command.Options{
				Logger:  l,
				Context: context.Background(),
			},
		}
		cmd := NewCmd(o)
		o.createdBy = "owner"

		s.Set(nil)

		err := cmd.RunE(cmd, []string{})
		assert.NoError(t, err)
	})
}

func fixServer(l *logrus.Logger, g agent.StateGetter) (stop func(), err error) {
	lis, err := agent.NewSocket(socketNetwork, socketAddress)
	if err != nil {
		return nil, err
	}

	grpcServer := googlerpc.NewServer(googlerpc.EmptyServerOption{})
	agentServer := agent.NewServer(&agent.ServerOption{
		Getter: g,
		Logger: l,
	})
	cloud_agent.RegisterAgentServer(grpcServer, agentServer)

	go grpcServer.Serve(lis)

	return grpcServer.Stop, nil
}
