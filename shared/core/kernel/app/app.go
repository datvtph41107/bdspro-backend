package app

import (
	"core/kernel/container"
	"core/kernel/provider"
	"core/kernel/runtime"
	"log"
)

type Application struct {
	Container  *container.Container
	Runtime    *runtime.Runtime
	providers  []provider.ServiceProvider
	providersM map[string]struct{}
}

func NewApp() *Application {
	c := container.New()
	return &Application{
		Container:  c,
		Runtime:    runtime.NewRuntime(),
		providers:  []provider.ServiceProvider{},
		providersM: map[string]struct{}{},
	}
}

func (a *Application) RegisterProvider(p provider.ServiceProvider) {
	// avoid double register of same type
	name := container.TypeKey(p)
	if _, ok := a.providersM[name]; ok {
		return
	}
	a.providers = append(a.providers, p)
	a.providersM[name] = struct{}{}
}

func (a *Application) Boot() error {
	// Register phase
	for _, p := range a.providers {
		if err := p.Register(a.Container); err != nil {
			return err
		}
	}
	// Boot phase
	for _, p := range a.providers {
		if err := p.Boot(a.Container); err != nil {
			return err
		}
	}
	return nil
}

func (a *Application) Run() error {
	log.Println("app: booting providers")
	if err := a.Boot(); err != nil {
		return err
	}
	log.Println("app: starting runtime")
	return a.Runtime.Start(a.Container)
}

func (a *Application) Make(key string) interface{} {
	v, err := a.Container.Resolve(key)
	if err != nil {
		panic(err)
	}
	return v
}
