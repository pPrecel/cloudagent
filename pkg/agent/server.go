package agent

import (
	"context"
	"math"

	cloud_agent "github.com/pPrecel/cloudagent/pkg/agent/proto"
	"github.com/pPrecel/cloudagent/pkg/cache"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/rand"
)

var _ cloud_agent.AgentServer = &server{}

type ResourceGetter interface {
	GetGardenerCache() cache.GardenerCache
	GetGeneralError() error
}

type ServerOption struct {
	ResourceGetter ResourceGetter
	Logger         *logrus.Entry
}

type server struct {
	cloud_agent.UnimplementedAgentServer
	resourceGetter ResourceGetter
	logger         *logrus.Entry
}

func NewServer(opts *ServerOption) cloud_agent.AgentServer {
	return &server{
		logger:         opts.Logger,
		resourceGetter: opts.ResourceGetter,
	}
}

func (s *server) GardenerShoots(ctx context.Context, _ *cloud_agent.Empty) (*cloud_agent.GardenerResponse, error) {
	logger := s.logger.WithField("request-id", rand.Intn(math.MaxInt))
	logger.Debug("handling request")

	if s.resourceGetter == nil {
		errMessage := "can't get latest shoots data"
		logger.Debug(errMessage)
		return nil, errors.New(errMessage)
	}

	resp := toGardenerResponse(s.resourceGetter)

	logger.Debugf("returning list of '%d' elements, with err: '%s'", len(resp.ShootList), resp.GeneralError)

	return resp, nil
}
