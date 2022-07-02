package formater

import (
	"strconv"
	"strings"

	"github.com/pPrecel/cloudagent/internal/output"
	cloud_agent "github.com/pPrecel/cloudagent/pkg/agent/proto"
)

const (
	gardenerProvider               = "Gardener"
	createdByLabel                 = `gardener.cloud/created-by`
	GardenerTextAllFormat          = "$a"
	GardenerTextUnknownFormat      = "$u"
	GardenerTextHealthyFormat      = "$r"
	GardenerTextHibernatedFormat   = "$h"
	GardenerTextEmptyFormat        = "$e"
	GardenerTextEmptyUnknownFormat = "$x"
	GardenerTextErrorFormat        = "$E"
)

var (
	gardenerHeaders = []string{"PROJECT", "NAME", "CREATED BY", "CONDITION", "PROVIDER"}

	preDirectives = gardenerDirectiveMap{
		GardenerTextAllFormat: func(_ *cloud_agent.Shoot) bool {
			return true
		},
	}

	postDirectives = gardenerDirectiveMap{
		GardenerTextEmptyUnknownFormat: func(s *cloud_agent.Shoot) bool {
			return s.Condition == cloud_agent.Condition_EMPTY ||
				s.Condition == cloud_agent.Condition_UNKNOWN
		},
		GardenerTextUnknownFormat: func(s *cloud_agent.Shoot) bool {
			return s.Condition == cloud_agent.Condition_UNKNOWN
		},
		GardenerTextHealthyFormat: func(s *cloud_agent.Shoot) bool {
			return s.Condition == cloud_agent.Condition_HEALTHY
		},
		GardenerTextHibernatedFormat: func(s *cloud_agent.Shoot) bool {
			return s.Condition == cloud_agent.Condition_HIBERNATED
		},
		GardenerTextEmptyFormat: func(s *cloud_agent.Shoot) bool {
			return s.Condition == cloud_agent.Condition_EMPTY
		},
	}
)

var _ output.Formater = &gardenerFormater{}

type gardenerFormater struct {
	err     error
	filters Filters
	shoots  *cloud_agent.ShootList
}

func NewGardener(err error, shoots *cloud_agent.ShootList, filters Filters) output.Formater {
	return &gardenerFormater{
		err:     err,
		shoots:  shoots,
		filters: filters,
	}
}

func (f *gardenerFormater) YAML() interface{} {
	if f.err != nil {
		return map[string]interface{}{}
	}

	shoots := f.filters.filter(f.shoots)

	return map[string]interface{}{
		"shoots": shoots.Shoots,
	}
}

func (f *gardenerFormater) JSON() interface{} {
	if f.err != nil {
		return map[string]interface{}{}
	}

	shoots := f.filters.filter(f.shoots)

	return map[string]interface{}{
		"shoots": shoots.Shoots,
	}
}

func (f *gardenerFormater) Table() ([]string, [][]string) {
	rows := [][]string{}

	if f.err != nil {
		return gardenerHeaders, rows
	}

	shoots := f.filters.filter(f.shoots)

	for i := range shoots.Shoots {
		shoot := shoots.Shoots[i]

		if shoot == nil {
			continue
		}

		rows = append(rows, []string{
			shoot.Namespace,
			shoot.Name,
			shoot.Annotations[createdByLabel],
			shoot.Condition.String(),
			gardenerProvider,
		})
	}

	return gardenerHeaders, rows
}

func (f *gardenerFormater) Text(outFormat, errFormat string) string {
	if f.err != nil {
		return strings.ReplaceAll(errFormat, GardenerTextErrorFormat, f.err.Error())
	}

	shoots := f.shoots
	directives := preDirectives.run(shoots, map[string]int{})

	shoots = f.filters.filter(shoots)

	directives = postDirectives.run(shoots, directives)

	str := outFormat
	for key, val := range directives {
		str = strings.ReplaceAll(str, key, strconv.Itoa(val))
	}

	return str
}

type gardenerDirectiveMap map[string]func(*cloud_agent.Shoot) bool

func (d gardenerDirectiveMap) run(s *cloud_agent.ShootList, m map[string]int) map[string]int {
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
