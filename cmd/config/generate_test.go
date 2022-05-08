package config

import (
	"testing"

	command "github.com/pPrecel/cloudagent/cmd"
	"github.com/stretchr/testify/assert"
)

func TestNewCmd(t *testing.T) {
	c := NewCmd(&command.Options{})

	assert.Equal(t, 1, len(c.Commands()))
}
