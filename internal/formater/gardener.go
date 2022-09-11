package formater

import (
	"strconv"
	"strings"
	"time"

	"github.com/pPrecel/cloudagent/internal/output"
	cloud_agent "github.com/pPrecel/cloudagent/pkg/agent/proto"
	"k8s.io/apimachinery/pkg/labels"
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
	gardenerHeaders = []string{"PROJECT", "NAME", "CREATED BY", "CONDITION", "UPDATED", "CREATED", "PROVIDER"}

	preGardenerDirectives = gardenerDirectiveMap{
		GardenerTextAllFormat: func(_ *cloud_agent.Shoot) bool {
			return true
		},
	}

	postGardenerDirectives = gardenerDirectiveMap{
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
	filters Filters
	shoots  map[string]*cloud_agent.ShootList
}

func NewGardener(shoots map[string]*cloud_agent.ShootList, filters Filters) output.Formater {
	return &gardenerFormater{
		shoots:  shoots,
		filters: filters,
	}
}

func (f *gardenerFormater) YAML() interface{} {
	shootList := mergeShoots(f.shoots)
	shoots := f.filters.filter(shootList)

	return map[string]interface{}{
		"shoots": shoots.Shoots,
	}
}

func (f *gardenerFormater) JSON() interface{} {
	shootList := mergeShoots(f.shoots)
	shoots := f.filters.filter(shootList)

	return map[string]interface{}{
		"shoots": shoots.Shoots,
	}
}

func (f *gardenerFormater) Table() ([]string, [][]string) {
	rows := [][]string{}

	shootList := mergeShoots(f.shoots)
	shoots := f.filters.filter(shootList)

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
			shoot.LastTransitionTime.AsTime().Local().Format("2006-01-02 15:04:05"),
			shoot.CreationTimestamp.AsTime().Local().Format("2006-01-02 15:04:05"),
			gardenerProvider,
		})
	}

	return gardenerHeaders, rows
}

func (f *gardenerFormater) Text(outFormat, errFormat string) string {
	shoots := mergeShoots(f.shoots)
	directives := preGardenerDirectives.run(shoots, map[string]int{})

	shoots = f.filters.filter(shoots)

	directives = postGardenerDirectives.run(shoots, directives)

	str := outFormat
	for key, val := range directives {
		str = strings.ReplaceAll(str, key, strconv.Itoa(val))
	}

	return str
}

func mergeShoots(m map[string]*cloud_agent.ShootList) *cloud_agent.ShootList {
	l := &cloud_agent.ShootList{}
	for key := range m {
		l.Shoots = append(l.Shoots, m[key].Shoots...)
	}

	return l
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
	CreatedBy     string
	Project       string
	Condition     string
	LabelSelector string
	UpdatedAfter  time.Time
	UpdatedBefore time.Time
	CreatedAfter  time.Time
	CreatedBefore time.Time
}

func (f *Filters) filter(s *cloud_agent.ShootList) *cloud_agent.ShootList {
	if f.CreatedBy != "" {
		s = filterByCreatedBy(s, f.CreatedBy)
	}
	if f.Project != "" {
		s = project(s, f.Project)
	}
	if f.Condition != "" {
		s = filterByCondition(s, f.Condition)
	}
	if f.LabelSelector != "" {
		s = filterByLabelSelector(s, f.LabelSelector)
	}
	if !f.UpdatedAfter.IsZero() {
		s = filterByUpdatedAfter(s, f.UpdatedAfter)
	}
	if !f.UpdatedBefore.IsZero() {
		s = filterByUpdatedBefore(s, f.UpdatedBefore)
	}
	if !f.CreatedAfter.IsZero() {
		s = filterByCreatedAfter(s, f.CreatedAfter)
	}
	if !f.CreatedBefore.IsZero() {
		s = filterByCreatedBefore(s, f.CreatedBefore)
	}
	return s
}

func filterByCreatedAfter(s *cloud_agent.ShootList, v time.Time) *cloud_agent.ShootList {
	return filterBy(s, func(shoot *cloud_agent.Shoot) bool {
		return shoot.CreationTimestamp.AsTime().After(v)
	})
}

func filterByCreatedBefore(s *cloud_agent.ShootList, v time.Time) *cloud_agent.ShootList {
	return filterBy(s, func(shoot *cloud_agent.Shoot) bool {
		return shoot.CreationTimestamp.AsTime().Before(v)
	})
}

func filterByUpdatedAfter(s *cloud_agent.ShootList, v time.Time) *cloud_agent.ShootList {
	return filterBy(s, func(shoot *cloud_agent.Shoot) bool {
		return shoot.LastTransitionTime.AsTime().After(v)
	})
}

func filterByUpdatedBefore(s *cloud_agent.ShootList, v time.Time) *cloud_agent.ShootList {
	return filterBy(s, func(shoot *cloud_agent.Shoot) bool {
		return shoot.LastTransitionTime.AsTime().Before(v)
	})
}

func filterByLabelSelector(s *cloud_agent.ShootList, v string) *cloud_agent.ShootList {
	return filterBy(s, func(shoot *cloud_agent.Shoot) bool {
		selector, err := labels.Parse(v)
		if err != nil {
			return false
		}

		return selector.Matches(labels.Set(shoot.Labels))
	})
}

func filterByCondition(s *cloud_agent.ShootList, v string) *cloud_agent.ShootList {
	return filterBy(s, func(shoot *cloud_agent.Shoot) bool {
		return shoot.Condition.String() == v
	})
}

func project(s *cloud_agent.ShootList, v string) *cloud_agent.ShootList {
	return filterBy(s, func(shoot *cloud_agent.Shoot) bool {
		return shoot.Namespace == v
	})
}

func filterByCreatedBy(s *cloud_agent.ShootList, v string) *cloud_agent.ShootList {
	return filterBy(s, func(shoot *cloud_agent.Shoot) bool {
		return shoot.Annotations[createdByLabel] == v
	})
}

type statement func(s *cloud_agent.Shoot) bool

func filterBy(list *cloud_agent.ShootList, state statement) *cloud_agent.ShootList {
	l := &cloud_agent.ShootList{}

	for i := range list.Shoots {
		s := list.Shoots[i]
		if s == nil {
			continue
		}

		if state(s) {
			l.Shoots = append(l.Shoots, s)
		}
	}

	return l
}
