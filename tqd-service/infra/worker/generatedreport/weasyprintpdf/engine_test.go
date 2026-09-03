package weasyprintpdf

import (
	"bytes"
	"context"
	"errors"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func TestEngineRendersUTF8HTMLWhenWeasyPrintIsAvailable(t *testing.T) {
	binary, err := exec.LookPath("weasyprint")
	if err != nil {
		t.Skip("weasyprint is not installed")
	}
	engine, err := New(binary, 15*time.Second)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	pdf, err := engine.RenderPDF(context.Background(), []byte(`<!doctype html><meta charset="utf-8"><h1>Báo cáo quy hoạch Hà Nội</h1>`))
	if err != nil {
		t.Fatalf("RenderPDF() error = %v", err)
	}
	if !bytes.HasPrefix(pdf, []byte("%PDF-")) || len(pdf) < 1000 {
		t.Fatalf("PDF output invalid: %d bytes", len(pdf))
	}

	// pdftotext is test-only evidence that the engine preserves Vietnamese text;
	// production rendering itself has no dependency on this utility.
	pdfToText, err := exec.LookPath("pdftotext")
	if err != nil {
		return
	}
	temp, err := os.CreateTemp("", "qhpro-render-proof-*.pdf")
	if err != nil {
		t.Fatalf("CreateTemp() error = %v", err)
	}
	path := temp.Name()
	defer os.Remove(path)
	if _, err := temp.Write(pdf); err != nil {
		_ = temp.Close()
		t.Fatalf("write proof PDF: %v", err)
	}
	if err := temp.Close(); err != nil {
		t.Fatalf("close proof PDF: %v", err)
	}
	text, err := exec.Command(pdfToText, path, "-").Output()
	if err != nil {
		t.Fatalf("pdftotext proof error = %v", err)
	}
	if !strings.Contains(string(text), "Báo cáo quy hoạch Hà Nội") {
		t.Fatalf("rendered PDF lost Vietnamese content: %q", string(text))
	}
}

func TestEngineHonorsCallerCancellation(t *testing.T) {
	engine := &Engine{binary: "does-not-run", timeout: time.Second}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := engine.RenderPDF(ctx, []byte("<p>x</p>"))
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("RenderPDF() error = %v, want context.Canceled", err)
	}
}
