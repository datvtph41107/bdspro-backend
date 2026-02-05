package provider

import "core/kernel/container"

type ServiceProvider interface {
	Register(c *container.Container) error
	Boot(c *container.Container) error
}
