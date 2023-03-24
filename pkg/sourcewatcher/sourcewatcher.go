package sourcewatcher

import (
	"fmt"
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/oblaxio/wingman/pkg/config"
	"github.com/oblaxio/wingman/pkg/swarm"
)

type SourceWatcherOption func(*SourceWatcher) error

type SourceWatcher struct {
	config  *config.Config
	watcher *fsnotify.Watcher
	swarm   *swarm.ServiceSwarm
}

func NewWatcher(options ...SourceWatcherOption) (*SourceWatcher, error) {
	w := &SourceWatcher{}
	for _, option := range options {
		if err := option(w); err != nil {
			return nil, err
		}
	}
	var err error
	w.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	for _, v := range config.Get().WatchDir {
		if err := addWatcher(".", w.watcher, v); err != nil {
			return nil, err
		}
	}
	return w, nil
}

func WithConfig(config *config.Config) SourceWatcherOption {
	return func(w *SourceWatcher) error {
		w.config = config
		return nil
	}
}

func WithSwarm(swarm *swarm.ServiceSwarm) SourceWatcherOption {
	return func(w *SourceWatcher) error {
		w.swarm = swarm
		return nil
	}
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
					// log.Println("file changed:", event.Name)
					for _, s := range w.swarm.List() {
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

func (w *SourceWatcher) Stop() {
	w.watcher.Close()
}
