package handlers

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/oblaxio/wingman/pkg/config"
	"github.com/oblaxio/wingman/pkg/print"
	"github.com/oblaxio/wingman/pkg/service"
	"github.com/oblaxio/wingman/pkg/swarm"
	"github.com/spf13/cobra"
)

func StartHandler(cmd *cobra.Command, args []string) {
	fmt.Println("Start handler")
	var stdOut, stdErr bytes.Buffer
	configFile := config.DefaultConfigFile
	if len(args) == 1 {
		configFile = args[0]
	}
	err := config.Read(configFile)
	if err != nil {
		log.Fatal(err)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	swarm := swarm.NewServiceSwarm(len(config.Get().Services))
	for _, v := range config.Get().WatchDir {
		if err := addWatcher(".", watcher, v); err != nil {
			fmt.Println(err)
		}
	}

	for serviceName := range config.Get().Services {
		s, err := service.NewService(serviceName, ".")
		if err != nil {
			fmt.Println(err)
		}
		s.StdOut = &stdOut
		s.StdErr = &stdErr
		s.GetDependencies()
		if err := s.Build(); err != nil {
			fmt.Println(err)
		}
		if err := s.Start(); err != nil {
			fmt.Println(err)
		}
		print.PrintInfo(s.Output + " service started")
		swarm.Append(s)
	}

	// Start listening for events.
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op == fsnotify.Create || event.Op == fsnotify.Write {
					// log.Println("file changed:", event.Name)
					for _, s := range *swarm {
						if s != nil {
							if s.CheckDependency(event.Name) {
								if err := s.Stop(); err != nil {
									fmt.Println(err)
								}
								if err := s.Build(); err != nil {
									fmt.Println(err)
								}
								if err := s.Start(); err != nil {
									fmt.Println(err)
								}
							}
						}
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	<-make(chan struct{})
}

func addWatcher(prefix string, watcher *fsnotify.Watcher, dir string) error {
	fullDir := fmt.Sprintf("%s/%s", prefix, dir)
	items, err := os.ReadDir(fullDir)
	if err != nil {
		return err
	}
	for _, i := range items {
		path := fmt.Sprintf("%s/%s", fullDir, i.Name())
		err = watcher.Add(path)
		if err != nil {
			return err
		}
		if i.IsDir() {
			if err := addWatcher(fullDir, watcher, i.Name()); err != nil {
				return err
			}
		}
	}
	return nil
}
