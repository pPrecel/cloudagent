package formater

import (
	"strconv"
	"strings"

	"github.com/pPrecel/cloudagent/internal/output"
	cloud_agent "github.com/pPrecel/cloudagent/pkg/agent/proto"
)

const (
	createdByLabel         = `gardener.cloud/created-by`
	TextAllFormat          = "$a"
	TextUnknownFormat      = "$u"
	TextHealthyFormat      = "$r"
	TextHibernatedFormat   = "$h"
	TextEmptyFormat        = "$e"
	TextEmptyUnknownFormat = "$x"
	TextErrorFormat        = "$E"
)

var (
	headers = []string{"NAME", "CREATED BY", "CONDITION"}

	preDirectives = directiveMap{
		TextAllFormat: func(_ *cloud_agent.Shoot) bool {
			return true
		},
	}

	postDirectives = directiveMap{
		TextEmptyUnknownFormat: func(s *cloud_agent.Shoot) bool {
			return s.Condition == cloud_agent.Condition_EMPTY ||
				s.Condition == cloud_agent.Condition_UNKNOWN
		},
		TextUnknownFormat: func(s *cloud_agent.Shoot) bool {
			return s.Condition == cloud_agent.Condition_UNKNOWN
		},
		TextHealthyFormat: func(s *cloud_agent.Shoot) bool {
			return s.Condition == cloud_agent.Condition_HEALTHY
		},
		TextHibernatedFormat: func(s *cloud_agent.Shoot) bool {
			return s.Condition == cloud_agent.Condition_HIBERNATED
		},
		TextEmptyFormat: func(s *cloud_agent.Shoot) bool {
			return s.Condition == cloud_agent.Condition_EMPTY
		},
	}
)

var _ output.Formater = &state{}

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

	shoots := s.filters.filter(s.shoots)

	return map[string]interface{}{
		"shoots": shoots.Shoots,
	}
}

func (s *state) JSON() interface{} {
	if s.err != nil {
		return map[string]interface{}{}
	}

	shoots := s.filters.filter(s.shoots)

	return map[string]interface{}{
		"shoots": shoots.Shoots,
	}
}

func (s *state) Table() ([]string, [][]string) {
	rows := [][]string{}

	if s.err != nil {
		return headers, rows
	}

	shoots := s.filters.filter(s.shoots)

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

	shoots := s.shoots
	directives := preDirectives.run(shoots, map[string]int{})

	shoots = s.filters.filter(shoots)

	directives = postDirectives.run(shoots, directives)

	str := outFormat
	for key, val := range directives {
		str = strings.ReplaceAll(str, key, strconv.Itoa(val))
	}

	return str
}

type directiveMap map[string]func(*cloud_agent.Shoot) bool

func (d directiveMap) run(s *cloud_agent.ShootList, m map[string]int) map[string]int {
	for key, val := range d {
		m[key] = 0

		for i := range s.Shoots {
			if s.Shoots[i] != nil && val(s.Shoots[i]) {
				m[key]++
			}
		}
	}

	return m
}

type Filters struct {
	CreatedBy string
}

func (f *Filters) filter(s *cloud_agent.ShootList) *cloud_agent.ShootList {
	if f.CreatedBy != "" {
		s = shootsCreatedBy(s, f.CreatedBy)
	}
	return s
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
