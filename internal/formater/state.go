package formater

import (
	"strconv"
	"strings"

	"github.com/pPrecel/cloudagent/internal/output"
	cloud_agent "github.com/pPrecel/cloudagent/pkg/agent/proto"
)

const (
	createdByLabel       = `gardener.cloud/created-by`
	TextAllFormat        = "$a"
	TextUnknownFormat    = "$u"
	TextRunningFormat    = "$r"
	TextHibernatedFormat = "$h"
	TextErrorFormat      = "$e"
)

var (
	headers = []string{"NAME", "CREATED BY", "CONDITION"}
)

var _ output.Formater = &state{}

type Filters struct {
	CreatedBy string
}

type state struct {
	err     error
	filters Filters
	shoots  *cloud_agent.ShootList
}

func NewForState(err error, shoots *cloud_agent.ShootList, filters Filters) output.Formater {
	return &state{
		err:     err,
		shoots:  shoots,
		filters: filters,
	}
}

func (s *state) YAML() interface{} {
	if s.err != nil {
		return map[string]interface{}{}
	}

	shoots := s.shoots
	if s.filters.CreatedBy != "" {
		shoots = shootsCreatedBy(shoots, s.filters.CreatedBy)
	}

	return map[string]interface{}{
		"shoots": shoots.Shoots,
	}
}

func (s *state) JSON() interface{} {
	if s.err != nil {
		return map[string]interface{}{}
	}

	shoots := s.shoots
	if s.filters.CreatedBy != "" {
		shoots = shootsCreatedBy(shoots, s.filters.CreatedBy)
	}

	return map[string]interface{}{
		"shoots": shoots.Shoots,
	}
}

func (s *state) Table() ([]string, [][]string) {
	rows := [][]string{}

	if s.err != nil {
		return headers, rows
	}

	shoots := s.shoots
	if s.filters.CreatedBy != "" {
		shoots = shootsCreatedBy(shoots, s.filters.CreatedBy)
	}

	for i := range shoots.Shoots {
		shoot := shoots.Shoots[i]

		if shoot == nil {
			continue
		}

		rows = append(rows, []string{
			shoot.Name,
			shoot.Annotations[createdByLabel],
			shoot.Condition.String(),
		})
	}

	return headers, rows
}

func (s *state) Text(outFormat, errFormat string) string {
	if s.err != nil {
		return strings.ReplaceAll(errFormat, TextErrorFormat, s.err.Error())
	}

	str := outFormat

	l := len(s.shoots.Shoots)
	str = strings.ReplaceAll(str, TextAllFormat, strconv.Itoa(l))

	shoots := s.shoots
	if s.filters.CreatedBy != "" {
		shoots = shootsCreatedBy(shoots, s.filters.CreatedBy)
	}

	r := 0
	h := 0
	u := 0
	for i := range shoots.Shoots {
		if shoots.Shoots[i] == nil {
			continue
		}

		if shoots.Shoots[i].Condition == cloud_agent.Condition_HEALTHY {
			r++
		} else if shoots.Shoots[i].Condition == cloud_agent.Condition_HIBERNATED {
			h++
		} else if shoots.Shoots[i].Condition == cloud_agent.Condition_UNKNOWN {
			u++
		}
	}

	str = strings.ReplaceAll(str, TextRunningFormat, strconv.Itoa(r))
	str = strings.ReplaceAll(str, TextHibernatedFormat, strconv.Itoa(h))
	str = strings.ReplaceAll(str, TextUnknownFormat, strconv.Itoa(u))

	return str
}

func shootsCreatedBy(s *cloud_agent.ShootList, c string) *cloud_agent.ShootList {
	list := &cloud_agent.ShootList{}

	for i := range s.Shoots {
		if s.Shoots[i] == nil {
			continue
		}

		if s.Shoots[i].Annotations[createdByLabel] == c {
			list.Shoots = append(list.Shoots, s.Shoots[i])
		}
	}

	return list
}
