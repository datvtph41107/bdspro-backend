package memoryusage

import (
	"context"
	"errors"
	"tqd/internal/usecase/usage"
)

type readCommandState struct {
	command       usage.ReadCommand
	accepted      bool
	reservationID string
}

// Map dùng struct key để không có va chạm do ghép chuỗi input.
type readCommandKey struct{ subjectType, subjectID, operation, commandKey string }

func commandKey(c usage.ReadCommand) readCommandKey {
	return readCommandKey{string(c.Subject.Type), c.Subject.ID, string(c.Operation), c.CommandKey}
}
func (s *Store) BindReadCommand(ctx context.Context, c usage.ReadCommand) (bool, error) {
	if err := ctx.Err(); err != nil {
		return false, err
	}
	if !c.IsValid() {
		return false, errors.New("invalid read command")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	key := commandKey(c)
	if row, found := s.commands[key]; found {
		if row.command.RequestHash != c.RequestHash {
			return false, usage.ErrReadCommandConflict
		}
		return row.accepted, nil
	}
	s.commands[key] = readCommandState{command: c}
	return false, nil
}
func (s *Store) AcceptReadCommand(ctx context.Context, c usage.ReadCommand, event *usage.Event) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	if !c.IsValid() || (event != nil && (!event.IsValid() || !c.Matches(*event))) {
		return "", errors.New("invalid read acceptance")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	key := commandKey(c)
	row, found := s.commands[key]
	if !found {
		return "", errors.New("read command is not bound")
	}
	if row.command.RequestHash != c.RequestHash {
		return "", usage.ErrReadCommandConflict
	}
	if row.accepted {
		return row.reservationID, nil
	}
	if event != nil {
		if saved, found := s.events[event.UsageKey]; found {
			row.reservationID = saved.ReservationID
		} else {
			s.events[event.UsageKey] = *event
			row.reservationID = event.ReservationID
		}
	}
	row.accepted = true
	s.commands[key] = row
	return row.reservationID, nil
}
