package watcher

import (
	"io"
	"reflect"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func Test_onDemandWatcher_GetGardenerCache(t *testing.T) {
	l := &logrus.Entry{
		Logger: logrus.New(),
	}
	l.Logger.Out = io.Discard

	t.Run("build new watcher", func(t *testing.T) {
		assert.NotNil(t, NewOnDemand(&NewOnDemandOptions{
			Logger: l,
		}))
	})
}

func TestNewOnDemand(t *testing.T) {
	type args struct {
		o *NewOnDemandOptions
	}
	tests := []struct {
		name string
		args args
		want *onDemandWatcher
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewOnDemand(tt.args.o); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewOnDemand() = %v, want %v", got, tt.want)
			}
		})
	}
}
