package gardener

import (
	"context"
	"errors"
	"io/ioutil"
	"testing"

	"github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/pPrecel/cloud-agent/internal/gardener/automock"
	"github.com/pPrecel/cloud-agent/pkg/agent"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

func TestNewWatchFunc(t *testing.T) {
	l := logrus.New()
	l.Out = ioutil.Discard

	type args struct {
		c Client
	}
	tests := []struct {
		name string
		args args
		want *v1beta1.ShootList
	}{
		{
			name: "watch resources",
			args: args{
				c: func() Client {
					c := automock.NewClient(t)
					c.On("List", mock.Anything, v1.ListOptions{}).Return(shootList, nil).Once()

					return c
				}(),
			},
			want: shootList,
		},
		{
			name: "list error",
			args: args{
				c: func() Client {
					c := automock.NewClient(t)
					c.On("List", mock.Anything, v1.ListOptions{}).Return(nil, errors.New("new error")).Once()

					return c
				}(),
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := agent.NewCache[*v1beta1.ShootList]()
			r := c.Register("test-data")
			got := NewWatchFunc(l, tt.args.c, r)

			got(context.Background())
			assert.Equal(t, tt.want, r.Get())
		})
	}
}
