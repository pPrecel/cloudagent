package logs

import (
	"io"
	"os"

	"github.com/hpcloud/tail"
	command "github.com/pPrecel/cloudagent/cmd"
)

type options struct {
	*command.Options

	filePath   string
	followLogs bool

	tailFile func(string, tail.Config) (*tail.Tail, error)
	writer   io.Writer
}

func NewOptions(o *command.Options) *options {
	return &options{
		tailFile: tail.TailFile,
		writer:   os.Stdout,
	}
}
