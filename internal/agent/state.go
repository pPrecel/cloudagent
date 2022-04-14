package agent

import (
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
}

func (s *LastState) Set(shoots *v1beta1.ShootList) {
	s.v = shoots
}

func (s *LastState) Get() *v1beta1.ShootList {
	return s.v
}
