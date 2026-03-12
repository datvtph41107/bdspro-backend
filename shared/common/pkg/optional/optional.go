package optional

type State uint8

const (
	StateAbsent  State = iota // 0: field không có trong request
	StateNull                 // 1: field có trong request với giá trị null
	StatePresent              // 2: field có trong request với giá trị
)

// Value là interface chung cho mọi kiểu optional
type Value interface {
	State() State
	IsAbsent() bool
	IsNull() bool
	IsPresent() bool
	Interface() interface{}
}
