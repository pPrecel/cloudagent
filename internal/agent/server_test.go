package agent

import (
	"context"
	"reflect"
	"testing"

	"github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/pPrecel/cloud-agent/internal/agent/automock"
	cloud_agent "github.com/pPrecel/cloud-agent/internal/agent/proto"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	testGardenerShootList = &v1beta1.ShootList{
		Items: []v1beta1.Shoot{
			{
				ObjectMeta: v1.ObjectMeta{
					Name:      "name1",
					Namespace: "namespace1",
					Labels: map[string]string{
						"label1": "val1",
						"label2": "val2",
					},
					Annotations: map[string]string{
						"annotation1": "val1",
						"annotation2": "val2",
					},
				},
				Status: v1beta1.ShootStatus{
					IsHibernated: true,
				},
			},
			{
				ObjectMeta: v1.ObjectMeta{
					Name:      "name2",
					Namespace: "namespace1",
					Labels: map[string]string{
						"label1": "val1",
					},
					Annotations: map[string]string{
						"annotation2": "val2",
					},
				},
			},
			{
				ObjectMeta: v1.ObjectMeta{
					Name:      "name2",
					Namespace: "namespace1",
				},
				Status: v1beta1.ShootStatus{
					Conditions: []v1beta1.Condition{
						{
							Status: v1beta1.ConditionTrue,
						},
						{
							Status: v1beta1.ConditionFalse,
						},
					},
				},
			},
			{},
		},
	}
	testAgentShootList = &cloud_agent.ShootList{
		Shoots: []*cloud_agent.Shoot{
			{
				Name:      "name1",
				Namespace: "namespace1",
				Labels: map[string]string{
					"label1": "val1",
					"label2": "val2",
				},
				Annotations: map[string]string{
					"annotation1": "val1",
					"annotation2": "val2",
				},
				Condition: cloud_agent.Condition_HIBERNATED,
			},
			{
				Name:      "name2",
				Namespace: "namespace1",
				Labels: map[string]string{
					"label1": "val1",
				},
				Annotations: map[string]string{
					"annotation2": "val2",
				},
				Condition: cloud_agent.Condition_HEALTHY,
			},
			{
				Name:      "name2",
				Namespace: "namespace1",
				Condition: cloud_agent.Condition_UNKNOWN,
			},
			{},
		},
	}
)

func Test_server_GardenerShoots(t *testing.T) {
	type fields struct {
		getter StateGetter
		logger *logrus.Logger
	}
	tests := []struct {
		name    string
		fields  fields
		want    *cloud_agent.ShootList
		wantErr bool
	}{
		{
			name: "empty state list",
			fields: fields{
				getter: func() StateGetter {
					m := automock.NewStateGetter(t)
					m.On("Get").Return(&v1beta1.ShootList{}).Once()

					return m
				}(),
				logger: logrus.New(),
			},
			want:    &cloud_agent.ShootList{},
			wantErr: false,
		},
		{
			name: "state list",
			fields: fields{
				getter: func() StateGetter {
					m := automock.NewStateGetter(t)
					m.On("Get").Return(testGardenerShootList).Once()

					return m
				}(),
				logger: logrus.New(),
			},
			want:    testAgentShootList,
			wantErr: false,
		},
		{
			name: "nil state list",
			fields: fields{
				getter: func() StateGetter {
					m := automock.NewStateGetter(t)
					m.On("Get").Return(nil).Once()

					return m
				}(),
				logger: logrus.New(),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewServer(&ServerOption{
				Getter: tt.fields.getter,
				Logger: tt.fields.logger,
			})
			got, err := s.GardenerShoots(context.Background(), &cloud_agent.Empty{})
			if (err != nil) != tt.wantErr {
				t.Errorf("server.GardenerShoots() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("server.GardenerShoots() = %v, want %v", got, tt.want)
			}
		})
	}
}
