package formater

import (
	"strconv"
	"strings"

	"github.com/pPrecel/cloudagent/internal/output"
	"github.com/pPrecel/cloudagent/pkg/config"
)

const (
	ConfigTextErrorFormat      = "$e"
	ConfigTextAllFormat        = "$a"
	ConfigTextPersistentFormat = "$p"
	ConfigTextGardenerFormat   = "$g"
	ConfigTextGCPFormat        = "$G"
)

var (
	configHeaders = []string{"PROJECT", "CREDENTIALS"}

	configDirectives = configDirectiveMap{
		ConfigTextAllFormat: func(c config.Config) string {
			return strconv.Itoa(len(c.GCPProjects) + len(c.GardenerProjects))
		},
		ConfigTextPersistentFormat: func(c config.Config) string {
			return c.PersistentSpec
		},
		ConfigTextGardenerFormat: func(c config.Config) string {
			return strconv.Itoa(len(c.GardenerProjects))
		},
		ConfigTextGCPFormat: func(c config.Config) string {
			return strconv.Itoa(len(c.GCPProjects))
		},
	}
)

var _ output.Formater = &configFormater{}

type configDirectiveMap map[string]func(config.Config) string

type configFormater struct {
	err error
	cfg *config.Config
}

func NewConfig(err error, cfg *config.Config) output.Formater {
	return &configFormater{
		err: err,
		cfg: cfg,
	}
}

func (f *configFormater) YAML() interface{} {
	if f.err != nil || f.cfg == nil {
		return map[string]interface{}{}
	}

	return f.cfg
}

func (f *configFormater) JSON() interface{} {
	if f.err != nil || f.cfg == nil {
		return map[string]interface{}{}
	}

	return f.cfg
}

func (f *configFormater) Table() ([]string, [][]string) {
	if f.err != nil || f.cfg == nil {
		return configHeaders, [][]string{}
	}

	rows := [][]string{}
	for i := range f.cfg.GardenerProjects {
		p := f.cfg.GardenerProjects[i]

		rows = append(rows, []string{
			p.Namespace,
			p.KubeconfigPath,
		})
	}

	return configHeaders, rows
}

func (f *configFormater) Text(outFormat, errFormat string) string {
	if f.err != nil {
		return strings.ReplaceAll(errFormat, ConfigTextErrorFormat, f.err.Error())
	}

	if f.cfg == nil {
		err := "nil config"
		return strings.ReplaceAll(errFormat, ConfigTextErrorFormat, err)
	}

	str := outFormat
	for key, val := range configDirectives {
		str = strings.ReplaceAll(str, key, val(*f.cfg))
	}

	return str
}
