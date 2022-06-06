package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCmd(t *testing.T) {
	c := NewCmd(&options{})

	assert.Equal(t, 2, len(c.Commands()))
}
