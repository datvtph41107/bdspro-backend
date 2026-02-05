// package contract

// import "context"

// type Event interface {
// 	Name() string
// 	Version() int
// 	Async() bool
// }

//	type OptionalValidator interface {
//		Validate(ctx context.Context) error
//	}
package contract

import "context"

type EventBus interface {
	Dispatch(ctx context.Context, payload interface{}) error
}
