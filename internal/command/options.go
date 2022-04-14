package command

import (
	"context"

	"github.com/sirupsen/logrus"
)

type Options struct {
	Ctx     context.Context
	Logger  *logrus.Logger
	Verbose bool
}
