package agent

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSocket(t *testing.T) {
	t.Run("create socket", func(t *testing.T) {
		tmpDir, e := ioutil.TempDir(os.TempDir(), "test-socket-")
		assert.NoError(t, e)
		defer os.RemoveAll(tmpDir)

		sockAddress := filepath.Join(tmpDir, "file.sock")

		l, e := NewSocket(Network, sockAddress)
		assert.NoError(t, e)
		assert.NotNil(t, l)

		f, e := os.Stat(sockAddress)
		assert.NoError(t, e)
		assert.False(t, f.IsDir())
	})

	t.Run("removing error", func(t *testing.T) {
		l, e := NewSocket(Network, ".")
		assert.Error(t, e)
		assert.Nil(t, l)
	})

	t.Run("make error", func(t *testing.T) {
		l, e := NewSocket(Network, "/any/path/\n")
		assert.Error(t, e)
		assert.Nil(t, l)
	})
}
