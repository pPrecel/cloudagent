package config

import (
	"bytes"
	"encoding/json"
	"errors"

	"github.com/spf13/cobra"
)

func newSchemaCmd(o *schemaOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "schema",
		Short: "Generate the config JSON schema.",
		Long:  "Use this command to generate a config JSON schema.",
		RunE: func(_ *cobra.Command, _ []string) error {
			return runSchema(o)
		},
	}
}

func runSchema(o *schemaOptions) error {
	b, err := o.jsonSchema()
	if err != nil {
		return errors.New("can't reflect schema")
	}

	var prettyJSON bytes.Buffer
	error := json.Indent(&prettyJSON, b, "", "    ")
	if error != nil {
		return errors.New("can't parse output json")
	}

	o.stdout.Write(prettyJSON.Bytes())
	return nil
}
