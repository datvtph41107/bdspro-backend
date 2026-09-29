package qh_usecase

import (
	"bytes"
	_db "common/db"
	_err "common/domain/err"
	"context"
	"fmt"
	"log/slog"
	"path/filepath"
	"time"
	qh_domain "tqd/internal/domain/qh"
	qh_dto "tqd/internal/dto/qh"
	"tqd/internal/enums"
	"tqd/internal/interface/provider"
	"tqd/internal/interface/repo"
)

// IQHPlanningFolderUsecase tạo 1 đồ án quy hoạch từ 1 folder admin upload: mỗi file trong folder
// trở thành 1 QHPlanningDocument (trạng thái Pending, chờ job classify quét và gọi Gemini phân loại).
// Đồng thời tạo ProAIJob type Đồ án gắn planning_project_id.
type IQHPlanningFolderUsecase interface {
	CreateFromFolder(ctx context.Context, project *qh_domain.QHPlanningProject, folderName string, files []qh_dto.FolderFileInput, jobType uint32) (*qh_domain.QHPlanningProject, []qh_domain.QHPlanningDocument, *qh_domain.ProAIJob, error)
}

type qhPlanningFolderUsecase struct {
	projectRepo  repo.IQHPlanningProjectRepo
	documentRepo repo.IQHPlanningDocumentRepo
	aiJobRepo    repo.IProAIJobRepo
	fileProvider provider.FileProvider
	tx           *_db.TransactionRepo
}

func NewQHPlanningFolderUsecase(
	projectRepo repo.IQHPlanningProjectRepo,
	documentRepo repo.IQHPlanningDocumentRepo,
	aiJobRepo repo.IProAIJobRepo,
	fileProvider provider.FileProvider,
	tx *_db.TransactionRepo,
) IQHPlanningFolderUsecase {
	return &qhPlanningFolderUsecase{
		projectRepo:  projectRepo,
		documentRepo: documentRepo,
		aiJobRepo:    aiJobRepo,
		fileProvider: fileProvider,
		tx:           tx,
	}
}

func (u *qhPlanningFolderUsecase) CreateFromFolder(
	ctx context.Context,
	project *qh_domain.QHPlanningProject,
	folderName string,
	files []qh_dto.FolderFileInput,
	jobType uint32,
) (*qh_domain.QHPlanningProject, []qh_domain.QHPlanningDocument, *qh_domain.ProAIJob, error) {
	if len(files) == 0 {
		return nil, nil, nil, &_err.ErrorDTO{Code: 400, Message: "folder không có file nào"}
	}

	if jobType == 0 {
		jobType = enums.ProAIJobTypePlanningProject
	}
	if !enums.IsSupportedProAIJobType(jobType) {
		return nil, nil, nil, &_err.ErrorDTO{Code: 400, Message: "job_type không được hỗ trợ"}
	}
	// Hiện chỉ type Đồ án mới đẩy sang project/document.
	if jobType != enums.ProAIJobTypePlanningProject {
		return nil, nil, nil, &_err.ErrorDTO{Code: 400, Message: "chỉ hỗ trợ job_type Đồ án (10)"}
	}

	if project.ValidityStatus == "" {
		project.ValidityStatus = "DRAFT"
	}
	if project.Code == "" {
		project.Code = fmt.Sprintf("DA-%d", time.Now().UnixMilli())
	}
	if project.Name == "" {
		project.Name = folderName
	}
	if project.PlanningType == 0 {
		project.PlanningType = enums.PlanningTypeGeneral
	}
	if project.PlanningLevel == 0 {
		project.PlanningLevel = enums.PlanningLevelProject
	}
	project.SourceFolderName = folderName
	project.SourceFolderPath = fmt.Sprintf("planning/%s", folderName)
	project.ProcessStatus = enums.PlanningProcessStatus(enums.PlanningProcessStatusPending)
	if project.Metadata == "" {
		project.Metadata = "{}"
	}

	// Upload file ngoài transaction để không giữ DB lock lâu; lỗi upload vẫn tạo doc Failed.
	pending := make([]qh_domain.QHPlanningDocument, 0, len(files))
	for _, f := range files {
		doc := qh_domain.QHPlanningDocument{
			DocumentType:   enums.PlanningDocumentType(enums.PlanningDocumentTypeOther),
			Title:          filepath.Base(f.RelativePath),
			RelativePath:   f.RelativePath,
			ValidityStatus: "DRAFT",
			ProcessStatus:  enums.PlanningProcessStatus(enums.PlanningProcessStatusPending),
			Metadata:       "{}",
		}

		contentType := f.ContentType
		if contentType == "" {
			contentType = "application/octet-stream"
		}
		path, err := u.fileProvider.UploadFile(ctx, bytes.NewReader(f.Content), f.RelativePath, contentType)
		if err != nil {
			slog.ErrorContext(ctx, fmt.Sprintf("[QHPlanningFolder] upload lỗi file %s: %v", f.RelativePath, err))
			doc.ProcessStatus = enums.PlanningProcessStatus(enums.PlanningProcessStatusFailed)
			doc.ClassifyError = "upload lỗi: " + err.Error()
		} else {
			doc.Filepath = path
		}
		pending = append(pending, doc)
	}

	var aiJob *qh_domain.ProAIJob
	err := u.tx.WithTransaction(ctx, func(txCtx context.Context) error {
		if err := u.projectRepo.Create(txCtx, project); err != nil {
			return &_err.ErrorDTO{Code: 500, Message: "Không thể tạo đồ án: " + err.Error()}
		}

		for i := range pending {
			pending[i].PlanningProjectID = project.ID
			if err := u.documentRepo.Create(txCtx, &pending[i]); err != nil {
				return &_err.ErrorDTO{Code: 500, Message: fmt.Sprintf("Không thể lưu tài liệu %s: %v", pending[i].RelativePath, err)}
			}
		}

		projectID := project.ID
		job := &qh_domain.ProAIJob{
			JobType:           jobType,
			Name:              project.Name,
			SourceFolderName:  folderName,
			ProcessStatus:     enums.PlanningProcessStatus(enums.PlanningProcessStatusPending),
			PlanningProjectID: &projectID,
			Metadata:          "{}",
		}
		if err := u.aiJobRepo.Create(txCtx, job); err != nil {
			return &_err.ErrorDTO{Code: 500, Message: "Không thể tạo job AI: " + err.Error()}
		}
		aiJob = job
		return nil
	})
	if err != nil {
		return nil, nil, nil, err
	}

	return project, pending, aiJob, nil
}
