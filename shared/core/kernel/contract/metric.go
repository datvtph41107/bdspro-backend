package contract

type Metric interface {
	Inc(name string, tags map[string]string)
	Observe(name string, value float64, tags map[string]string)
}
