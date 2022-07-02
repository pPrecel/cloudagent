package config

import (
	"context"
	"io"
	"testing"

	command "github.com/pPrecel/cloudagent/cmd"
	"github.com/pPrecel/cloudagent/pkg/config"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestNewCmd(t *testing.T) {
	o := &options{}
	c := NewCmd(o)

	t.Run("defaults", func(t *testing.T) {
		assert.Equal(t, 2, len(c.Commands()))
		assert.Equal(t, config.ConfigPath, o.configPath)
	})

	t.Run("parse flags", func(t *testing.T) {
		c.ParseFlags([]string{
			"--config-path", "path",
		})

		assert.Equal(t, "path", o.configPath)
	})

	t.Run("parse shortcuts", func(t *testing.T) {
		c.ParseFlags([]string{
			"-c", "other-path",
		})

		assert.Equal(t, "other-path", o.configPath)
	})
}

type mockWriter struct {
	val    []byte
	rInt   int
	rError error
}

func (w *mockWriter) Write(b []byte) (int, error) {
	w.val = append(w.val, b...)
	return w.rInt, w.rError
}

func Test_run(t *testing.T) {
	l := logrus.New()
	l.Out = io.Discard

	type args struct {
		o *options
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "print config",
			args: args{
				o: func() *options {
					o := NewOptions(&command.Options{
						Context: context.Background(),
						Logger:  l,
					})

					o.readConfig = func(s string) (*config.Config, error) {
						return &config.Config{}, nil
					}
					o.stdout = io.Discard

					return o
				}(),
			},
			wantErr: false,
		},
		{
			name: "read config error",
			args: args{
				o: func() *options {
					o := NewOptions(&command.Options{
						Context: context.Background(),
						Logger:  l,
					})

					o.readConfig = func(s string) (*config.Config, error) {
						return &config.Config{}, errors.New("test error")
					}
					o.stdout = io.Discard

					return o
				}(),
			},
			wantErr: true,
		},
		{
			name: "marshal error",
			args: args{
				o: func() *options {
					o := NewOptions(&command.Options{
						Context: context.Background(),
						Logger:  l,
					})

					o.readConfig = func(_ string) (*config.Config, error) {
						return &config.Config{}, nil
					}
					o.marshal = func(_ interface{}) ([]byte, error) {
						return []byte{}, errors.New("test error")
					}
					o.stdout = io.Discard

					return o
				}(),
			},
			wantErr: true,
		},
		{
			name: "write error",
			args: args{
				o: func() *options {
					o := NewOptions(&command.Options{
						Context: context.Background(),
						Logger:  l,
					})

					o.readConfig = func(_ string) (*config.Config, error) {
						return &config.Config{}, nil
					}
					o.stdout = &mockWriter{
						rError: errors.New("test error"),
					}

					return o
				}(),
			},
			wantErr: false, // expected
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCmd(tt.args.o)

			if err := c.RunE(c, nil); (err != nil) != tt.wantErr {
				t.Errorf("run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
