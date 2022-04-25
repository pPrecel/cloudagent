package agent

import (
	"context"

	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

type WatchFn func(context.Context)

type WatcherOptions struct {
	Spec    string
	Context context.Context
	Logger  *logrus.Logger
}

func NewWatcher(opts WatcherOptions, fn ...WatchFn) (*cron.Cron, error) {
	cron := cron.New()

	context := context.Background()
	if opts.Context != nil {
		context = opts.Context
	}

	for i := range fn {
		_, err := cron.AddFunc(opts.Spec, func() {
			fn[i](context)
		})
		if err != nil {
			return nil, err
		}
	}

	return cron, nil
}
