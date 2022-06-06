package config

import (
	"errors"
	"io/ioutil"
	"testing"

	"github.com/pPrecel/cloudagent/pkg/config"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func Test_newSchemaCmd(t *testing.T) {
	o := newSchemaOptions(&options{})
	c := newSchemaCmd(o)

	assert.Equal(t, 0, len(c.Commands()))
}

func Test_runSchema(t *testing.T) {
	l := logrus.New()
	l.Out = ioutil.Discard

	t.Run("validate and generate plist", func(t *testing.T) {
		o := &schemaOptions{
			stdout:     ioutil.Discard,
			jsonSchema: config.JSONSchema,
			options:    &options{},
		}
		c := newSchemaCmd(o)

		err := c.RunE(c, []string{})
		assert.NoError(t, err)
	})

	t.Run("reflect error", func(t *testing.T) {
		err := runSchema(&schemaOptions{
			stdout:     ioutil.Discard,
			jsonSchema: func() ([]byte, error) { return []byte{}, errors.New("test error") },
			options:    &options{},
		})

		assert.Error(t, err)
	})

	t.Run("reflect error", func(t *testing.T) {
		err := runSchema(&schemaOptions{
			stdout:     ioutil.Discard,
			jsonSchema: func() ([]byte, error) { return []byte("test test"), nil },
			options:    &options{},
		})

		assert.Error(t, err)
	})
}
