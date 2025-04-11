package sourcewatcher

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/oblaxio/wingman/pkg/config"
	"github.com/oblaxio/wingman/pkg/print"
	"github.com/oblaxio/wingman/pkg/service"
	"github.com/oblaxio/wingman/pkg/swarm"
)

type SourceWatcherOption func(*SourceWatcher) error

type SourceWatcher struct {
	watcher *fsnotify.Watcher
	swarm   *swarm.ServiceSwarm
}

func NewWatcher(swarm *swarm.ServiceSwarm) (*SourceWatcher, error) {
	w := &SourceWatcher{
		swarm: swarm,
	}
	var err error
	w.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	for _, v := range config.Get().Watchers.IncludeDirs {
		if err := addWatcher(".", w.watcher, v); err != nil {
			return nil, err
		}
	}
	return w, nil
}

func (w *SourceWatcher) Start() error {
	go func() {
		for {
			select {
			case event, ok := <-w.watcher.Events:
				if !ok {
					return
				}
				if event.Op == fsnotify.Create || event.Op == fsnotify.Write {
					print.Rebuild("file change: " + event.Name)
					for _, s := range w.swarm.List() {
						go func(svc *service.Service) {
							svc.Lock()
							if svc != nil {
								if svc.CheckDependency(event.Name) {
									if err := svc.Stop(); err != nil {
										fmt.Println(err)
									}
									if err := svc.Build(); err != nil {
										fmt.Println(err)
									}
									if err := svc.Start(); err != nil {
										fmt.Println(err)
									}
									print.Rebuild("rebuilt " + svc.Executable)
								}
							}
							svc.Unlock()
						}(s)
					}
				}
			case err, ok := <-w.watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()
	<-make(chan struct{})
	return nil
}

func (w *SourceWatcher) Watcher() *fsnotify.Watcher {
	return w.watcher
}

func addWatcher(prefix string, watcher *fsnotify.Watcher, dir string) error {
	fullDir := fmt.Sprintf("%s/%s", prefix, dir)
	// check whether this dir needs to be ignorred
	for _, d := range config.Get().Watchers.ExcludeDirs {
		if fullDir == d {
			return nil // skip dir watch
		}
	}
	items, err := os.ReadDir(fullDir)
	if err != nil {
		return err
	}
	for _, i := range items {
		path := fmt.Sprintf("%s/%s", fullDir, i.Name())
		// check whether the file is in the negative watch list
		if len(config.Get().Watchers.ExcludeFiles) > 0 && !i.IsDir() {
			for _, f := range config.Get().Watchers.ExcludeFiles {
				match, err := filepath.Match(fullDir+"/"+f, path)
				if err != nil {
					return err
				}
				if match {
					continue
				}
			}
		}
		// check whether the file is in the positive watch list
		if len(config.Get().Watchers.IncludeFiles) > 0 && !i.IsDir() {
			for _, f := range config.Get().Watchers.IncludeFiles {
				match, err := filepath.Match(fullDir+"/"+f, path)
				if err != nil {
					return err
				}
				if !match {
					continue
				}
			}
		}
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

func (w *SourceWatcher) Stop() {
	w.watcher.Close()
}
