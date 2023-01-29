package agent

import (
	"time"

	"github.com/gardener/gardener/pkg/apis/core/v1beta1"
	cloud_agent "github.com/pPrecel/cloudagent/pkg/agent/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func toGardenerResponse(serverCache ResourceGetter) *cloud_agent.GardenerResponse {
	err := ""
	if serverCache.GetGeneralError() != nil {
		err = serverCache.GetGeneralError().Error()
	}

	resources := serverCache.GetGardenerCache().Resources()
	shootList := map[string]*cloud_agent.ShootList{}
	for key := range resources {
		r := resources[key]
		shootList[key] = toShootList(r)
	}

	return &cloud_agent.GardenerResponse{
		ShootList:    shootList,
		GeneralError: err,
	}
}

func toShootList(resource RegisteredResource[*v1beta1.ShootList]) *cloud_agent.ShootList {
	r := resource.Get()

	err := ""
	if r.Error != nil {
		err = r.Error.Error()
	}

	list := &cloud_agent.ShootList{
		Error:  err,
		Time:   timestamppb.New(r.Time),
		Shoots: []*cloud_agent.Shoot{},
	}

	if r.Value == nil {
		return list
	}

	for i := range r.Value.Items {
		s := &r.Value.Items[i]
		list.Shoots = append(list.Shoots, toShoot(s))
	}

	return list
}

func toShoot(shoot *v1beta1.Shoot) *cloud_agent.Shoot {
	cond := cloud_agent.Condition_HEALTHY
	if shoot.Status.IsHibernated {
		cond = cloud_agent.Condition_HIBERNATED
	} else if isConditionUnknown(shoot) {
		cond = cloud_agent.Condition_UNKNOWN
	} else if len(shoot.Status.Conditions) == 0 {
		cond = cloud_agent.Condition_EMPTY
	}

	return &cloud_agent.Shoot{
		Name:               shoot.Name,
		Namespace:          shoot.Namespace,
		Labels:             shoot.Labels,
		Annotations:        shoot.Annotations,
		Condition:          cond,
		LastTransitionTime: lastConditionUpdate(shoot),
		CreationTimestamp:  timestamppb.New(shoot.ObjectMeta.CreationTimestamp.Time),
	}
}

func lastConditionUpdate(shoot *v1beta1.Shoot) *timestamppb.Timestamp {
	last := time.Time{}
	for i := range shoot.Status.Conditions {
		cond := shoot.Status.Conditions[i]

		if cond.LastTransitionTime.Time.After(last) {
			last = cond.LastTransitionTime.Time
		}
	}

	return timestamppb.New(last)
}

func isConditionUnknown(shoot *v1beta1.Shoot) bool {
	for i := range shoot.Status.Conditions {
		if shoot.Status.Conditions[i].Status != v1beta1.ConditionTrue {
			return true
		}
	}
	return false
}
