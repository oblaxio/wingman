package swarm

import (
	"bytes"

	"github.com/oblaxio/wingman/pkg/config"
	"github.com/oblaxio/wingman/pkg/print"
	"github.com/oblaxio/wingman/pkg/service"
)

type ServiceSwarm struct {
	swarm []*service.Service
}

func NewServiceSwarm() *ServiceSwarm {
	return &ServiceSwarm{}
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
		print.Info(s.Executable + " service started")
		sw.Append(s)
	}

	return nil
}

func (sw *ServiceSwarm) List() []*service.Service {
	return sw.swarm
}

func (sw *ServiceSwarm) KillAll() error {
	for _, s := range sw.swarm {
		if err := s.Stop(); err != nil {
			return err
		}
	}
	return nil
}
