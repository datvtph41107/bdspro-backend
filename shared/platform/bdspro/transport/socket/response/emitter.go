package socket_resp

import (
	"context"
	"time"

	corectx "bdspro/core/context"
	apperr "bdspro/core/pkg/errors"
)

type Connection interface {
	WriteJSON(v any) error
}

type Emitter struct{}

func NewEmitter() *Emitter {
	return &Emitter{}
}

func (e *Emitter) Emit(ctx context.Context, conn Connection, typ string, payload any) error {
	return conn.WriteJSON(Event{
		Type:      typ,
		Payload:   payload,
		RequestID: corectx.RequestID(ctx),
		Timestamp: time.Now().UTC(),
	})
}

func (e *Emitter) EmitError(ctx context.Context, conn Connection, err error) error {
	ae, ok := apperr.As(err)

	evt := Event{
		Type:      "error",
		RequestID: corectx.RequestID(ctx),
		Timestamp: time.Now().UTC(),
	}

	if !ok {
		evt.Error = &ErrorEvent{
			Code:    "INTERNAL",
			Message: "internal server error",
		}
		return conn.WriteJSON(evt)
	}

	evt.Error = &ErrorEvent{
		Code:    string(ae.Code()),
		Message: ae.Message(),
		Fields:  toFields(ae.Fields()),
	}

	return conn.WriteJSON(evt)
}

func toFields(in []apperr.FieldError) []FieldError {
	out := make([]FieldError, len(in))
	for i, f := range in {
		out[i] = FieldError{
			Field:   f.Field,
			Message: f.Message,
			Code:    f.Code,
		}
	}
	return out
}
