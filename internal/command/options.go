package command

import (
	"context"

	"github.com/sirupsen/logrus"
)

type Options struct {
	Context context.Context
	Logger  *logrus.Logger
	Verbose bool
}
