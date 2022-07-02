package config

import (
	"testing"

	"github.com/pPrecel/cloudagent/pkg/config"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_newGardenerCmd(t *testing.T) {
	o := &gardenerOptions{}
	c := newGardenerCmd(o)

	t.Run("defaults", func(t *testing.T) {
		assert.Equal(t, "", o.namespace)
		assert.Equal(t, "", o.kubeconfig)
	})

	t.Run("validation errors", func(t *testing.T) {
		assert.Error(t, c.PreRunE(c, []string{"add"}))
		assert.Error(t, c.PreRunE(c, []string{"del"}))
		assert.Error(t, c.PreRunE(c, []string{""}))
		assert.Error(t, c.RunE(c, []string{""}))
	})

	t.Run("parse flags", func(t *testing.T) {
		c.ParseFlags([]string{
			"--namespace", "namespace",
			"--kubeconfig", "kubeconfig",
		})

		assert.Equal(t, "namespace", o.namespace)
		assert.Equal(t, "kubeconfig", o.kubeconfig)
	})

	t.Run("validate", func(t *testing.T) {
		assert.NoError(t, c.PreRunE(c, []string{"add"}))
		assert.NoError(t, c.PreRunE(c, []string{"del"}))
	})

	t.Run("parse shortcuts", func(t *testing.T) {
		c.ParseFlags([]string{
			"-n", "other-namespace",
			"-k", "other-kubeconfig",
		})

		assert.Equal(t, "other-namespace", o.namespace)
		assert.Equal(t, "other-kubeconfig", o.kubeconfig)
	})
}

func Test_runAddGardener(t *testing.T) {
	type args struct {
		o *gardenerOptions
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "read write config",
			args: args{
				o: &gardenerOptions{
					options: &options{
						readConfig: func(s string) (*config.Config, error) {
							return &config.Config{
								PersistentSpec: "any",
								GardenerProjects: []config.GardenerProject{
									{
										Namespace:      "any",
										KubeconfigPath: "any",
									},
								},
							}, nil
						},
						writeConfig: func(s string, i interface{}) error {
							c, ok := i.(*config.Config)
							assert.True(t, ok)

							assert.Equal(t, "any", c.PersistentSpec)
							assert.Equal(t, []config.GardenerProject{
								{
									Namespace:      "any",
									KubeconfigPath: "any",
								},
								{
									Namespace:      "",
									KubeconfigPath: "",
								},
							}, c.GardenerProjects)

							return nil
						},
					},
				},
			},
		},
		{
			name: "read error",
			args: args{
				o: &gardenerOptions{
					options: &options{
						readConfig: func(s string) (*config.Config, error) {
							return nil, errors.New("test error")
						},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newGardenerCmd(tt.args.o)

			if err := c.RunE(c, []string{"add"}); (err != nil) != tt.wantErr {
				t.Errorf("runAddGardener() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_runDelGardener(t *testing.T) {

	tests := []struct {
		name    string
		cmd     *cobra.Command
		wantErr bool
	}{
		{
			name: "read write config",
			cmd: func() *cobra.Command {
				o := &gardenerOptions{
					options: &options{
						readConfig: func(s string) (*config.Config, error) {
							return &config.Config{
								GardenerProjects: []config.GardenerProject{
									{
										Namespace:      "any",
										KubeconfigPath: "any",
									},
									{
										Namespace:      "any2",
										KubeconfigPath: "any2",
									},
								},
							}, nil
						},
						writeConfig: func(s string, i interface{}) error {
							return nil
						},
					},
				}
				c := newGardenerCmd(o)

				o.namespace = "any"

				return c
			}(),
			wantErr: false,
		},
		{
			name: "read error",
			cmd: func() *cobra.Command {
				return newGardenerCmd(&gardenerOptions{
					options: &options{
						readConfig: func(s string) (*config.Config, error) {
							return nil, errors.New("test error")
						},
					},
				})
			}(),
			wantErr: true,
		},
		{
			name: "no match error",
			cmd: func() *cobra.Command {
				return newGardenerCmd(&gardenerOptions{
					options: &options{
						readConfig: func(s string) (*config.Config, error) {
							return &config.Config{}, nil
						},
					},
				})
			}(),
			wantErr: true,
		},
		{
			name: "write error",
			cmd: func() *cobra.Command {
				o := &gardenerOptions{
					options: &options{
						readConfig: func(s string) (*config.Config, error) {
							return &config.Config{
								GardenerProjects: []config.GardenerProject{
									{
										Namespace:      "any",
										KubeconfigPath: "any",
									},
									{
										Namespace:      "any2",
										KubeconfigPath: "any2",
									},
								},
							}, nil
						},
						writeConfig: func(s string, i interface{}) error {
							return errors.New("test error")
						},
					},
				}
				c := newGardenerCmd(o)

				o.kubeconfig = "any"

				return c
			}(),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.cmd.RunE(tt.cmd, []string{"del"}); (err != nil) != tt.wantErr {
				t.Errorf("runDelGardener() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
