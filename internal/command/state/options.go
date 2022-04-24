package state

import (
	"io"
	"os"
	"time"

	"github.com/pPrecel/cloud-agent/internal/agent"
	"github.com/pPrecel/cloud-agent/internal/command"
	"github.com/pPrecel/cloud-agent/internal/output"
)

type options struct {
	*command.Options

	createdBy string
	outFormat output.Output
	timeout   time.Duration

	socketAddress string
	socketNetwork string
	writer        io.Writer
}

func NewOptions(opts *command.Options) *options {
	return &options{
		Options:       opts,
		socketAddress: agent.Address,
		socketNetwork: agent.Network,
		writer:        os.Stdout,
	}
}
