package handlers

import (
	"log"

	"github.com/oblaxio/wingman/pkg/config"
	"github.com/oblaxio/wingman/pkg/proxy"
	"github.com/oblaxio/wingman/pkg/sourcewatcher"
	"github.com/oblaxio/wingman/pkg/swarm"
	"github.com/spf13/cobra"
)

func StartHandler(configFile *string) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		// load config
		if err := config.NewConfigFromFile(*configFile); err != nil {
			log.Fatal(err)
		}
		group := ""
		if len(args) == 1 {
			group = args[0]
		}
		// create service swarm
		serviceSwarm := swarm.NewServiceSwarm(group)
		defer serviceSwarm.KillAll()
		// start services
		if err := serviceSwarm.RunServices(); err != nil {
			log.Fatal(err)
		}
		// setup source watcher
		watcher, err := sourcewatcher.NewWatcher(serviceSwarm)
		if err != nil {
			log.Fatal(err)
		}
		defer watcher.Stop()
		// check whether to start proxy
		if config.Get().Proxy.Enabled {
			// start proxy
			p, err := proxy.NewServer()
			if err != nil {
				log.Fatal(err)
			}
			go p.Serve()
		}
		// start source watcher
		watcher.Start()
	}
}
