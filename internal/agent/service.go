package agent

import (
	"context"

	"github.com/gardener/gardener/pkg/apis/core/v1beta1"
	gardener_agent "github.com/pPrecel/gardener-agent/internal/agent/proto"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	Address = "/tmp/gardener-agent.sock"
	Network = "unix"
)

var _ gardener_agent.AgentServer = &server{}

type ServerOption struct {
	Getter StateGetter
	Logger *logrus.Logger
}

type server struct {
	gardener_agent.UnimplementedAgentServer

	logger *logrus.Logger
	getter StateGetter
}

func NewServer(opts *ServerOption) gardener_agent.AgentServer {
	return &server{
		logger: opts.Logger,
		getter: opts.Getter,
	}
}

func (s *server) Shoots(ctx context.Context, _ *gardener_agent.Empty) (*gardener_agent.ShootList, error) {
	s.logger.Debug("handling reque")

	state := s.getter.Get()
	if state == nil {
		return nil, errors.New("can't get latest shoots list")
	}

	list := &gardener_agent.ShootList{}
	for i := range state.Items {
		item := state.Items[i]

		cond := gardener_agent.Condition_HEALTHY
		if item.Status.IsHibernated {
			cond = gardener_agent.Condition_HIBERNATED
		} else if isConditionUnknown(item) {
			cond = gardener_agent.Condition_UNKNOWN
		}

		list.Shoots = append(list.Shoots, &gardener_agent.Shoot{
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
