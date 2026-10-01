package usage

import (
	"common/operation"
	"common/request"
	"encoding/hex"
	"errors"
	"userv4/modules/tqd/access"
)

var ErrReadCommandConflict = errors.New("read command key was already bound to another payload")
var ErrLegacyReadCommand = errors.New("legacy read command has no payload evidence; start a new action")

// ReadCommand là identity bền vững của ý định đọc. Không chứa meter/chu kỳ:
// thay gói không được biến retry thành giao dịch tiêu thụ mới.
// RequestHash do backend tính từ protobuf request, không tin hash từ client.
type ReadCommand struct {
	Subject     access.Subject
	Operation   operation.Code
	CommandKey  string
	RequestHash string
}

func (c ReadCommand) IsValid() bool {
	hash, err := hex.DecodeString(c.RequestHash)
	return c.Subject.IsValid() && c.Operation.IsValid() && request.IsValidIdempotencyKey(c.CommandKey) && err == nil && len(hash) == 32
}

func (c ReadCommand) Matches(event Event) bool {
	return event.Subject == c.Subject && event.Operation == c.Operation && event.CommandKey == c.CommandKey
}
