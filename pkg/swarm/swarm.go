package swarm

import (
	"bytes"
	"fmt"
	"maps"
	"slices"
	"strings"
	"sync"

	"github.com/oblaxio/wingman/pkg/config"
	"github.com/oblaxio/wingman/pkg/print"
	"github.com/oblaxio/wingman/pkg/service"
)

const groupPrefix = "@"

type ServiceSwarm struct {
	swarm        []*service.Service
	servicesList map[string]bool
	loadedGroups map[string]bool
	startGroups  [][]string
}

func NewServiceSwarm(group ...string) *ServiceSwarm {
	sw := &ServiceSwarm{
		servicesList: make(map[string]bool),
		loadedGroups: make(map[string]bool),
		startGroups:  [][]string{},
	}
	if len(group) > 0 {
		sw.getGroupedServices(group[0])
		return sw
	}
	for serviceName := range config.Get().Services {
		sw.servicesList[serviceName] = true
	}
	return sw
}

func (sw *ServiceSwarm) getGroupedServices(group string) map[string]bool {
	sw.loadedGroups[group] = true
	for _, svc := range config.Get().ServiceGroups[group] {
		if strings.HasPrefix(svc, groupPrefix) {
			groupName := svc[1:]
			if _, loaded := sw.loadedGroups[groupName]; loaded {
				continue
			} else {
				nestedGroupedServices := sw.getGroupedServices(groupName)
				maps.Copy(sw.servicesList, nestedGroupedServices)
			}
		} else {
			sw.servicesList[svc] = true
		}
	}
	return sw.servicesList
}

func (sw *ServiceSwarm) Append(svc *service.Service) {
	sw.swarm = append(sw.swarm, svc)
}

func (sw *ServiceSwarm) RunServices() error {
	originalList := make([]string, 0)
	for serviceName, _ := range sw.servicesList {
		originalList = append(originalList, serviceName)
	}
	if err := sw.getDependencyOrderedList(originalList, []string{}); err != nil {
		return err
	}
	var stdOut, stdErr bytes.Buffer
	var wg sync.WaitGroup
	for _, group := range sw.startGroups {
		wg.Add(len(group))
		for _, serviceName := range group {
			go func(serviceName string) {
				s, err := service.NewService(serviceName, ".")
				if err != nil {
					panic(err)
				}
				s.StdOut = &stdOut
				s.StdErr = &stdErr
				print.Info("calculating " + serviceName + " deepndencies")
				go s.GetDependencies()
				if s.Entrypoint != "" {
					if err := s.Build(); err != nil {
						print.SvcErr(s.Executable, "\n"+s.StdErr.String())
						panic(err)
					}
					if err := s.Start(); err != nil {
						panic(err)
					}
				}
				print.Info(s.Executable + " service started")
				wg.Done()
				sw.Append(s)
			}(serviceName)
		}
		wg.Wait()
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

func (sw *ServiceSwarm) getDependencyOrderedList(originalList []string, orderedList []string) error {
	countOriginal := len(originalList)
	countOrdered := len(orderedList)
	groupServices := []string{}
	for k := 0; k < len(originalList); k++ {
		serviceName := originalList[k]
		if config.Get().Services[serviceName].DependsOn != nil {
			clearedDeps := 0
			for _, dep := range config.Get().Services[serviceName].DependsOn {
				if !slices.Contains(orderedList, dep) {
					clearedDeps++
				}
			}
			if clearedDeps == 0 {
				orderedList = append(orderedList, serviceName)
				groupServices = append(groupServices, serviceName)
				originalList = append(originalList[:k], originalList[k+1:]...)
			}
		} else {
			orderedList = append(orderedList, serviceName)
			groupServices = append(groupServices, serviceName)
			originalList = append(originalList[:k], originalList[k+1:]...)
		}
	}
	sw.startGroups = append(sw.startGroups, groupServices)
	if countOriginal == len(originalList) && countOrdered == len(orderedList) {
		return fmt.Errorf("circular dependency detected or bad dependency configuration")
	}
	if len(orderedList) != len(sw.servicesList) {
		if err := sw.getDependencyOrderedList(originalList, orderedList); err != nil {
			return err
		}
	}
	return nil
}
