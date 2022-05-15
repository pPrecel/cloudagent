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

type Watcher struct {
	c *cron.Cron
}

func NewWatcher(opts WatcherOptions, fn ...WatchFn) (*Watcher, error) {
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

	return &Watcher{
		c: cron,
	}, nil
}

func (w *Watcher) Start() {
	go func() {
		e := w.c.Entries()
		for i := range e {
			e[i].Job.Run()
		}
	}()

	w.c.Start()
}

func (w *Watcher) Stop() context.Context {
	return w.c.Stop()
}
