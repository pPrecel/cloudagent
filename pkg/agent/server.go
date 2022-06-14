package agent

import (
	"context"
	"errors"

	"github.com/gardener/gardener/pkg/apis/core/v1beta1"
	cloud_agent "github.com/pPrecel/cloudagent/pkg/agent/proto"
	"github.com/sirupsen/logrus"
)

var _ cloud_agent.AgentServer = &server{}

type ServerOption struct {
	GardenerCache Cache[*v1beta1.ShootList]
	Logger        *logrus.Logger
}

type server struct {
	cloud_agent.UnimplementedAgentServer
	gardenerCache Cache[*v1beta1.ShootList]
	logger        *logrus.Logger
}

func NewServer(opts *ServerOption) cloud_agent.AgentServer {
	return &server{
		logger:        opts.Logger,
		gardenerCache: opts.GardenerCache,
	}
}

func (s *server) GardenerShoots(ctx context.Context, _ *cloud_agent.Empty) (*cloud_agent.ShootList, error) {
	s.logger.Debug("handling request")

	if s.gardenerCache == nil {
		return nil, errors.New("can't get latest shoots list")
	}

	v1beta1List := &v1beta1.ShootList{}
	r := s.gardenerCache.Resources()
	for key := range r {
		if r[key] != nil && r[key].Get() != nil {
			v1beta1List.Items = append(v1beta1List.Items, r[key].Get().Items...)
		}
	}

	agentList := &cloud_agent.ShootList{}
	for i := range v1beta1List.Items {
		item := v1beta1List.Items[i]

		cond := cloud_agent.Condition_HEALTHY
		if item.Status.IsHibernated {
			cond = cloud_agent.Condition_HIBERNATED
		} else if isConditionUnknown(item) {
			cond = cloud_agent.Condition_UNKNOWN
		} else if len(item.Status.Conditions) == 0 {
			cond = cloud_agent.Condition_EMPTY
		}

		agentList.Shoots = append(agentList.Shoots, &cloud_agent.Shoot{
			Name:        item.Name,
			Namespace:   item.Namespace,
			Labels:      item.Labels,
			Annotations: item.Annotations,
			Condition:   cond,
		})
	}

	s.logger.Debugf("returning list of %v items", len(agentList.Shoots))
	return agentList, nil
}

func isConditionUnknown(shoot v1beta1.Shoot) bool {
	for i := range shoot.Status.Conditions {
		if shoot.Status.Conditions[i].Status != v1beta1.ConditionTrue {
			return true
		}
	}
	return false
}
