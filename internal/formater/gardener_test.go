package formater

import (
	"testing"
	"time"

	cloud_agent "github.com/pPrecel/cloudagent/pkg/agent/proto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	testShoots = &cloud_agent.ShootList{
		Shoots: []*cloud_agent.Shoot{
			{
				Name: "test",
				Annotations: map[string]string{
					createdByLabel: "me",
				},
				Condition:          1,
				LastTransitionTime: timestamppb.New(fixRFC3339Time("2022-09-10T10:08:17Z")),
				CreationTimestamp:  timestamppb.New(fixRFC3339Time("2022-09-10T10:06:17Z")),
			},
			{
				Name:      "test2",
				Namespace: "test-namespace",
				Annotations: map[string]string{
					createdByLabel: "me2",
				},
				Condition:          2,
				LastTransitionTime: timestamppb.New(fixRFC3339Time("2022-09-10T10:10:10Z")),
				CreationTimestamp:  timestamppb.New(fixRFC3339Time("2022-09-10T10:08:10Z")),
			},
			{
				Name:      "test3",
				Namespace: "test-namespace",
				Annotations: map[string]string{
					createdByLabel: "me2",
				},
				Condition:          3,
				LastTransitionTime: timestamppb.New(fixRFC3339Time("2022-09-10T10:02:23Z")),
				CreationTimestamp:  timestamppb.New(fixRFC3339Time("2022-09-10T10:00:23Z")),
			},
			{
				Name: "test4",
				Annotations: map[string]string{
					createdByLabel: "me2",
				},
				LastTransitionTime: timestamppb.New(time.Time{}),
				CreationTimestamp:  timestamppb.New(time.Time{}),
			},
			nil,
			nil,
		},
	}

	testShootsCreatedBy = &cloud_agent.ShootList{
		Shoots: []*cloud_agent.Shoot{
			{
				Name: "test",
				Annotations: map[string]string{
					createdByLabel: "me",
				},
				Condition:          1,
				LastTransitionTime: timestamppb.New(fixRFC3339Time("2022-09-10T10:08:17Z")),
				CreationTimestamp:  timestamppb.New(fixRFC3339Time("2022-09-10T10:06:17Z")),
			},
		},
	}

	testRows = [][]string{
		{"", "test", "me", "HEALTHY", fixLocalTime(fixRFC3339Time("2022-09-10T10:08:17Z")), fixLocalTime(fixRFC3339Time("2022-09-10T10:06:17Z")), "Gardener"},
		{"test-namespace", "test2", "me2", "HIBERNATED", fixLocalTime(fixRFC3339Time("2022-09-10T10:10:10Z")), fixLocalTime(fixRFC3339Time("2022-09-10T10:08:10Z")), "Gardener"},
		{"test-namespace", "test3", "me2", "UNKNOWN", fixLocalTime(fixRFC3339Time("2022-09-10T10:02:23Z")), fixLocalTime(fixRFC3339Time("2022-09-10T10:00:23Z")), "Gardener"},
		{"", "test4", "me2", "EMPTY", fixLocalTime(time.Time{}), fixLocalTime(time.Time{}), "Gardener"},
	}

	testFilteredRows = [][]string{
		{"", "test", "me", "HEALTHY", fixLocalTime(fixRFC3339Time("2022-09-10T10:08:17Z")), fixLocalTime(fixRFC3339Time("2022-09-10T10:06:17Z")), "Gardener"},
	}
)

func Test_state_YAML(t *testing.T) {
	type fields struct {
		filters Filters
		shoots  map[string]*cloud_agent.ShootList
	}
	tests := []struct {
		name   string
		fields fields
		want   interface{}
	}{
		{
			name: "get yaml",
			fields: fields{
				shoots: map[string]*cloud_agent.ShootList{
					"test1": testShoots,
				},
			},
			want: map[string]interface{}{
				"shoots": testShoots.Shoots,
			},
		},
		{
			name: "with filters",
			fields: fields{
				shoots: map[string]*cloud_agent.ShootList{
					"test1": testShoots,
				},
				filters: Filters{
					CreatedBy: "me",
				},
			},
			want: map[string]interface{}{
				"shoots": testShootsCreatedBy.Shoots,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewGardener(tt.fields.shoots, tt.fields.filters)
			assert.Equal(t, tt.want, s.YAML())
		})
	}
}

func Test_state_JSON(t *testing.T) {
	type fields struct {
		filters Filters
		shoots  map[string]*cloud_agent.ShootList
	}
	tests := []struct {
		name   string
		fields fields
		want   interface{}
	}{
		{
			name: "get yaml",
			fields: fields{
				shoots: map[string]*cloud_agent.ShootList{
					"test1": testShoots,
				},
			},
			want: map[string]interface{}{
				"shoots": testShoots.Shoots,
			},
		},
		{
			name: "with filters",
			fields: fields{
				shoots: map[string]*cloud_agent.ShootList{
					"test1": testShoots,
				},
				filters: Filters{
					CreatedBy: "me",
				},
			},
			want: map[string]interface{}{
				"shoots": testShootsCreatedBy.Shoots,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewGardener(tt.fields.shoots, tt.fields.filters)
			assert.Equal(t, tt.want, s.JSON())
		})
	}
}

func Test_state_Table(t *testing.T) {
	type fields struct {
		filters Filters
		shoots  map[string]*cloud_agent.ShootList
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
		want1  [][]string
	}{
		{
			name: "get table",
			fields: fields{
				shoots: map[string]*cloud_agent.ShootList{
					"test1": testShoots,
				},
			},
			want:  gardenerHeaders,
			want1: testRows,
		},
		{
			name: "with filters",
			fields: fields{
				shoots: map[string]*cloud_agent.ShootList{
					"test1": testShoots,
				},
				filters: Filters{
					CreatedBy: "me",
				},
			},
			want:  gardenerHeaders,
			want1: testFilteredRows,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewGardener(tt.fields.shoots, tt.fields.filters)
			got, got1 := s.Table()
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.want1, got1)
		})
	}
}

func Test_state_Text(t *testing.T) {
	type fields struct {
		filters Filters
		shoots  map[string]*cloud_agent.ShootList
	}
	type args struct {
		outFormat string
		errFormat string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "get text",
			fields: fields{
				shoots: map[string]*cloud_agent.ShootList{
					"test1": testShoots,
				},
			},
			args: args{
				outFormat: "$r $h $u $a",
				errFormat: " $E ",
			},
			want: "1 1 1 4",
		},
		{
			name: "with filters",
			fields: fields{
				shoots: map[string]*cloud_agent.ShootList{
					"test1": testShoots,
				},
				filters: Filters{
					CreatedBy: "me2",
				},
			},
			args: args{
				outFormat: "$r $h $u $a $e $x",
				errFormat: " $E ",
			},
			want: "0 1 1 4 1 2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewGardener(tt.fields.shoots, tt.fields.filters)
			assert.Equal(t, tt.want, s.Text(tt.args.outFormat, tt.args.errFormat))
		})
	}
}

func fixRFC3339Time(value string) time.Time {
	t, _ := time.Parse(time.RFC3339, value)
	return t.Local()
}

func fixLocalTime(value time.Time) string {
	return value.Local().Format("2006-01-02 15:04:05")
}
