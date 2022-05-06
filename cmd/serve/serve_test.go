package serve

import (
	"context"
	"io"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gardener/gardener/pkg/apis/core/v1beta1"
	command "github.com/pPrecel/cloudagent/cmd"
	"github.com/pPrecel/cloudagent/internal/gardener"
	"github.com/pPrecel/cloudagent/pkg/agent"
	"github.com/pPrecel/cloudagent/pkg/config"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

func TestNewCmd(t *testing.T) {
	o := NewOptions(&command.Options{})
	c := NewCmd(o)

	t.Run("defaults", func(t *testing.T) {
		assert.Equal(t, config.ConfigPath, o.configPath)
	})

	t.Run("parse flags", func(t *testing.T) {
		c.ParseFlags([]string{
			"--configPath", "path",
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

	fixConfigFn = func(s string) (*config.Config, error) {
		return &config.Config{
			PersistentSpec: "@every 2m",
			GardenerProjects: []config.GardenerProject{
				{
					Namespace:      "test-namespace",
					KubeconfigPath: "path",
				},
			},
		}, nil
	}
)

func Test_run(t *testing.T) {
	l := logrus.New()
	l.Out = io.Discard

	t.Run("validate and run", func(t *testing.T) {
		o := &options{
			Options: &command.Options{
				Logger:  l,
				Context: context.Background(),
			},
			socketNetwork: testNetwork,
			socketAddress: testAddress,
			newSocket:     agent.NewSocket,
			getConfig: func(s string) (*config.Config, error) {
				return &config.Config{}, nil
			},
			newClusterConfig: func(s string) (*rest.Config, error) {
				return fixRestClient()
			},
			newWatchFunc: func(l *logrus.Logger, c gardener.Client, s agent.RegisteredResource[*v1beta1.ShootList]) agent.WatchFn {
				return func(ctx context.Context) {}
			},
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

	t.Run("config error", func(t *testing.T) {
		o := &options{
			Options: &command.Options{
				Logger:  l,
				Context: context.Background(),
			},
			socketNetwork: testNetwork,
			socketAddress: testAddress,
			newSocket:     agent.NewSocket,
			getConfig: func(s string) (*config.Config, error) {
				return nil, errors.New("test error")
			},
		}
		c := NewCmd(o)

		err := c.PreRunE(c, []string{})
		assert.NoError(t, err)

		assert.Error(t, c.RunE(c, []string{}))
	})

	t.Run("cluster config error", func(t *testing.T) {
		o := &options{
			Options: &command.Options{
				Logger:  l,
				Context: context.Background(),
			},
			socketNetwork: testNetwork,
			socketAddress: testAddress,
			getConfig:     fixConfigFn,
			newSocket:     agent.NewSocket,
			newClusterConfig: func(s string) (*rest.Config, error) {
				return nil, errors.New("test error")
			},
		}
		c := NewCmd(o)

		err := c.PreRunE(c, []string{})
		assert.NoError(t, err)

		assert.Error(t, c.RunE(c, []string{}))
	})

	t.Run("wrong config error", func(t *testing.T) {
		o := &options{
			Options: &command.Options{
				Logger:  l,
				Context: context.Background(),
			},
			socketNetwork: testNetwork,
			socketAddress: testAddress,
			getConfig:     fixConfigFn,
			newClusterConfig: func(s string) (*rest.Config, error) {
				return fixWrongRestClient()
			},
		}
		c := NewCmd(o)

		err := c.PreRunE(c, []string{})
		assert.NoError(t, err)

		assert.Error(t, c.RunE(c, []string{}))
	})

	t.Run("wrong socket path", func(t *testing.T) {
		o := &options{
			Options: &command.Options{
				Logger:  l,
				Context: context.Background(),
			},
			newSocket: func(network, address string) (net.Listener, error) {
				return nil, errors.New("test error")
			},
			getConfig: fixConfigFn,
			newClusterConfig: func(s string) (*rest.Config, error) {
				return fixRestClient()
			},
			newWatchFunc: func(l *logrus.Logger, c gardener.Client, s agent.RegisteredResource[*v1beta1.ShootList]) agent.WatchFn {
				return func(ctx context.Context) {}
			},
		}
		c := NewCmd(o)

		err := c.PreRunE(c, []string{})
		assert.NoError(t, err)

		assert.Error(t, c.RunE(c, []string{}))
	})
}

func fixWrongRestClient() (*rest.Config, error) {
	client, err := fixRestClient()
	if err != nil {
		return nil, err
	}

	client.AuthProvider = &api.AuthProviderConfig{}
	client.ExecProvider = &api.ExecConfig{}

	return client, err
}

func fixRestClient() (*rest.Config, error) {
	config := createValidTestConfig()

	clientBuilder := clientcmd.NewNonInteractiveClientConfig(*config, "clean", &clientcmd.ConfigOverrides{
		ClusterInfo: api.Cluster{
			TLSServerName: "overridden-server-name",
		},
	}, nil)

	return clientBuilder.ClientConfig()
}

func createValidTestConfig() *api.Config {
	const (
		server = "https://anything.com:8080"
		token  = "the-token"
	)

	config := api.NewConfig()
	config.Clusters["clean"] = &api.Cluster{
		Server: server,
	}
	config.AuthInfos["clean"] = &api.AuthInfo{
		Token: token,
	}
	config.Contexts["clean"] = &api.Context{
		Cluster:  "clean",
		AuthInfo: "clean",
	}
	config.CurrentContext = "clean"

	return config
}
