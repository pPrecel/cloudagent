package gardener

import (
	"context"
	"errors"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/pPrecel/cloudagent/internal/gardener/automock"
	"github.com/pPrecel/cloudagent/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

func TestNewClusterConfig(t *testing.T) {
	tests := []struct {
		name           string
		kubeconfigPath string
		want           *rest.Config
		wantErr        bool
	}{
		{
			name:           "create config",
			kubeconfigPath: fixKubeconfigPath(t),
			want: &rest.Config{
				Host: "http://localhost:8080",
			},
			wantErr: false,
		},
		{
			name:           "path does not exist",
			kubeconfigPath: "/this/path/does/not/exist",
			want:           nil,
			wantErr:        true,
		},
		{
			name:           "empty kubeconfig",
			kubeconfigPath: fixEmptyKubeconfigPath(t),
			want:           nil,
			wantErr:        true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newClusterConfig(tt.kubeconfigPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClusterConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewClusterConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewClient(t *testing.T) {
	t.Run("create client", func(t *testing.T) {
		actualCfg, err := fixRestClient()
		assert.NoError(t, err)

		c, err := newShootClient(actualCfg, "")
		assert.NoError(t, err)
		assert.NotNil(t, c)
	})

	t.Run("client error", func(t *testing.T) {
		actualCfg, err := fixWrongRestClient()
		assert.NoError(t, err)

		c, err := newShootClient(actualCfg, "")
		assert.Error(t, err)
		assert.Nil(t, c)
	})
}

func fixKubeconfigPath(t *testing.T) string {
	_, filename, _, ok := runtime.Caller(0)
	assert.True(t, ok)

	return filepath.Join(filepath.Dir(filename), "/testdata/kubeconfig.yml")
}

func fixEmptyKubeconfigPath(t *testing.T) string {
	_, filename, _, ok := runtime.Caller(0)
	assert.True(t, ok)

	return filepath.Join(filepath.Dir(filename), "/testdata/empty_kubeconfig.yml")
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

func Test_shootClient_List(t *testing.T) {
	type fields struct {
		resourceInterface ResourceInterface
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.ShootList
		wantErr bool
	}{
		{
			name: "no resources",
			fields: fields{
				resourceInterface: func() ResourceInterface {
					ri := automock.NewResourceInterface(t)
					ri.On("List", mock.Anything, mock.Anything).Return(&unstructured.UnstructuredList{}, nil)
					return ri
				}(),
			},
			args: args{
				ctx: context.Background(),
			},
			want:    &types.ShootList{},
			wantErr: false,
		},
		{
			name: "list resources",
			fields: fields{
				resourceInterface: func() ResourceInterface {
					ri := automock.NewResourceInterface(t)
					ri.On("List", mock.Anything, mock.Anything).Return(&unstructured.UnstructuredList{
						Items: []unstructured.Unstructured{
							{}, {},
						},
					}, nil)
					return ri
				}(),
			},
			args: args{
				ctx: context.Background(),
			},
			want: &types.ShootList{
				Items: []types.Shoot{
					{}, {},
				},
			},
			wantErr: false,
		},
		{
			name: "client error",
			fields: fields{
				resourceInterface: func() ResourceInterface {
					ri := automock.NewResourceInterface(t)
					ri.On("List", mock.Anything, mock.Anything).Return(nil, errors.New("test error"))
					return ri
				}(),
			},
			args: args{
				ctx: context.Background(),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "client unexpected list response",
			fields: fields{
				resourceInterface: func() ResourceInterface {
					ri := automock.NewResourceInterface(t)
					ri.On("List", mock.Anything, mock.Anything).Return(&unstructured.UnstructuredList{
						Object: map[string]interface{}{
							"metadata": "wrong value",
						},
					}, nil)
					return ri
				}(),
			},
			args: args{
				ctx: context.Background(),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "client unexpected item list",
			fields: fields{
				resourceInterface: func() ResourceInterface {
					ri := automock.NewResourceInterface(t)
					ri.On("List", mock.Anything, mock.Anything).Return(&unstructured.UnstructuredList{
						Items: []unstructured.Unstructured{
							{
								Object: map[string]interface{}{
									"metadata": "wrong value",
								},
							},
						},
					}, nil)
					return ri
				}(),
			},
			args: args{
				ctx: context.Background(),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sc := &shootClient{
				resourceInterface: tt.fields.resourceInterface,
			}
			got, err := sc.List(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("shootClient.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("shootClient.List() = %v, want %v", got, tt.want)
			}
		})
	}
}
