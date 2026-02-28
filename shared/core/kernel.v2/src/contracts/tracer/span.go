package tracer

import "time"

type SpanContext struct {
	TraceID    string
	SpanID     string
	TraceFlags byte
	TraceState string
	IsRemote   bool
}

type Code int

const (
	CodeUnset Code = iota
	CodeOk
	CodeError
)

type Attribute struct {
	Key   string
	Value interface{}
}

type SpanOption func(*SpanConfig)

type SpanConfig struct {
	Attributes []Attribute
	Links      []Link
	StartTime  time.Time
}

type Link struct {
	SpanContext SpanContext
	Attributes  []Attribute
}
