package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_newGardenerOptions(t *testing.T) {
	t.Run("validate add", func(t *testing.T) {
		o := newGardenerOptions(&options{})
		o.kubeconfig = "/any/path"
		o.namespace = "any-namespace"

		assert.NoError(t, o.validateAdd())
	})

	t.Run("validate add error", func(t *testing.T) {
		o := newGardenerOptions(&options{})
		o.namespace = "any-namespace"

		assert.Error(t, o.validateAdd())

		o = newGardenerOptions(&options{})
		o.kubeconfig = "/any/path"

		assert.Error(t, o.validateAdd())
	})

	t.Run("validate del", func(t *testing.T) {
		o := newGardenerOptions(&options{})
		o.namespace = "any-namespace"

		assert.NoError(t, o.validateDel())

		o = newGardenerOptions(&options{})
		o.kubeconfig = "/any/path"

		assert.NoError(t, o.validateDel())

		o = newGardenerOptions(&options{})
		o.kubeconfig = "/any/path"
		o.namespace = "any-namespace"

		assert.NoError(t, o.validateDel())
	})

	t.Run("validate del error", func(t *testing.T) {
		o := newGardenerOptions(&options{})

		assert.Error(t, o.validateDel())
	})
}
