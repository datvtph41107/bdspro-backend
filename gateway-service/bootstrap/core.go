package bootstrap

import "core/kernel/config"

func Init() {
	_ = config.Load
}
