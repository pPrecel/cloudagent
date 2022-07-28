package state

import (
	"io"
	"os"
	"time"

	command "github.com/pPrecel/cloudagent/cmd"
	"github.com/pPrecel/cloudagent/internal/output"
	"github.com/pPrecel/cloudagent/pkg/agent"
)

type options struct {
	*command.Options

	createdBy string
	outFormat output.Flag
	timeout   time.Duration

	socketAddress string
	socketNetwork string
	writer        io.Writer
}

func NewOptions(opts *command.Options) *options {
	return &options{
		Options:       opts,
		socketNetwork: agent.Network,
		writer:        os.Stdout,
	}
}
