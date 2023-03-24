package handlers

import (
	"fmt"
	"log"

	"github.com/oblaxio/wingman/pkg/config"
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
	swm, err := swarm.NewServiceSwarm(
		swarm.WithConfig(config.Get()),
	)
	if err != nil {
		log.Fatal(err)
	}
	// start services
	if err := swm.RunServices(); err != nil {
		log.Fatal(err)
	}
	// setup source watcher
	sw, err := sourcewatcher.NewWatcher(
		sourcewatcher.WithConfig(config.Get()),
		sourcewatcher.WithSwarm(swm),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer sw.Stop()
	// start source watcher
	sw.Start()
}
