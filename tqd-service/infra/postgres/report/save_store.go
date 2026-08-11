package reportpostgres

import (
	"context"
	"errors"
	"fmt"
	"time"
	reportjobpostgres "tqd/infra/postgres/reportjob"
	usageeventpostgres "tqd/infra/postgres/usageevent"
	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/report"
	"tqd/internal/reportjob"
	"tqd/internal/usage"

	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

/**
 * ReportStore giữ durable state của generated report trong PostgreSQL.
 */
type ReportStore struct {
	db *gorm.DB
}

func NewReportStore(db *gorm.DB) *ReportStore {
	return &ReportStore{db: db}
}

/**
 * FindReportByCommand tìm report đã accepted cho cùng user + command key.
 */
func (s *ReportStore) FindReportByCommand(
	ctx context.Context,
	userID uint64,
	commandKey string,
) (report.Report, bool, error) {
	if s == nil || s.db == nil {
		return report.Report{}, false, errors.New("report database is not configured")
	}
	if userID == 0 || commandKey == "" {
		return report.Report{}, false, errors.New("report command is invalid")
	}

	row, err := findByCommand(s.db.WithContext(ctx), userID, commandKey)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return report.Report{}, false, nil
	}
	if err != nil {
		return report.Report{}, false, err
	}
	return fromRow(row), true, nil
}

/**
 * AcceptReport lưu report và usage event trong cùng transaction.
 *
 * Cùng command chạy lại sẽ nhận report cũ và không tạo thêm usage event.
 */
func (s *ReportStore) AcceptReport(
	ctx context.Context,
	input report.Report,
	usageEvent *usage.Event,
) (report.Report, bool, bool, error) {
	if s == nil || s.db == nil {
		return report.Report{}, false, false, errors.New("report database is not configured")
	}
	if input.UserID == 0 || input.CommandKey == "" || !input.ReportType.IsValid() ||
		!input.Operation.IsValid() || input.OperationID == "" || input.JobID == "" {
		return report.Report{}, false, false, errors.New("report data is invalid")
	}

	var saved report.Report
	created := false
	finalizeReservation := usageEvent == nil

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		row := toRow(input)
		result := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&row)
		if result.Error != nil {
			return fmt.Errorf("create generated report: %w", result.Error)
		}

		if result.RowsAffected == 0 {
			existing, err := findByCommand(tx, input.UserID, input.CommandKey)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return errors.New("report conflict occurred but existing report was not found")
				}
				return err
			}
			if existing.RequestHash != "" && input.RequestHash != "" && existing.RequestHash != input.RequestHash {
				return report.ErrCommandConflict
			}
			row = existing
		} else {
			created = true
		}

		if row.JobID == "" {
			row.JobID = input.JobID
			if err := tx.Model(&qh_domain.QHUserReported{}).
				Where("id = ?", row.ID).
				Update("job_id", row.JobID).Error; err != nil {
				return fmt.Errorf("set generated report job ID: %w", err)
			}
		}

		if usageEvent != nil {
			usageStore := usageeventpostgres.NewStore(tx)
			if created {
				usageCreated, err := usageStore.SaveUsageEvent(ctx, *usageEvent)
				if err != nil {
					return fmt.Errorf("save quota usage event: %w", err)
				}
				if !usageCreated {
					return errors.New("new report did not create its quota usage event")
				}
				finalizeReservation = true
			} else {
				existingUsage, found, err := usageStore.FindUsageEventByKey(ctx, usageEvent.UsageKey)
				if err != nil {
					return fmt.Errorf("load existing quota usage event: %w", err)
				}
				if found {
					finalizeReservation = existingUsage.ReservationID == usageEvent.ReservationID
				}
			}
		}

		jobTime := time.Now().UTC()
		jobStore := reportjobpostgres.NewStore(tx)
		if _, err := jobStore.SaveJob(ctx, reportjob.Job{
			ID:          row.JobID,
			ReportID:    row.ID,
			UserID:      row.UserID,
			Operation:   input.Operation,
			OperationID: input.OperationID,
			CommandKey:  input.CommandKey,
			Status:      reportjob.StatusPending,
			AvailableAt: jobTime,
			CreatedAt:   jobTime,
			UpdatedAt:   jobTime,
		}); err != nil {
			return fmt.Errorf("save generated report job: %w", err)
		}

		saved = fromRow(row)
		saved.Operation = input.Operation
		saved.OperationID = input.OperationID
		return nil
	})
	if err != nil {
		return report.Report{}, false, false, err
	}

	return saved, created, finalizeReservation, nil
}

func findByCommand(
	tx *gorm.DB,
	userID uint64,
	commandKey string,
) (qh_domain.QHUserReported, error) {
	var row qh_domain.QHUserReported
	err := tx.
		Where("user_id = ? AND client_request_id = ? AND deleted_at IS NULL", userID, commandKey).
		First(&row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return qh_domain.QHUserReported{}, gorm.ErrRecordNotFound
		}
		return qh_domain.QHUserReported{}, fmt.Errorf("find generated report by command: %w", err)
	}
	return row, nil
}

func toRow(input report.Report) qh_domain.QHUserReported {
	return qh_domain.QHUserReported{
		UserID:          input.UserID,
		ReportType:      input.ReportType,
		Status:          input.Status,
		Title:           input.Title,
		Subtitle:        input.Subtitle,
		Address:         input.Location.Address,
		Province:        input.Location.Province,
		ProvinceCode:    input.Location.ProvinceCode,
		WardCode:        input.Location.WardCode,
		ParcelID:        input.ParcelID,
		RegionID:        input.RegionID,
		MinLon:          input.Spatial.Bounds.MinLon,
		MinLat:          input.Spatial.Bounds.MinLat,
		MaxLon:          input.Spatial.Bounds.MaxLon,
		MaxLat:          input.Spatial.Bounds.MaxLat,
		CenterLat:       input.Spatial.Centroid.Lat,
		CenterLon:       input.Spatial.Centroid.Lon,
		Format:          input.Format,
		Comparison:      datatypes.JSON(input.ComparisonJSON),
		Metadata:        datatypes.JSON(input.MetadataJSON),
		ClientRequestID: input.CommandKey,
		RequestHash:     input.RequestHash,
		JobID:           input.JobID,
	}
}

func fromRow(row qh_domain.QHUserReported) report.Report {
	return report.Report{
		ID:         row.ID,
		UserID:     row.UserID,
		ReportType: row.ReportType,
		Status:     row.Status,
		Title:      row.Title,
		Subtitle:   row.Subtitle,
		Location: report.Location{
			Address:      row.Address,
			Province:     row.Province,
			ProvinceCode: row.ProvinceCode,
			WardCode:     row.WardCode,
		},
		Spatial: report.Spatial{
			Bounds: report.Bounds{
				MinLon: row.MinLon,
				MinLat: row.MinLat,
				MaxLon: row.MaxLon,
				MaxLat: row.MaxLat,
			},
			Centroid: report.Point{
				Lat: row.CenterLat,
				Lon: row.CenterLon,
			},
		},
		ParcelID:       row.ParcelID,
		RegionID:       row.RegionID,
		Format:         row.Format,
		ComparisonJSON: append([]byte(nil), row.Comparison...),
		MetadataJSON:   append([]byte(nil), row.Metadata...),
		CommandKey:     row.ClientRequestID,
		RequestHash:    row.RequestHash,
		JobID:          row.JobID,
		ThumbnailURL:   row.ThumbnailURL,
		ImageURL:       row.ImageURL,
		PDFURL:         row.PDFURL,
		ShareURL:       row.ShareURL,
		FileSize:       row.FileSize,
		ErrorMessage:   row.ErrorMessage,
		CreatedAt:      row.CreatedAt,
		UpdatedAt:      row.UpdatedAt,
		ExpiresAt:      row.ExpiresAt,
	}
}
