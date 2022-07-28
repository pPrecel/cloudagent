package check

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

	stdout        io.Writer
	outFormat     output.Flag
	timeout       time.Duration
	socketAddress string
	socketNetwork string
}

func NewOptions(o *command.Options) *options {
	return &options{
		Options:       o,
		stdout:        os.Stdout,
		socketNetwork: agent.Network,
	}
}
