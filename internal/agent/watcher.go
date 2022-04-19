package agent

import (
	"context"

	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

const defaultSpec = "@every 15m"

type WatchFn func(context.Context)

type WatcherOptions struct {
	Spec    string
	Context context.Context
	Logger  *logrus.Logger
}

func NewWatcher(opts WatcherOptions, fn ...WatchFn) (*cron.Cron, error) {
	cron := cron.New()
	spec := defaultSpec
	if opts.Spec != "" {
		spec = opts.Spec
	}

	context := context.Background()
	if opts.Context != nil {
		context = opts.Context
	}

	for i := range fn {
		_, err := cron.AddFunc(spec, func() {
			fn[i](context)
		})
		if err != nil {
			return nil, err
		}
	}

	return cron, nil
}
