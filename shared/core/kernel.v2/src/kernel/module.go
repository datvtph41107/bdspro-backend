package kernel

import "fmt"

type Module interface {
	ID() ComponentID
	Build() (Component, error)
}

func RegisterModules(app *App, modules ...Module) error {
	for _, module := range modules {
		if module == nil {
			continue
		}
		comp, err := module.Build()
		if err != nil {
			return fmt.Errorf("build module %s: %w", module.ID(), err)
		}
		if err := app.Register(comp); err != nil {
			return fmt.Errorf("register module %s: %w", module.ID(), err)
		}
	}
	return nil
}
