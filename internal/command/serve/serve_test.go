package serve

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/pPrecel/cloud-agent/internal/agent"
	"github.com/pPrecel/cloud-agent/internal/command"
	"github.com/pPrecel/cloud-agent/internal/gardener"
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
		assert.Equal(t, "", o.kubeconfigPath)
		assert.Equal(t, "", o.namespace)
		assert.Equal(t, "@every 15m", o.cronSpec)
	})

	t.Run("parse flags", func(t *testing.T) {
		c.ParseFlags([]string{
			"--kubeconfigPath", "path",
			"--namespace", "namespace",
			"--cronSpec", "@every 15m",
		})

		assert.Equal(t, "path", o.kubeconfigPath)
		assert.Equal(t, "namespace", o.namespace)
		assert.Equal(t, "@every 15m", o.cronSpec)
	})

	t.Run("parse shortcuts", func(t *testing.T) {
		c.ParseFlags([]string{
			"-k", "other-path",
			"-n", "other-namespace",
			"-c", "@every 20m",
		})

		assert.Equal(t, "other-path", o.kubeconfigPath)
		assert.Equal(t, "other-namespace", o.namespace)
		assert.Equal(t, "@every 20m", o.cronSpec)
	})
}

var (
	testNetwork = "unix"
	testAddress = filepath.Join(os.TempDir(), "serve-test-socket.sock")
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
			newClusterConfig: func(s string) (*rest.Config, error) {
				return fixRestClient()
			},
			newWatchFunc: func(l *logrus.Logger, c gardener.Client, s gardener.StateSetter) agent.WatchFn {
				return func(ctx context.Context) {}
			},
		}
		c := NewCmd(o)
		c.ParseFlags([]string{
			"--kubeconfigPath", "/path",
			"--namespace", "namespace",
		})

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
			newClusterConfig: func(s string) (*rest.Config, error) {
				return nil, errors.New("test error")
			},
		}
		c := NewCmd(o)
		c.ParseFlags([]string{
			"--kubeconfigPath", "/path",
			"--namespace", "namespace",
		})

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
			newClusterConfig: func(s string) (*rest.Config, error) {
				return fixWrongRestClient()
			},
		}
		c := NewCmd(o)
		c.ParseFlags([]string{
			"--kubeconfigPath", "/path",
			"--namespace", "namespace",
		})

		err := c.PreRunE(c, []string{})
		assert.NoError(t, err)

		assert.Error(t, c.RunE(c, []string{}))
	})

	t.Run("wrong WatchFn type error", func(t *testing.T) {
		o := &options{
			Options: &command.Options{
				Logger:  l,
				Context: context.Background(),
			},
			socketNetwork: testNetwork,
			socketAddress: testAddress,
			newClusterConfig: func(s string) (*rest.Config, error) {
				return fixRestClient()
			},
			newWatchFunc: func(l *logrus.Logger, c gardener.Client, s gardener.StateSetter) agent.WatchFn {
				return func(ctx context.Context) {}
			},
		}
		c := NewCmd(o)
		c.ParseFlags([]string{
			"--kubeconfigPath", "/path",
			"--namespace", "namespace",
		})

		err := c.PreRunE(c, []string{})
		assert.NoError(t, err)

		o.cronSpec = ""
		assert.Error(t, c.RunE(c, []string{}))
	})

	t.Run("wrong WatchFn type error", func(t *testing.T) {
		o := &options{
			Options: &command.Options{
				Logger:  l,
				Context: context.Background(),
			},
			socketNetwork: testNetwork,
			socketAddress: testAddress,
			newClusterConfig: func(s string) (*rest.Config, error) {
				return fixRestClient()
			},
			newWatchFunc: func(l *logrus.Logger, c gardener.Client, s gardener.StateSetter) agent.WatchFn {
				return func(ctx context.Context) {}
			},
		}
		c := NewCmd(o)
		c.ParseFlags([]string{
			"--kubeconfigPath", "/path",
			"--namespace", "namespace",
		})

		err := c.PreRunE(c, []string{})
		assert.NoError(t, err)

		o.socketAddress = "."
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
