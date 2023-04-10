package handlers

import (
	"fmt"
	"log"

	"github.com/oblaxio/wingman/pkg/config"
	"github.com/oblaxio/wingman/pkg/proxy"
	"github.com/oblaxio/wingman/pkg/sourcewatcher"
	"github.com/oblaxio/wingman/pkg/swarm"
	"github.com/spf13/cobra"
)

func StartHandler(cmd *cobra.Command, args []string) {
	fmt.Println("Start handler")
	// get the right config file
	configFile := config.DefaultConfigFile
	if len(args) == 1 {
		configFile = args[0]
	}
	// read configuration
	err := config.NewConfigFromFile(configFile)
	if err != nil {
		log.Fatal(err)
	}
	// create service swarm
	swm := swarm.NewServiceSwarm()
	defer swm.KillAll()
	// start services
	if err := swm.RunServices(); err != nil {
		log.Fatal(err)
	}
	// setup source watcher
	sw, err := sourcewatcher.NewWatcher(swm)
	if err != nil {
		log.Fatal(err)
	}
	defer sw.Stop()
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
	sw.Start()
}
