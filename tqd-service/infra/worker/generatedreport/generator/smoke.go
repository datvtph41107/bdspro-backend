package generator

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"tqd/internal/usecase/generatedreport/processing"
)

// Uploader is the smallest file-service contract needed by the smoke renderer.
type Uploader interface {
	PutOwnedFile(
		ctx context.Context,
		ownerNamespace string,
		ownerKey string,
		file io.Reader,
		filename string,
		contentType string,
	) (string, error)
}

// DeliveryReference is the public File reference capability required by the
// smoke renderer after owned upload.
type DeliveryReference interface {
	PublicURL(path string) (string, error)
}

// SmokeGenerator creates a deterministic, minimal PDF for local end-to-end tests.
// It is deliberately not a production planning renderer and is only wired when
// QHPRO_REPORT_GENERATOR_MODE=smoke is explicitly configured.
type SmokeGenerator struct {
	uploader Uploader
	delivery DeliveryReference
}

func NewSmokeGenerator(
	uploader Uploader,
	delivery DeliveryReference,
) (*SmokeGenerator, error) {
	if uploader == nil || delivery == nil {
		return nil, errors.New(
			"smoke report File dependencies are required",
		)
	}

	return &SmokeGenerator{
		uploader: uploader,
		delivery: delivery,
	}, nil
}

func (g *SmokeGenerator) GenerateReport(ctx context.Context, job processing.Job) (processing.Output, error) {
	if g == nil ||
		g.uploader == nil ||
		g.delivery == nil {
		return processing.Output{}, errors.New("smoke report generator is not configured")
	}
	if job.ID == "" || job.ReportID == 0 || job.UserID == 0 {
		return processing.Output{}, errors.New("smoke report job is invalid")
	}

	content := minimalPDF([]string{
		"QHPro local end-to-end smoke report",
		"Report ID: " + strconv.FormatUint(job.ReportID, 10),
		"User ID: " + strconv.FormatUint(job.UserID, 10),
		"Operation: " + string(job.Operation),
		"Operation ID: " + job.OperationID,
		"Job ID: " + job.ID,
	})
	filename := fmt.Sprintf("qhpro-report-%d.pdf", job.ReportID)
	path, err := g.uploader.PutOwnedFile(
		ctx,
		processing.FileOwnerNamespace,
		job.ID,
		bytes.NewReader(content),
		filename,
		"application/pdf",
	)
	if err != nil {
		return processing.Output{}, fmt.Errorf("upload smoke report: %w", err)
	}
	path = strings.TrimSpace(path)
	if path == "" {
		return processing.Output{}, errors.New("upload smoke report returned an empty path")
	}

	pdfURL, err := g.delivery.PublicURL(path)
	if err != nil {
		return processing.Output{}, fmt.Errorf(
			"build smoke report PDF delivery URL: %w",
			err,
		)
	}
	pdfURL = strings.TrimSpace(pdfURL)
	if pdfURL == "" {
		return processing.Output{}, errors.New(
			"smoke report PDF delivery URL is empty",
		)
	}

	return processing.Output{
		PDFURL:   pdfURL,
		FileSize: uint64(len(content)),
		Format:   "pdf",
	}, nil
}

func minimalPDF(lines []string) []byte {
	var content strings.Builder
	content.WriteString("BT\n/F1 14 Tf\n72 760 Td\n")
	for index, line := range lines {
		if index > 0 {
			content.WriteString("0 -22 Td\n")
		}
		content.WriteString("(")
		content.WriteString(escapePDFText(line))
		content.WriteString(") Tj\n")
	}
	content.WriteString("ET\n")
	stream := content.String()

	objects := []string{
		"<< /Type /Catalog /Pages 2 0 R >>",
		"<< /Type /Pages /Kids [3 0 R] /Count 1 >>",
		"<< /Type /Page /Parent 2 0 R /MediaBox [0 0 612 792] /Resources << /Font << /F1 4 0 R >> >> /Contents 5 0 R >>",
		"<< /Type /Font /Subtype /Type1 /BaseFont /Helvetica >>",
		fmt.Sprintf("<< /Length %d >>\nstream\n%sendstream", len(stream), stream),
	}

	var document bytes.Buffer
	document.WriteString("%PDF-1.4\n")
	offsets := make([]int, len(objects)+1)
	for index, object := range objects {
		offsets[index+1] = document.Len()
		fmt.Fprintf(&document, "%d 0 obj\n%s\nendobj\n", index+1, object)
	}
	xrefOffset := document.Len()
	fmt.Fprintf(&document, "xref\n0 %d\n", len(objects)+1)
	document.WriteString("0000000000 65535 f \n")
	for index := 1; index <= len(objects); index++ {
		fmt.Fprintf(&document, "%010d 00000 n \n", offsets[index])
	}
	fmt.Fprintf(
		&document,
		"trailer\n<< /Size %d /Root 1 0 R >>\nstartxref\n%d\n%%%%EOF\n",
		len(objects)+1,
		xrefOffset,
	)
	return document.Bytes()
}

func escapePDFText(value string) string {
	value = strings.Map(func(character rune) rune {
		if character < 32 || character > 126 {
			return '?'
		}
		return character
	}, value)
	value = strings.ReplaceAll(value, "\\", "\\\\")
	value = strings.ReplaceAll(value, "(", "\\(")
	return strings.ReplaceAll(value, ")", "\\)")
}
