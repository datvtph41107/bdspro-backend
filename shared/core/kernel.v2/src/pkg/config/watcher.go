package config

import (
	"log"

	"github.com/fsnotify/fsnotify"
)

func (c *Config[T]) watch() {
	c.v.WatchConfig()

	c.v.OnConfigChange(func(e fsnotify.Event) {
		old := c.Get()

		if err := c.reload(); err != nil {
			log.Println("config reload failed:", err)
			return
		}

		newCfg := c.Get()

		c.mu.RLock()
		defer c.mu.RUnlock()

		for _, h := range c.hooks {
			go h(old, newCfg)
		}
	})
}
