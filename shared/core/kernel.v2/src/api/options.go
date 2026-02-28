package api

func WithName(name string) Option {
	return func(k *Kernel) { k.name = name }
}

func WithVersion(version string) Option {
	return func(k *Kernel) { k.version = version }
}

func WithEnvironment(env string) Option {
	return func(k *Kernel) { k.env = env }
}
