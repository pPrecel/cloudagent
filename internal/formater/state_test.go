package formater

import (
	"errors"
	"testing"

	cloud_agent "github.com/pPrecel/cloudagent/pkg/agent/proto"
	"github.com/stretchr/testify/assert"
)

var (
	testShoots = &cloud_agent.ShootList{
		Shoots: []*cloud_agent.Shoot{
			{
				Name: "test",
				Annotations: map[string]string{
					createdByLabel: "me",
				},
				Condition: 1,
			},
			{
				Name:      "test2",
				Namespace: "test-namespace",
				Annotations: map[string]string{
					createdByLabel: "me2",
				},
				Condition: 2,
			},
			{
				Name:      "test3",
				Namespace: "test-namespace",
				Annotations: map[string]string{
					createdByLabel: "me2",
				},
				Condition: 3,
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
				Condition: 1,
			},
		},
	}

	testRows = [][]string{
		{"test", "me", "HEALTHY"},
		{"test2", "me2", "HIBERNATED"},
		{"test3", "me2", "UNKNOWN"},
	}

	testFilteredRows = [][]string{
		{"test", "me", "HEALTHY"},
	}
)

func Test_state_YAML(t *testing.T) {
	type fields struct {
		err     error
		filters Filters
		shoots  *cloud_agent.ShootList
	}
	tests := []struct {
		name   string
		fields fields
		want   interface{}
	}{
		{
			name: "get yaml",
			fields: fields{
				shoots: testShoots,
			},
			want: map[string]interface{}{
				"shoots": testShoots.Shoots,
			},
		},
		{
			name: "with error",
			fields: fields{
				err: errors.New("error"),
			},
			want: map[string]interface{}{},
		},
		{
			name: "with filters",
			fields: fields{
				shoots: testShoots,
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
			s := NewForState(tt.fields.err, tt.fields.shoots, tt.fields.filters)
			assert.Equal(t, tt.want, s.YAML())
		})
	}
}

func Test_state_JSON(t *testing.T) {
	type fields struct {
		err     error
		filters Filters
		shoots  *cloud_agent.ShootList
	}
	tests := []struct {
		name   string
		fields fields
		want   interface{}
	}{
		{
			name: "get yaml",
			fields: fields{
				shoots: testShoots,
			},
			want: map[string]interface{}{
				"shoots": testShoots.Shoots,
			},
		},
		{
			name: "with error",
			fields: fields{
				err: errors.New("error"),
			},
			want: map[string]interface{}{},
		},
		{
			name: "with filters",
			fields: fields{
				shoots: testShoots,
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
			s := NewForState(tt.fields.err, tt.fields.shoots, tt.fields.filters)
			assert.Equal(t, tt.want, s.JSON())
		})
	}
}

func Test_state_Table(t *testing.T) {
	type fields struct {
		err     error
		filters Filters
		shoots  *cloud_agent.ShootList
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
				shoots: testShoots,
			},
			want:  headers,
			want1: testRows,
		},
		{
			name: "with error",
			fields: fields{
				err: errors.New("test"),
			},
			want:  headers,
			want1: [][]string{},
		},
		{
			name: "with filters",
			fields: fields{
				shoots: testShoots,
				filters: Filters{
					CreatedBy: "me",
				},
			},
			want:  headers,
			want1: testFilteredRows,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewForState(tt.fields.err, tt.fields.shoots, tt.fields.filters)
			got, got1 := s.Table()
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.want1, got1)
		})
	}
}

func Test_state_Text(t *testing.T) {
	type fields struct {
		err     error
		filters Filters
		shoots  *cloud_agent.ShootList
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
			name: "get table",
			fields: fields{
				shoots: testShoots,
			},
			args: args{
				outFormat: "$r $h $u $a",
				errFormat: " $e ",
			},
			want: "1 1 1 5",
		},
		{
			name: "with error",
			fields: fields{
				err: errors.New("test error"),
			},
			args: args{
				outFormat: "$r $h $u $a",
				errFormat: "$e.",
			},
			want: "test error.",
		},
		{
			name: "with filters",
			fields: fields{
				shoots: testShoots,
				filters: Filters{
					CreatedBy: "me2",
				},
			},
			args: args{
				outFormat: "$r $h $u $a",
				errFormat: " $e ",
			},
			want: "0 1 1 5",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewForState(tt.fields.err, tt.fields.shoots, tt.fields.filters)
			assert.Equal(t, tt.want, s.Text(tt.args.outFormat, tt.args.errFormat))
		})
	}
}
