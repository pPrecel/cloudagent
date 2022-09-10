package agent

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/gardener/gardener/pkg/apis/core/v1beta1"
	cloud_agent "github.com/pPrecel/cloudagent/pkg/agent/proto"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
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
					Name:              "name2",
					Namespace:         "namespace1",
					CreationTimestamp: v1.NewTime(fixRFC3339Time("2022-09-10T01:08:00Z")),
				},
				Status: v1beta1.ShootStatus{
					Conditions: []v1beta1.Condition{
						{
							Status:             v1beta1.ConditionTrue,
							LastTransitionTime: v1.NewTime(fixRFC3339Time("2022-09-10T10:08:17Z")),
						},
						{
							Status:             v1beta1.ConditionFalse,
							LastTransitionTime: v1.NewTime(fixRFC3339Time("2022-09-10T10:02:00Z")),
						},
					},
				},
			},
			{},
		},
	}
	testGardenerShootList2 = &v1beta1.ShootList{
		Items: []v1beta1.Shoot{
			{
				ObjectMeta: v1.ObjectMeta{
					Name:      "name1",
					Namespace: "namespace2",
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
				Condition:          cloud_agent.Condition_HIBERNATED,
				LastTransitionTime: timestamppb.New(time.Time{}),
				CreationTimestamp:  timestamppb.New(time.Time{}),
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
				Condition:          cloud_agent.Condition_HEALTHY,
				LastTransitionTime: timestamppb.New(time.Time{}),
				CreationTimestamp:  timestamppb.New(time.Time{}),
			},
			{
				Name:               "name2",
				Namespace:          "namespace1",
				Condition:          cloud_agent.Condition_UNKNOWN,
				LastTransitionTime: timestamppb.New(fixRFC3339Time("2022-09-10T10:08:17Z")),
				CreationTimestamp:  timestamppb.New(fixRFC3339Time("2022-09-10T01:08:00Z")),
			},
			{
				LastTransitionTime: timestamppb.New(time.Time{}),
				CreationTimestamp:  timestamppb.New(time.Time{}),
			},
		},
	}
	testAgentShootList2 = &cloud_agent.ShootList{
		Error: "test error",
		Shoots: []*cloud_agent.Shoot{
			{
				Name:      "name1",
				Namespace: "namespace2",
				Labels: map[string]string{
					"label1": "val1",
					"label2": "val2",
				},
				Annotations: map[string]string{
					"annotation1": "val1",
					"annotation2": "val2",
				},
				Condition:          cloud_agent.Condition_HIBERNATED,
				LastTransitionTime: timestamppb.New(time.Time{}),
				CreationTimestamp:  timestamppb.New(time.Time{}),
			},
		},
	}
)

func Test_server_GardenerShoots(t *testing.T) {
	l := &logrus.Entry{
		Logger: logrus.New(),
	}
	l.Logger.Out = io.Discard

	type fields struct {
		gardenerCache *ServerCache
		logger        *logrus.Entry
	}
	tests := []struct {
		name    string
		fields  fields
		want    map[string]*cloud_agent.ShootList
		wantErr bool
	}{
		{
			name: "nil server cache",
			fields: fields{
				gardenerCache: nil,
				logger:        l,
			}, wantErr: true,
		},
		{
			name: "nil gardener cache",
			fields: fields{
				gardenerCache: &ServerCache{
					GardenerCache: nil,
				},
				logger: l,
			},
			wantErr: true,
		},
		{
			name: "state list",
			fields: fields{
				gardenerCache: fixShootListCache(testGardenerShootList),
				logger:        l,
			},
			want: map[string]*cloud_agent.ShootList{
				"test": testAgentShootList,
			},
			wantErr: false,
		},
		{
			name: "nil state list",
			fields: fields{
				gardenerCache: fixShootListCache(nil),
				logger:        l,
			},
			want: map[string]*cloud_agent.ShootList{
				"test": {
					Shoots: []*cloud_agent.Shoot{},
				},
			},
			wantErr: false,
		},
		{
			name: "multiple cache keys",
			fields: fields{
				gardenerCache: fixShootListCache2(),
				logger:        l,
			},
			want: map[string]*cloud_agent.ShootList{
				"test":  testAgentShootList,
				"test2": testAgentShootList2,
			},
			wantErr: false,
		},
		{
			name: "geenral error",
			fields: fields{
				gardenerCache: &ServerCache{
					GardenerCache: NewCache[*v1beta1.ShootList](),
					GeneralError:  errors.New("test error"),
				},
				logger: l,
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "empty cache",
			fields: fields{
				gardenerCache: fixShootListCache(&v1beta1.ShootList{}),
				logger:        l,
			},
			want: map[string]*cloud_agent.ShootList{
				"test": {
					Shoots: []*cloud_agent.Shoot{},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewServer(&ServerOption{
				Cache:  tt.fields.gardenerCache,
				Logger: tt.fields.logger,
			})
			got, err := s.GardenerShoots(context.Background(), &cloud_agent.Empty{})
			if (err != nil) != tt.wantErr {
				t.Errorf("server.GardenerShoots() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.want != nil {
				compareMaps(t, tt.want, got.ShootList)
			}
		})
	}
}

func compareMaps(t *testing.T, m1, m2 map[string]*cloud_agent.ShootList) {

	// check maps len
	assert.Equal(t, len(m1), len(m2))

	// check if maps are nil
	if m1 == nil {
		assert.Nil(t, m2)
		return
	}

	for key := range m1 {

		// if first map elem == nil then second should be nil
		if m1[key] == nil {
			assert.Nil(t, m2[key])
			continue
		}

		// if not then compare values
		assert.Equal(t, m1[key].Shoots, m2[key].Shoots)
		assert.Equal(t, m1[key].Error, m2[key].Error)
	}
}

func fixShootListCache(s *v1beta1.ShootList) *ServerCache {
	c := NewCache[*v1beta1.ShootList]()

	c.Clean()

	r := c.Register("test")
	r.Set(s, nil)

	return &ServerCache{
		GardenerCache: c,
	}
}

func fixShootListCache2() *ServerCache {
	c := fixShootListCache(testGardenerShootList)

	r := c.GardenerCache.Register("test2")

	// set value only
	r.Set(testGardenerShootList2, nil)

	// set error separately to not override test value (to keep both)
	r.Set(nil, errors.New("test error"))

	return c
}

func fixRFC3339Time(value string) time.Time {
	t, _ := time.Parse(time.RFC3339, value)
	return t
}
