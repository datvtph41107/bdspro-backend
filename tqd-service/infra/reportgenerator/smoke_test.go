package reportgenerator

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"

	"common/businessoperation"
	"tqd/internal/reportjob"
)

type recordingUploader struct {
	contentType string
	filename    string
	content     []byte
	path        string
}

func (u *recordingUploader) UploadFile(
	_ context.Context,
	file io.Reader,
	filename string,
	contentType string,
) (string, error) {
	u.filename = filename
	u.contentType = contentType
	u.content, _ = io.ReadAll(file)
	return u.path, nil
}

func TestSmokeGeneratorCreatesDownloadablePDFOutput(t *testing.T) {
	uploader := &recordingUploader{path: "pencoded-path"}
	generator, err := NewSmokeGenerator(uploader, "http://127.0.0.1:8002/")
	if err != nil {
		t.Fatalf("NewSmokeGenerator() error = %v", err)
	}

	output, err := generator.GenerateReport(context.Background(), reportjob.Job{
		ID:          "report_job_test",
		ReportID:    12,
		UserID:      42,
		Operation:   businessoperation.WorkspaceReportGenerate,
		OperationID: "op-smoke-test",
	})
	if err != nil {
		t.Fatalf("GenerateReport() error = %v", err)
	}
	if uploader.filename != "qhpro-report-12.pdf" || uploader.contentType != "application/pdf" {
		t.Fatalf("upload metadata = %q, %q", uploader.filename, uploader.contentType)
	}
	if !bytes.HasPrefix(uploader.content, []byte("%PDF-1.4")) ||
		!bytes.HasSuffix(uploader.content, []byte("%%EOF\n")) {
		t.Fatal("generated content is not a minimal PDF")
	}
	if !strings.Contains(output.PDFURL, "/v1/file/load/pencoded-path/local") {
		t.Fatalf("PDFURL = %q", output.PDFURL)
	}
	if output.FileSize != uint64(len(uploader.content)) {
		t.Fatalf("FileSize = %d, want %d", output.FileSize, len(uploader.content))
	}
}

func TestSmokeGeneratorRequiresValidConfiguration(t *testing.T) {
	if _, err := NewSmokeGenerator(nil, "http://127.0.0.1:8002"); err == nil {
		t.Fatal("NewSmokeGenerator(nil) error = nil")
	}
	if _, err := NewSmokeGenerator(&recordingUploader{}, "not-a-url"); err == nil {
		t.Fatal("NewSmokeGenerator(invalid URL) error = nil")
	}
}
