package agent

import (
	"context"
	"reflect"
	"testing"

	"github.com/gardener/gardener/pkg/apis/core/v1beta1"
	cloud_agent "github.com/pPrecel/cloud-agent/pkg/agent/proto"
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
				Status: v1beta1.ShootStatus{
					Conditions: []v1beta1.Condition{
						{
							Status: v1beta1.ConditionTrue,
						},
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
		gardenerCache Cache[*v1beta1.ShootList]
		logger        *logrus.Logger
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
				gardenerCache: fixShootListCache(&v1beta1.ShootList{}),
				logger:        logrus.New(),
			},
			want:    &cloud_agent.ShootList{},
			wantErr: false,
		},
		{
			name: "state list",
			fields: fields{
				gardenerCache: fixShootListCache(testGardenerShootList),
				logger:        logrus.New(),
			},
			want:    testAgentShootList,
			wantErr: false,
		},
		{
			name: "nil state list",
			fields: fields{
				gardenerCache: fixShootListCache(nil),
				logger:        logrus.New(),
			},
			want:    &cloud_agent.ShootList{},
			wantErr: false,
		},
		{
			name: "nil cache",
			fields: fields{
				gardenerCache: nil,
				logger:        logrus.New(),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewServer(&ServerOption{
				GardenerCache: tt.fields.gardenerCache,
				Logger:        tt.fields.logger,
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

func fixShootListCache(s *v1beta1.ShootList) Cache[*v1beta1.ShootList] {
	c := NewCache[*v1beta1.ShootList]()

	r := c.Register("test")
	r.Set(s)

	return c
}
