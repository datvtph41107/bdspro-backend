package rendering

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"
	"tqd/internal/usecase/generatedreport/processing"
)

// Service is the production Generated Report renderer. It owns orchestration
// across the accepted Report snapshot, PDF engine and File-service upload.
type Service struct {
	source   Source
	engine   PDFEngine
	uploader Uploader
	delivery DeliveryReference
}

func NewService(
	source Source,
	engine PDFEngine,
	uploader Uploader,
	delivery DeliveryReference,
) (*Service, error) {
	if source == nil ||
		engine == nil ||
		uploader == nil ||
		delivery == nil {
		return nil, errors.New(
			"production report renderer dependencies are required",
		)
	}

	return &Service{
		source:   source,
		engine:   engine,
		uploader: uploader,
		delivery: delivery,
	}, nil
}

func (s *Service) GenerateReport(ctx context.Context, job processing.Job) (processing.Output, error) {
	if s == nil ||
		s.source == nil ||
		s.engine == nil ||
		s.uploader == nil ||
		s.delivery == nil {
		return processing.Output{}, errors.New("production report renderer is not configured")
	}
	if ctx == nil {
		return processing.Output{}, errors.New("production report context is nil")
	}
	if err := ctx.Err(); err != nil {
		return processing.Output{}, err
	}
	if job.ID == "" || job.ReportID == 0 || job.UserID == 0 {
		return processing.Output{}, errors.New("production report job is invalid")
	}

	document, err := s.source.LoadForRender(ctx, job.ReportID, job.UserID, job.ID)
	if err != nil {
		return processing.Output{}, fmt.Errorf("load accepted report snapshot: %w", err)
	}
	if document.ReportID != job.ReportID || document.UserID != job.UserID || document.JobID != job.ID {
		return processing.Output{}, errors.New("accepted report snapshot does not match claimed job")
	}

	html, err := RenderHTML(document)
	if err != nil {
		return processing.Output{}, err
	}
	pdf, err := s.engine.RenderPDF(ctx, html)
	if err != nil {
		return processing.Output{}, fmt.Errorf("render report PDF: %w", err)
	}
	if len(pdf) < 8 || !bytes.HasPrefix(pdf, []byte("%PDF-")) {
		return processing.Output{}, errors.New("report PDF engine returned invalid content")
	}

	filename := fmt.Sprintf("qhpro-report-%d.pdf", job.ReportID)
	path, err := s.uploader.PutOwnedFile(
		ctx,
		processing.FileOwnerNamespace,
		job.ID,
		bytes.NewReader(pdf),
		filename,
		"application/pdf",
	)
	if err != nil {
		return processing.Output{}, fmt.Errorf("upload generated report PDF: %w", err)
	}
	path = strings.TrimSpace(path)
	if path == "" {
		return processing.Output{}, errors.New("upload generated report PDF returned an empty path")
	}

	pdfURL, err := s.delivery.PublicURL(path)
	if err != nil {
		return processing.Output{}, fmt.Errorf(
			"build generated report PDF delivery URL: %w",
			err,
		)
	}
	pdfURL = strings.TrimSpace(pdfURL)
	if pdfURL == "" {
		return processing.Output{}, errors.New(
			"generated report PDF delivery URL is empty",
		)
	}

	return processing.Output{
		PDFURL:   pdfURL,
		FileSize: uint64(len(pdf)),
		Format:   "pdf",
	}, nil
}
