package swarm

import (
	"bytes"

	"github.com/oblaxio/wingman/pkg/config"
	"github.com/oblaxio/wingman/pkg/print"
	"github.com/oblaxio/wingman/pkg/service"
)

type ServiceSwarmOption func(*ServiceSwarm) error

// type ServiceSwarm []*service.Service

type ServiceSwarm struct {
	config *config.Config
	swarm  []*service.Service
}

func NewServiceSwarm(options ...ServiceSwarmOption) (*ServiceSwarm, error) {
	s := &ServiceSwarm{}
	for _, o := range options {
		if err := o(s); err != nil {
			return nil, err
		}
	}
	return s, nil
}

func WithConfig(config *config.Config) ServiceSwarmOption {
	return func(s *ServiceSwarm) error {
		s.config = config
		return nil
	}
}

func (sw *ServiceSwarm) Append(svc *service.Service) {
	sw.swarm = append(sw.swarm, svc)
}

func (sw *ServiceSwarm) RunServices() error {
	var stdOut, stdErr bytes.Buffer
	for serviceName := range config.Get().Services {
		s, err := service.NewService(serviceName, ".")
		if err != nil {
			return err
		}
		s.StdOut = &stdOut
		s.StdErr = &stdErr
		s.GetDependencies()
		if err := s.Build(); err != nil {
			return err
		}
		if err := s.Start(); err != nil {
			return err
		}
		print.PrintInfo(s.Executable + " service started")
		sw.Append(s)
	}

	return nil
}

func (sw *ServiceSwarm) List() []*service.Service {
	return sw.swarm
}
