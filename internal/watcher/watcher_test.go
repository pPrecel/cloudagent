package watcher

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/pPrecel/cloudagent/internal/system"
	"github.com/pPrecel/cloudagent/pkg/agent"
	"github.com/pPrecel/cloudagent/pkg/config"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

var (
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

func TestNewWatcher(t *testing.T) {
	t.Run("new watcher", func(t *testing.T) {
		assert.NotNil(t, NewWatcher())
	})
}

func Test_watcher_Start(t *testing.T) {
	l := logrus.New()
	l.Out = io.Discard

	type fields struct {
		getConfig        func(string) (*config.Config, error)
		newClusterConfig func(string) (*rest.Config, error)
		notifyChange     func(string) (*system.Notifier, error)
	}
	tests := []struct {
		name    string
		fields  fields
		args    *Options
		wantErr bool
	}{
		{
			name: "notify change",
			fields: fields{
				getConfig:        fixConfigFn,
				newClusterConfig: fixRestClient,
				notifyChange: func(s string) (*system.Notifier, error) {
					n := &system.Notifier{
						IsMotified: make(chan interface{}),
						Errors:     make(chan error),
						Stop:       func() {},
					}
					go func() {
						n.IsMotified <- 1
					}()

					return n, nil
				},
			},
			args: &Options{
				Context:    context.Background(),
				Logger:     l,
				ConfigPath: "",
				Cache:      agent.NewCache[*v1beta1.ShootList](),
			},
			wantErr: false,
		},
		{
			name: "notify change error",
			fields: fields{
				getConfig:        fixConfigFn,
				newClusterConfig: fixRestClient,
				notifyChange: func(s string) (*system.Notifier, error) {
					n := &system.Notifier{
						IsMotified: make(chan interface{}),
						Errors:     make(chan error),
						Stop:       func() {},
					}
					go func() {
						n.Errors <- errors.New("test error")
					}()

					return n, nil
				},
			},
			args: &Options{
				Context:    context.Background(),
				Logger:     l,
				ConfigPath: "",
				Cache:      agent.NewCache[*v1beta1.ShootList](),
			},
			wantErr: true,
		},
		{
			name: "getConfig error",
			fields: fields{
				getConfig: func(s string) (*config.Config, error) {
					return nil, errors.New("test error")
				},
			},
			args: &Options{
				Context:    context.Background(),
				Logger:     l,
				ConfigPath: "",
				Cache:      agent.NewCache[*v1beta1.ShootList](),
			},
			wantErr: true,
		},
		{
			name: "newClusterConfig error",
			fields: fields{
				getConfig: fixConfigFn,
				newClusterConfig: func(s string) (*rest.Config, error) {
					return nil, errors.New("test error")
				},
			},
			args: &Options{
				Context:    context.Background(),
				Logger:     l,
				ConfigPath: "",
				Cache:      agent.NewCache[*v1beta1.ShootList](),
			},
			wantErr: true,
		},
		{
			name: "wrong rest config",
			fields: fields{
				getConfig:        fixConfigFn,
				newClusterConfig: fixWrongRestClient,
			},
			args: &Options{
				Context:    context.Background(),
				Logger:     l,
				ConfigPath: "",
				Cache:      agent.NewCache[*v1beta1.ShootList](),
			},
			wantErr: true,
		},
		{
			name: "notifyChange error",
			fields: fields{
				getConfig:        fixConfigFn,
				newClusterConfig: fixRestClient,
				notifyChange: func(s string) (*system.Notifier, error) {
					return nil, errors.New("test error")
				},
			},
			args: &Options{
				Context:    context.Background(),
				Logger:     l,
				ConfigPath: "",
				Cache:      agent.NewCache[*v1beta1.ShootList](),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &watcher{
				getConfig:        tt.fields.getConfig,
				newClusterConfig: tt.fields.newClusterConfig,
				notifyChange:     tt.fields.notifyChange,
			}

			if err := w.Start(tt.args); (err != nil) != tt.wantErr {
				t.Errorf("watcher.Start() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func fixWrongRestClient(_ string) (*rest.Config, error) {
	client, err := fixRestClient("")
	if err != nil {
		return nil, err
	}

	client.AuthProvider = &api.AuthProviderConfig{}
	client.ExecProvider = &api.ExecConfig{}

	return client, err
}

func fixRestClient(_ string) (*rest.Config, error) {
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
