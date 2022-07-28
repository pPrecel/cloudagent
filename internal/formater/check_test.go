package formater

import (
	"reflect"
	"testing"
	"time"

	cloud_agent "github.com/pPrecel/cloudagent/pkg/agent/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	timeNow  = time.Now()
	testResp = &cloud_agent.GardenerResponse{
		GeneralError: "error",
		ShootList: map[string]*cloud_agent.ShootList{
			"p1": {
				Error: "project error",
				Time:  timestamppb.New(timeNow),
			},
			"p2": {
				Time: timestamppb.New(timeNow),
			},
		},
	}
	testRespRows = [][]string{
		{"p1", "ERROR", "project error", testResp.ShootList["p1"].Time.AsTime().Local().Format("2006-01-02 15:04:05"), "Gardener"},
		{"p2", "OK", "OK", testResp.ShootList["p1"].Time.AsTime().Local().Format("2006-01-02 15:04:05"), "Gardener"},
	}
)

func Test_checkFormater_YAML(t *testing.T) {
	type fields struct {
		resp *cloud_agent.GardenerResponse
	}
	tests := []struct {
		name   string
		fields fields
		want   interface{}
	}{
		{
			name: "nil arguments",
			fields: fields{
				resp: nil,
			},
			want: map[string]interface{}{},
		},
		{
			name: "empty arguments",
			fields: fields{
				resp: &cloud_agent.GardenerResponse{},
			},
			want: map[string]interface{}{
				"generalError": "",
				"shoots":       mergeShoots(cloud_agent.GardenerResponse{}.ShootList),
			},
		},
		{
			name: "get json",
			fields: fields{
				resp: testResp,
			},
			want: map[string]interface{}{
				"generalError": "error",
				"shoots":       mergeShoots(testResp.ShootList),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := NewCheck(tt.fields.resp)
			if got := f.YAML(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("checkFormater.YAML() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_checkFormater_JSON(t *testing.T) {
	type fields struct {
		resp *cloud_agent.GardenerResponse
	}
	tests := []struct {
		name   string
		fields fields
		want   interface{}
	}{
		{
			name: "nil arguments",
			fields: fields{
				resp: nil,
			},
			want: map[string]interface{}{},
		},
		{
			name: "empty arguments",
			fields: fields{
				resp: &cloud_agent.GardenerResponse{},
			},
			want: map[string]interface{}{
				"generalError": "",
				"shoots":       mergeShoots(cloud_agent.GardenerResponse{}.ShootList),
			},
		},
		{
			name: "get json",
			fields: fields{
				resp: testResp,
			},
			want: map[string]interface{}{
				"generalError": "error",
				"shoots":       mergeShoots(testResp.ShootList),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &checkFormater{
				resp: tt.fields.resp,
			}
			if got := f.JSON(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("checkFormater.JSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_checkFormater_Table(t *testing.T) {
	type fields struct {
		resp *cloud_agent.GardenerResponse
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
				resp: testResp,
			},
			want:  checkHeaders,
			want1: testRespRows,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &checkFormater{
				resp: tt.fields.resp,
			}
			got, got1 := f.Table()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("checkFormater.Table() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("checkFormater.Table() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_checkFormater_Text(t *testing.T) {
	type fields struct {
		resp *cloud_agent.GardenerResponse
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
				resp: testResp,
			},
			args: args{
				outFormat: "$h $e $a $E",
				errFormat: "$E",
			},
			want: "1 1 2 error, project error",
		},
		{
			name: "nil resp",
			fields: fields{
				resp: nil,
			},
			args: args{
				outFormat: "$h $e $a $E",
				errFormat: "$E",
			},
			want: "nil response",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &checkFormater{
				resp: tt.fields.resp,
			}
			if got := f.Text(tt.args.outFormat, tt.args.errFormat); got != tt.want {
				t.Errorf("checkFormater.Text() = %v, want %v", got, tt.want)
			}
		})
	}
}
