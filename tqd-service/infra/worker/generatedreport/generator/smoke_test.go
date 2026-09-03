package generator

import (
	"bytes"
	"context"
	"io"
	"testing"

	"tqd/internal/usecase/generatedreport/processing"
)

type recordingUploader struct {
	ownerNamespace string
	ownerKey       string
	contentType    string
	filename       string
	content        []byte
	path           string
}

func (u *recordingUploader) PutOwnedFile(
	_ context.Context,
	ownerNamespace string,
	ownerKey string,
	file io.Reader,
	filename string,
	contentType string,
) (string, error) {
	u.ownerNamespace = ownerNamespace
	u.ownerKey = ownerKey
	u.filename = filename
	u.contentType = contentType
	u.content, _ = io.ReadAll(file)
	return u.path, nil
}

type fixedDelivery struct {
	url string
}

func (d fixedDelivery) PublicURL(string) (string, error) {
	return d.url, nil
}

func TestSmokeGeneratorCreatesDownloadablePDFOutput(t *testing.T) {
	uploader := &recordingUploader{path: "pencoded-path"}
	generator, err := NewSmokeGenerator(
		uploader,
		fixedDelivery{url: "https://delivery.example/report.pdf"},
	)
	if err != nil {
		t.Fatalf("NewSmokeGenerator() error = %v", err)
	}

	output, err := generator.GenerateReport(context.Background(), processing.Job{
		ID:          "report_job_test",
		ReportID:    12,
		UserID:      42,
		Operation:   "workspace.report.generate",
		OperationID: "op-smoke-test",
	})
	if err != nil {
		t.Fatalf("GenerateReport() error = %v", err)
	}
	if uploader.ownerNamespace != processing.FileOwnerNamespace {
		t.Fatalf(
			"owner namespace = %q, want %q",
			uploader.ownerNamespace,
			processing.FileOwnerNamespace,
		)
	}
	if uploader.ownerKey != "report_job_test" {
		t.Fatalf(
			"owner key = %q, want report_job_test",
			uploader.ownerKey,
		)
	}
	if uploader.filename != "qhpro-report-12.pdf" || uploader.contentType != "application/pdf" {
		t.Fatalf("upload metadata = %q, %q", uploader.filename, uploader.contentType)
	}
	if !bytes.HasPrefix(uploader.content, []byte("%PDF-1.4")) ||
		!bytes.HasSuffix(uploader.content, []byte("%%EOF\n")) {
		t.Fatal("generated content is not a minimal PDF")
	}
	if output.PDFURL != "https://delivery.example/report.pdf" {
		t.Fatalf("PDFURL = %q", output.PDFURL)
	}
	if output.FileSize != uint64(len(uploader.content)) {
		t.Fatalf("FileSize = %d, want %d", output.FileSize, len(uploader.content))
	}
	if output.Format != "pdf" {
		t.Fatalf("Format = %q, want pdf", output.Format)
	}
}

func TestSmokeGeneratorRequiresValidConfiguration(t *testing.T) {
	delivery := fixedDelivery{
		url: "https://delivery.example/report.pdf",
	}

	if _, err := NewSmokeGenerator(nil, delivery); err == nil {
		t.Fatal("NewSmokeGenerator(nil uploader) error = nil")
	}

	if _, err := NewSmokeGenerator(
		&recordingUploader{},
		nil,
	); err == nil {
		t.Fatal("NewSmokeGenerator(nil delivery) error = nil")
	}
}
