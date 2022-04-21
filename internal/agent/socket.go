package agent

import (
	"net"
	"os"
	"path/filepath"
)

func NewSocket(network, address string) (net.Listener, error) {
	err := os.RemoveAll(address)
	if err != nil {
		return nil, err
	}

	err = os.MkdirAll(filepath.Dir(address), os.ModePerm)
	if err != nil {
		return nil, err
	}

	return net.Listen(network, address)
}
