// shared/pkg/config/watcher.go (nếu cần hot-reload)
package config

import (
	"context"
	"log"
	"sync/atomic"
	"unsafe"

	"github.com/fsnotify/fsnotify"
)

type AtomicConfig struct {
	ptr unsafe.Pointer
}

func NewAtomicConfig(initial interface{}) *AtomicConfig {
	ac := &AtomicConfig{}
	ac.Store(initial)
	return ac
}

func (ac *AtomicConfig) Load() interface{} {
	return atomic.LoadPointer(&ac.ptr)
}

func (ac *AtomicConfig) Store(cfg interface{}) {
	atomic.StorePointer(&ac.ptr, unsafe.Pointer(&cfg))
}

func WatchFile(ctx context.Context, path string, onReload func()) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	if err := watcher.Add(path); err != nil {
		return err
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
				watcher.Close()
				return
			case event := <-watcher.Events:
				if event.Op&fsnotify.Write != 0 {
					log.Println("[config] file changed, reloading...")
					onReload()
				}
			}
		}
	}()
	return nil
}
