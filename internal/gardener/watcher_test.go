package gardener

import (
	"context"
	"errors"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/pPrecel/cloudagent/internal/gardener/automock"
	"github.com/pPrecel/cloudagent/pkg/agent"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

var (
	shootList = &v1beta1.ShootList{
		Items: []v1beta1.Shoot{
			{
				ObjectMeta: v1.ObjectMeta{
					Name:      "name1",
					Namespace: "namespace1",
				},
			},
			{
				ObjectMeta: v1.ObjectMeta{
					Name:      "name2",
					Namespace: "namespace1",
				},
			},
			{
				ObjectMeta: v1.ObjectMeta{
					Name:      "name3",
					Namespace: "namespace3",
				},
			},
		},
	}
)

func Test_newWatchFunc(t *testing.T) {
	l := logrus.New()
	l.Out = ioutil.Discard

	t.Run("Fn not nil", func(t *testing.T) {
		c := NewWatchFunc(l, nil, "", "")
		assert.NotNil(t, c)
	})

	type args struct {
		l             *logrus.Logger
		r             agent.RegisteredResource[*v1beta1.ShootList]
		clientBuilder func() (Client, error)
	}
	tests := []struct {
		name string
		args args
		want *v1beta1.ShootList
	}{
		{
			name: "list resources",
			args: args{
				l: l,
				r: agent.NewCache[*v1beta1.ShootList]().Register("test"),
				clientBuilder: func() (Client, error) {
					c := automock.NewClient(t)
					c.On("List", mock.Anything, v1.ListOptions{}).Return(shootList, nil).Once()

					return c, nil
				},
			},
			want: shootList,
		},
		{
			name: "list error",
			args: args{
				l: l,
				r: agent.NewCache[*v1beta1.ShootList]().Register("test"),
				clientBuilder: func() (Client, error) {
					c := automock.NewClient(t)
					c.On("List", mock.Anything, v1.ListOptions{}).Return(nil, errors.New("test error")).Once()

					return c, nil
				},
			},
			want: nil,
		},
		{
			name: "client built error",
			args: args{
				l: l,
				r: agent.NewCache[*v1beta1.ShootList]().Register("test"),
				clientBuilder: func() (Client, error) {
					return nil, errors.New("test error")
				},
			},
			want: nil,
		},
		{
			name: "first client built error",
			args: args{
				l: l,
				r: agent.NewCache[*v1beta1.ShootList]().Register("test"),
				clientBuilder: func() func() (Client, error) {
					// return error on first run only
					i := 0
					return func() (Client, error) {
						if i == 0 {
							i++
							return nil, errors.New("test error")
						}

						c := automock.NewClient(t)
						c.On("List", mock.Anything, v1.ListOptions{}).Return(shootList, nil).Once()

						return c, nil
					}
				}(),
			},
			want: shootList,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newWatchFunc(tt.args.l, tt.args.r, tt.args.clientBuilder)(context.Background())

			got := tt.args.r.Get()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newWatchFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newClientBuilder(t *testing.T) {
	l := logrus.New()
	l.Out = ioutil.Discard

	type args struct {
		buildConfig func(string) (*rest.Config, error)
		namespace   string
		kubeconfig  string
	}
	tests := []struct {
		name       string
		args       args
		wantClient bool
		wantErr    bool
	}{
		{
			name: "return client",
			args: args{
				buildConfig: newClusterConfig,
				namespace:   "namespace1",
				kubeconfig:  fixKubeconfigPath(t),
			},
			wantClient: true,
			wantErr:    false,
		},
		{
			name: "cluster config error",
			args: args{
				buildConfig: newClusterConfig,
				namespace:   "namespace1",
				kubeconfig:  fixEmptyKubeconfigPath(t),
			},
			wantClient: false,
			wantErr:    true,
		},
		{
			name: "cluster client error",
			args: args{
				buildConfig: func(s string) (*rest.Config, error) {
					c, err := fixWrongRestClient()
					assert.NoError(t, err)
					return c, nil
				},
				namespace:  "namespace1",
				kubeconfig: fixKubeconfigPath(t),
			},
			wantClient: false,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newClient := newClientBuilder(l, tt.args.buildConfig, tt.args.namespace, tt.args.kubeconfig)
			got, err := newClient()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got != nil) != tt.wantClient {
				t.Errorf("NewClient() client = %v, wantClient %v", err, tt.wantErr)
				return
			}
		})
	}
}
