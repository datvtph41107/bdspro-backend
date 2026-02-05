package job

import (
	"context"
	"core/kernel/container"
	"core/kernel/db"
	"core/kernel/event"
	"encoding/json"
	"fmt"

	"gorm.io/gorm"
)

type SendBalancePayload struct {
	UserID     uint64   `json:"user_id"`
	Currencies []string `json:"currencies"`
	Store      string   `json:"store"`
}

func HandleSendBalanceJob(c *container.Container, payload []byte) error {
	var p SendBalancePayload
	if err := json.Unmarshal(payload, &p); err != nil {
		return err
	}

	fmt.Printf("SendBalance job invoked for user %d store %s\n", p.UserID, p.Store)

	dbm := c.Resolve("db.manager").(*db.Manager)
	gdb := dbm.Default.(*gorm.DB)

	return db.WithTransaction(context.Background(), gdb, func(tx *gorm.DB) error {
		balances := map[string]interface{}{
			"dummy": "value",
		}

		return event.SaveOutbox(tx, "BalanceUpdated", map[string]interface{}{
			"user_id":  p.UserID,
			"balances": balances,
			"store":    p.Store,
		})
	})
}
