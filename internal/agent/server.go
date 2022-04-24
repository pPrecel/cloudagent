package agent

import (
	"context"
	"errors"

	"github.com/gardener/gardener/pkg/apis/core/v1beta1"
	cloud_agent "github.com/pPrecel/cloud-agent/internal/agent/proto"
	"github.com/sirupsen/logrus"
)

var _ cloud_agent.AgentServer = &server{}

//go:generate mockery --name=StateGetter --output=automock --outpkg=automock
type StateGetter interface {
	Get() *v1beta1.ShootList
}

type ServerOption struct {
	Getter StateGetter
	Logger *logrus.Logger
}

type server struct {
	cloud_agent.UnimplementedAgentServer
	getter StateGetter
	logger *logrus.Logger
}

func NewServer(opts *ServerOption) cloud_agent.AgentServer {
	return &server{
		logger: opts.Logger,
		getter: opts.Getter,
	}
}

func (s *server) GardenerShoots(ctx context.Context, _ *cloud_agent.Empty) (*cloud_agent.ShootList, error) {
	s.logger.Debug("handling request")

	state := s.getter.Get()
	if state == nil {
		return nil, errors.New("can't get latest shoots list")
	}

	list := &cloud_agent.ShootList{}
	for i := range state.Items {
		item := state.Items[i]

		cond := cloud_agent.Condition_HEALTHY
		if item.Status.IsHibernated {
			cond = cloud_agent.Condition_HIBERNATED
		} else if isConditionUnknown(item) {
			cond = cloud_agent.Condition_UNKNOWN
		}

		list.Shoots = append(list.Shoots, &cloud_agent.Shoot{
			Name:        item.Name,
			Namespace:   item.Namespace,
			Labels:      item.Labels,
			Annotations: item.Annotations,
			Condition:   cond,
		})
	}

	s.logger.Debugf("returning list of %v items", len(list.Shoots))
	return list, nil
}

func isConditionUnknown(shoot v1beta1.Shoot) bool {
	for i := range shoot.Status.Conditions {
		if shoot.Status.Conditions[i].Status != v1beta1.ConditionTrue {
			return true
		}
	}
	return false
}
