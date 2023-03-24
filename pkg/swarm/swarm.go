package swarm

import (
	"github.com/oblaxio/wingman/pkg/service"
)

type ServiceSwarm []*service.Service

func NewServiceSwarm(length int) *ServiceSwarm {
	var s ServiceSwarm = make([]*service.Service, length)
	return &s
}

func (s *ServiceSwarm) Append(svc *service.Service) {
	*s = append(*s, svc)
}
