package postgresreport

import (
	"context"
	"errors"
	"fmt"
	"tqd/infra/postgres/generatedreport/job"
	"tqd/infra/postgres/usage"
	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/usecase/generatedreport/application"
	"tqd/internal/usecase/generatedreport/processing"

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
) (application.Report, bool, error) {
	if s == nil || s.db == nil {
		return application.Report{}, false, errors.New("report database is not configured")
	}
	if userID == 0 || commandKey == "" {
		return application.Report{}, false, errors.New("report command is invalid")
	}

	row, err := findByCommand(s.db.WithContext(ctx), userID, commandKey)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return application.Report{}, false, nil
	}
	if err != nil {
		return application.Report{}, false, err
	}
	return fromRow(row), true, nil
}

/**
 * Accept là durable acceptance transaction của create-report.
 *
 * Transaction chỉ làm ba việc:
 *   - arbitrate/persist Report command;
 *   - persist optional Usage evidence;
 *   - persist Report Job đã được application tạo intent.
 *
 * Redis reservation lifecycle không được quyết ở adapter này.
 */
func (s *ReportStore) AcceptReport(
	ctx context.Context,
	input application.Acceptance,
) (application.AcceptanceResult, error) {
	if s == nil || s.db == nil {
		return application.AcceptanceResult{}, errors.New("report database is not configured")
	}

	reportInput := input.Report
	if reportInput.UserID == 0 || reportInput.CommandKey == "" || !reportInput.ReportType.IsValid() ||
		!reportInput.Operation.IsValid() || reportInput.OperationID == "" || reportInput.JobID == "" {
		return application.AcceptanceResult{}, errors.New("report data is invalid")
	}
	if input.Job.ID == "" || input.Job.UserID != reportInput.UserID ||
		input.Job.Operation != reportInput.Operation || input.Job.OperationID != reportInput.OperationID ||
		input.Job.CommandKey != reportInput.CommandKey || input.Job.Status != processing.StatusPending ||
		input.Job.AvailableAt.IsZero() || input.Job.CreatedAt.IsZero() || input.Job.UpdatedAt.IsZero() {
		return application.AcceptanceResult{}, errors.New("report acceptance job is invalid")
	}

	var accepted application.AcceptanceResult

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		created := false
		row := toRow(reportInput)
		result := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&row)
		if result.Error != nil {
			return fmt.Errorf("create generated report: %w", result.Error)
		}

		if result.RowsAffected == 0 {
			existing, err := findByCommand(tx, reportInput.UserID, reportInput.CommandKey)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return errors.New("report conflict occurred but existing report was not found")
				}
				return err
			}
			if existing.RequestHash != "" && reportInput.RequestHash != "" && existing.RequestHash != reportInput.RequestHash {
				return application.ErrCommandConflict
			}
			row = existing
		} else {
			created = true
		}

		if row.JobID == "" {
			row.JobID = reportInput.JobID
			if err := tx.Model(&qh_domain.QHUserReported{}).
				Where("id = ?", row.ID).
				Update("job_id", row.JobID).Error; err != nil {
				return fmt.Errorf("set generated report job ID: %w", err)
			}
		}

		if input.Usage != nil {
			usageStore := postgresusage.NewStore(tx)
			if created {
				usageCreated, err := usageStore.SaveUsageEvent(ctx, *input.Usage)
				if err != nil {
					return fmt.Errorf("save quota usage event: %w", err)
				}
				if !usageCreated {
					return errors.New("new report did not create its quota usage event")
				}
				usage := *input.Usage
				accepted.DurableUsage = &usage
			} else {
				existingUsage, found, err := usageStore.FindUsageEventByKey(ctx, input.Usage.UsageKey)
				if err != nil {
					return fmt.Errorf("load existing quota usage event: %w", err)
				}
				if found {
					usage := existingUsage
					accepted.DurableUsage = &usage
				}
			}
		}

		job := input.Job
		job.ID = row.JobID
		job.ReportID = row.ID
		job.UserID = row.UserID

		jobStore := postgresjob.NewStore(tx)

		// A new acceptance creates its Job in this transaction.
		// A concurrent replay validates the Job owned by the durable winner.
		// OperationID belongs to the winning durable lineage, not to the
		// transport attempt that happened to observe the winner.
		var durableJob processing.Job

		if created {
			jobCreated, err := jobStore.SaveJob(ctx, job)
			if err != nil {
				return fmt.Errorf("save generated report job: %w", err)
			}
			if !jobCreated {
				return errors.New("new report did not create its processing job")
			}
			durableJob = job
		} else {
			existingJob, found, err := jobStore.FindJobByReportID(ctx, row.ID)
			if err != nil {
				return fmt.Errorf("load existing report processing job: %w", err)
			}
			if !found {
				return errors.New("existing accepted report has no processing job")
			}

			if !sameAcceptedJobOwner(existingJob, job) {
				return errors.New("existing report processing job does not match accepted command")
			}

			// Usage and Job are both durable winner facts. When metering
			// evidence exists, their lineage must agree with each other.
			if accepted.DurableUsage != nil &&
				existingJob.OperationID != accepted.DurableUsage.OperationID {
				return errors.New(
					"existing report processing job does not match durable usage lineage",
				)
			}

			durableJob = existingJob
		}

		saved := fromRow(row)
		saved.Operation = reportInput.Operation
		saved.OperationID = durableJob.OperationID
		accepted.Report = saved
		accepted.Created = created
		return nil
	})
	if err != nil {
		return application.AcceptanceResult{}, err
	}

	return accepted, nil
}

// sameAcceptedJobOwner validates the stable ownership of a Report Job.
//
// OperationID is intentionally not compared here. Concurrent attempts for one
// CommandKey may carry different transport Operation-IDs after they all miss
// the initial durable precheck. The winner's durable Operation-ID is validated
// against durable Usage evidence when metering is present.
func sameAcceptedJobOwner(actual, expected processing.Job) bool {
	return actual.ID == expected.ID &&
		actual.ReportID == expected.ReportID &&
		actual.UserID == expected.UserID &&
		actual.Operation == expected.Operation &&
		actual.CommandKey == expected.CommandKey
}

func findByCommand(
	tx *gorm.DB,
	userID uint64,
	commandKey string,
) (qh_domain.QHUserReported, error) {
	var row qh_domain.QHUserReported
	err := tx.
		Where("user_id = ? AND command_key = ? AND deleted_at IS NULL", userID, commandKey).
		First(&row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return qh_domain.QHUserReported{}, gorm.ErrRecordNotFound
		}
		return qh_domain.QHUserReported{}, fmt.Errorf("find generated report by command: %w", err)
	}
	return row, nil
}

func toRow(input application.Report) qh_domain.QHUserReported {
	return qh_domain.QHUserReported{
		UserID:       input.UserID,
		ReportType:   input.ReportType,
		Status:       input.Status,
		Title:        input.Title,
		Subtitle:     input.Subtitle,
		Address:      input.Location.Address,
		Province:     input.Location.Province,
		ProvinceCode: input.Location.ProvinceCode,
		WardCode:     input.Location.WardCode,
		ParcelID:     input.ParcelID,
		RegionID:     input.RegionID,
		MinLon:       input.Spatial.Bounds.MinLon,
		MinLat:       input.Spatial.Bounds.MinLat,
		MaxLon:       input.Spatial.Bounds.MaxLon,
		MaxLat:       input.Spatial.Bounds.MaxLat,
		CenterLat:    input.Spatial.Centroid.Lat,
		CenterLon:    input.Spatial.Centroid.Lon,
		Format:       input.Format,
		Comparison:   datatypes.JSON(input.ComparisonJSON),
		Metadata:     datatypes.JSON(input.MetadataJSON),
		CommandKey:   input.CommandKey,
		RequestHash:  input.RequestHash,
		JobID:        input.JobID,
	}
}

func fromRow(row qh_domain.QHUserReported) application.Report {
	return application.Report{
		ID:         row.ID,
		UserID:     row.UserID,
		ReportType: row.ReportType,
		Status:     row.Status,
		Title:      row.Title,
		Subtitle:   row.Subtitle,
		Location: application.Location{
			Address:      row.Address,
			Province:     row.Province,
			ProvinceCode: row.ProvinceCode,
			WardCode:     row.WardCode,
		},
		Spatial: application.Spatial{
			Bounds: application.Bounds{
				MinLon: row.MinLon,
				MinLat: row.MinLat,
				MaxLon: row.MaxLon,
				MaxLat: row.MaxLat,
			},
			Centroid: application.Point{
				Lat: row.CenterLat,
				Lon: row.CenterLon,
			},
		},
		ParcelID:       row.ParcelID,
		RegionID:       row.RegionID,
		Format:         row.Format,
		ComparisonJSON: append([]byte(nil), row.Comparison...),
		MetadataJSON:   append([]byte(nil), row.Metadata...),
		CommandKey:     row.CommandKey,
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
