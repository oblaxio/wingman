package swarm

import (
	"bytes"

	"github.com/oblaxio/wingman/pkg/config"
	"github.com/oblaxio/wingman/pkg/print"
	"github.com/oblaxio/wingman/pkg/service"
)

type ServiceSwarm struct {
	group string
	swarm []*service.Service
}

func NewServiceSwarm(group ...string) *ServiceSwarm {
	sw := &ServiceSwarm{}
	if len(group) > 0 {
		sw.group = group[0]
	}
	return sw
}

func (sw *ServiceSwarm) Append(svc *service.Service) {
	sw.swarm = append(sw.swarm, svc)
}

func (sw *ServiceSwarm) RunServices() error {
	var stdOut, stdErr bytes.Buffer
	for serviceName := range config.Get().Services {
		go func(serviceName string) {
			if sw.group != "" && !contains(config.Get().ServiceGroups[sw.group], serviceName) {
				return
			}
			s, err := service.NewService(serviceName, ".")
			if err != nil {
				panic(err)
			}
			s.StdOut = &stdOut
			s.StdErr = &stdErr
			print.Info("calculating " + serviceName + " deepndencies")
			go s.GetDependencies()
			if err := s.Build(); err != nil {
				print.SvcErr(s.Executable, "\n"+s.StdErr.String())
				panic(err)
			}
			if err := s.Start(); err != nil {
				panic(err)
			}
			print.Info(s.Executable + " service started")
			sw.Append(s)
		}(serviceName)
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

func contains[T comparable](s []T, e T) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
