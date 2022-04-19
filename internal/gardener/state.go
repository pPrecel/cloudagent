package gardener

import (
	"sync"

	"github.com/gardener/gardener/pkg/apis/core/v1beta1"
)

var _ StateGetter = &LastState{}
var _ StateSetter = &LastState{}

type StateGetter interface {
	Get() *v1beta1.ShootList
}

type StateSetter interface {
	Set(*v1beta1.ShootList)
}

type LastState struct {
	v *v1beta1.ShootList
	m sync.Mutex
}

func (s *LastState) Set(shoots *v1beta1.ShootList) {
	s.m.Lock()
	s.v = shoots
	s.m.Unlock()
}

func (s *LastState) Get() *v1beta1.ShootList {
	s.m.Lock()
	defer s.m.Unlock()

	return s.v
}
