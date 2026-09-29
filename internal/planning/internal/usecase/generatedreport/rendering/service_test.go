package rendering

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"testing"
	"time"
	"tqd/internal/usecase/generatedreport/processing"
)

type fakeSource struct {
	document Document
	err      error
}

func (s fakeSource) LoadForRender(context.Context, uint64, uint64, string) (Document, error) {
	return s.document, s.err
}

type fakeEngine struct {
	pdf  []byte
	err  error
	html []byte
}

func (e *fakeEngine) RenderPDF(ctx context.Context, html []byte) ([]byte, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	e.html = append([]byte(nil), html...)
	return append([]byte(nil), e.pdf...), e.err
}

type fakeUploader struct {
	path           string
	err            error
	ownerNamespace string
	ownerKey       string
	filename       string
	contentType    string
	content        []byte
}

func (u *fakeUploader) PutOwnedFile(
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
	return u.path, u.err
}

func renderDocument() Document {
	parcelID := uint64(99)
	return Document{
		ReportID:       12,
		UserID:         42,
		JobID:          "report_job_12",
		Title:          "Báo cáo quy hoạch thửa 88",
		Subtitle:       "Tờ 5 · Phường Minh Khai",
		Address:        "12 Đường Mẫu",
		Province:       "Hà Nội",
		ParcelID:       &parcelID,
		CenterLat:      21.0278,
		CenterLon:      105.8342,
		MinLon:         105.83,
		MinLat:         21.02,
		MaxLon:         105.84,
		MaxLat:         21.03,
		MetadataJSON:   []byte(`{"mapNumber":"5","landNumber":"88","source":"ignored"}`),
		ComparisonJSON: []byte(`{"fromPlanName":"QH 2020","toPlanName":"QH 2030","fromYear":2020,"toYear":2030}`),
		CreatedAt:      time.Date(2026, 8, 19, 8, 30, 0, 0, time.UTC),
	}
}

type fixedDelivery struct {
	url string
}

func (d fixedDelivery) PublicURL(string) (string, error) {
	return d.url, nil
}

func TestRenderHTMLIsDeterministicAndEscapesUntrustedText(t *testing.T) {
	doc := renderDocument()
	doc.Subtitle = `<script>alert("x")</script> · Tiếng Việt`
	first, err := RenderHTML(doc)
	if err != nil {
		t.Fatalf("RenderHTML() error = %v", err)
	}
	second, err := RenderHTML(doc)
	if err != nil {
		t.Fatalf("RenderHTML() retry error = %v", err)
	}
	if !bytes.Equal(first, second) {
		t.Fatal("RenderHTML() is not deterministic")
	}
	text := string(first)
	if strings.Contains(text, `<script>alert`) || !strings.Contains(text, `&lt;script&gt;`) {
		t.Fatal("HTML template did not escape untrusted text")
	}
	for _, expected := range []string{"Báo cáo quy hoạch thửa 88", "Hà Nội", "Tờ bản đồ", "QH 2020 2020 → QH 2030 2030"} {
		if !strings.Contains(text, expected) {
			t.Fatalf("HTML missing %q", expected)
		}
	}
}

func TestServiceGeneratesAndUploadsProductionPDF(t *testing.T) {
	engine := &fakeEngine{pdf: []byte("%PDF-1.7\nproduction\n%%EOF\n")}
	uploader := &fakeUploader{path: "pencoded-path"}
	service, err := NewService(fakeSource{document: renderDocument()}, engine, uploader, fixedDelivery{url: "https://delivery.example/report.pdf"})
	if err != nil {
		t.Fatalf("NewService() error = %v", err)
	}
	output, err := service.GenerateReport(context.Background(), processing.Job{ID: "report_job_12", ReportID: 12, UserID: 42})
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
	if uploader.ownerKey != "report_job_12" {
		t.Fatalf(
			"owner key = %q, want report_job_12",
			uploader.ownerKey,
		)
	}
	if uploader.filename != "qhpro-report-12.pdf" || uploader.contentType != "application/pdf" {
		t.Fatalf("upload metadata = %q %q", uploader.filename, uploader.contentType)
	}
	if !bytes.Equal(uploader.content, engine.pdf) {
		t.Fatal("uploaded bytes differ from PDF engine output")
	}
	if output.PDFURL != "https://delivery.example/report.pdf" {
		t.Fatalf("PDFURL = %q", output.PDFURL)
	}
	if output.FileSize != uint64(len(engine.pdf)) {
		t.Fatalf("FileSize = %d", output.FileSize)
	}
	if output.Format != "pdf" {
		t.Fatalf("Format = %q, want pdf", output.Format)
	}
	if !strings.Contains(string(engine.html), "Báo cáo quy hoạch thửa 88") {
		t.Fatal("engine did not receive accepted report content")
	}
}

func TestServiceFailsClosedOnInvalidEngineOutputOrSnapshotMismatch(t *testing.T) {
	doc := renderDocument()
	engine := &fakeEngine{pdf: []byte("not-pdf")}
	service, _ := NewService(fakeSource{document: doc}, engine, &fakeUploader{path: "p"}, fixedDelivery{url: "https://delivery.example/report.pdf"})
	if _, err := service.GenerateReport(context.Background(), processing.Job{ID: doc.JobID, ReportID: doc.ReportID, UserID: doc.UserID}); err == nil {
		t.Fatal("invalid PDF engine output was accepted")
	}

	doc.JobID = "other"
	engine.pdf = []byte("%PDF-1.7\n%%EOF")
	service, _ = NewService(fakeSource{document: doc}, engine, &fakeUploader{path: "p"}, fixedDelivery{url: "https://delivery.example/report.pdf"})
	if _, err := service.GenerateReport(context.Background(), processing.Job{ID: "report_job_12", ReportID: 12, UserID: 42}); err == nil {
		t.Fatal("snapshot/job mismatch was accepted")
	}
}

func TestServicePropagatesCancellationAndSourceErrors(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	service, _ := NewService(fakeSource{document: renderDocument()}, &fakeEngine{pdf: []byte("%PDF-1.7")}, &fakeUploader{path: "p"}, fixedDelivery{url: "https://delivery.example/report.pdf"})
	if _, err := service.GenerateReport(ctx, processing.Job{ID: "report_job_12", ReportID: 12, UserID: 42}); !errors.Is(err, context.Canceled) {
		t.Fatalf("cancel error = %v", err)
	}

	service, _ = NewService(fakeSource{err: errors.New("db down")}, &fakeEngine{pdf: []byte("%PDF-1.7")}, &fakeUploader{path: "p"}, fixedDelivery{url: "https://delivery.example/report.pdf"})
	if _, err := service.GenerateReport(context.Background(), processing.Job{ID: "report_job_12", ReportID: 12, UserID: 42}); err == nil || !strings.Contains(err.Error(), "load accepted report snapshot") {
		t.Fatalf("source error = %v", err)
	}
}
