package config

import (
	"github.com/invopop/jsonschema"
)

func JSONSchema() ([]byte, error) {
	s := jsonschema.Reflect(&Config{})

	return s.MarshalJSON()
}
