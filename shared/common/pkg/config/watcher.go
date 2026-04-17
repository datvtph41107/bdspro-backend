package config

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

// WatchConfig watches a config file and calls onReload when changed.
// It validates the new config before applying.
func WatchConfig(ctx context.Context, configPath string, loader func() (interface{}, error), onReload func(interface{})) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	// Watch the directory (not the file) because some editors save to temp and rename.
	dir := filepath.Dir(configPath)
	if err := watcher.Add(dir); err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				watcher.Close()
				return
			case event := <-watcher.Events:
				if event.Op&fsnotify.Write != 0 || event.Op&fsnotify.Create != 0 {
					// Avoid duplicate events
					info, err := os.Stat(event.Name)
					if err != nil || info.IsDir() {
						continue
					}
					log.Printf("[config] file %s changed, reloading...", event.Name)
					newCfg, err := loader()
					if err != nil {
						log.Printf("[config] reload failed: %v", err)
						continue
					}
					onReload(newCfg)
				}
			}
		}
	}()
	return nil
}
