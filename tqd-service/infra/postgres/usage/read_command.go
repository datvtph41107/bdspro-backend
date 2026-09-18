package postgresusage

import (
	"context"
	"errors"
	"time"
	"tqd/internal/usecase/usage"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type readCommandModel struct {
	SubjectType      string `gorm:"primaryKey"`
	SubjectID        string `gorm:"primaryKey"`
	Operation        string `gorm:"primaryKey"`
	CommandKey       string `gorm:"primaryKey"`
	RequestHash      *string
	AcceptedUsageKey *string
	CreatedAt        time.Time
	AcceptedAt       *time.Time
}

func (readCommandModel) TableName() string { return "quota_read_commands" }
func readCommandWhere(db *gorm.DB, c usage.ReadCommand) *gorm.DB {
	return db.Where("subject_type = ? AND subject_id = ? AND operation = ? AND command_key = ?", string(c.Subject.Type), c.Subject.ID, string(c.Operation), c.CommandKey)
}
func verifyReadCommand(row readCommandModel, c usage.ReadCommand) error {
	if row.RequestHash == nil {
		return usage.ErrLegacyReadCommand
	}
	if *row.RequestHash != c.RequestHash {
		return usage.ErrReadCommandConflict
	}
	return nil
}

// BindReadCommand phân xử key trước khi đọc dữ liệu hoặc reserve. Insert
// ON CONFLICT cùng unique key chặn hai payload khác nhau chiếm cùng command.
func (s *Store) BindReadCommand(ctx context.Context, c usage.ReadCommand) (bool, error) {
	if s == nil || s.db == nil || !c.IsValid() {
		return false, errors.New("invalid read command/store")
	}
	row := readCommandModel{SubjectType: string(c.Subject.Type), SubjectID: c.Subject.ID, Operation: string(c.Operation), CommandKey: c.CommandKey, RequestHash: &c.RequestHash, CreatedAt: time.Now().UTC()}
	if err := s.db.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(&row).Error; err != nil {
		return false, err
	}
	var saved readCommandModel
	if err := readCommandWhere(s.db.WithContext(ctx), c).First(&saved).Error; err != nil {
		return false, err
	}
	if err := verifyReadCommand(saved, c); err != nil {
		return false, err
	}
	return saved.AcceptedAt != nil, nil
}

// AcceptReadCommand ghi usage và kết quả acceptance nguyên tử. Khóa command
// bao trùm mọi meter/chu kỳ; transaction thua trả reservation của winner.
// Event nil là quyền không tính quota, vẫn giữ acceptance cho retry sau này.
func (s *Store) AcceptReadCommand(ctx context.Context, c usage.ReadCommand, event *usage.Event) (string, error) {
	if s == nil || s.db == nil || !c.IsValid() || (event != nil && (!event.IsValid() || !c.Matches(*event))) {
		return "", errors.New("invalid read acceptance")
	}
	var reservationID string
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var row readCommandModel
		if err := readCommandWhere(tx, c).Clauses(clause.Locking{Strength: "UPDATE"}).First(&row).Error; err != nil {
			return err
		}
		if err := verifyReadCommand(row, c); err != nil {
			return err
		}
		store := NewStore(tx)
		if row.AcceptedAt != nil {
			if row.AcceptedUsageKey != nil {
				saved, found, err := store.FindUsageEventByKey(ctx, *row.AcceptedUsageKey)
				if err != nil {
					return err
				}
				if !found {
					return errors.New("accepted read usage is missing")
				}
				reservationID = saved.ReservationID
			}
			return nil
		}
		var usageKey *string
		if event != nil {
			if _, err := store.SaveUsageEvent(ctx, *event); err != nil {
				return err
			}
			saved, found, err := store.FindUsageEventByKey(ctx, event.UsageKey)
			if err != nil {
				return err
			}
			if !found {
				return errors.New("read usage was not persisted")
			}
			reservationID = saved.ReservationID
			usageKey = &saved.UsageKey
		}
		return readCommandWhere(tx.Model(&readCommandModel{}), c).Updates(map[string]any{"accepted_at": time.Now().UTC(), "accepted_usage_key": usageKey}).Error
	})
	return reservationID, err
}
