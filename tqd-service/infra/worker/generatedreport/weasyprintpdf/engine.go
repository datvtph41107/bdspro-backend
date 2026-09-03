package weasyprintpdf

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// Engine renders the self-contained Generated Report HTML using a process-owned
// WeasyPrint executable. The timeout bounds one render independently of process
// shutdown, while CommandContext still propagates caller cancellation.
type Engine struct {
	binary  string
	timeout time.Duration
}

func New(binary string, timeout time.Duration) (*Engine, error) {
	binary = strings.TrimSpace(binary)
	if binary == "" {
		return nil, errors.New("WeasyPrint binary is required")
	}
	if timeout <= 0 {
		return nil, errors.New("WeasyPrint render timeout must be positive")
	}
	resolved, err := exec.LookPath(binary)
	if err != nil {
		return nil, fmt.Errorf("find WeasyPrint binary: %w", err)
	}
	return &Engine{binary: resolved, timeout: timeout}, nil
}

func (e *Engine) RenderPDF(ctx context.Context, html []byte) ([]byte, error) {
	if e == nil || e.binary == "" || e.timeout <= 0 {
		return nil, errors.New("WeasyPrint PDF engine is not configured")
	}
	if ctx == nil {
		return nil, errors.New("WeasyPrint PDF context is nil")
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if len(bytes.TrimSpace(html)) == 0 {
		return nil, errors.New("WeasyPrint PDF HTML is empty")
	}

	directory, err := os.MkdirTemp("", "qhpro-report-render-*")
	if err != nil {
		return nil, fmt.Errorf("create WeasyPrint render directory: %w", err)
	}
	defer os.RemoveAll(directory)

	inputPath := filepath.Join(directory, "report.html")
	outputPath := filepath.Join(directory, "report.pdf")
	if err := os.WriteFile(inputPath, html, 0o600); err != nil {
		return nil, fmt.Errorf("write WeasyPrint report HTML: %w", err)
	}

	renderCtx, cancel := context.WithTimeout(ctx, e.timeout)
	defer cancel()

	command := exec.CommandContext(renderCtx, e.binary, inputPath, outputPath)
	var stderr bytes.Buffer
	command.Stderr = &stderr
	if err := command.Run(); err != nil {
		if ctxErr := renderCtx.Err(); ctxErr != nil {
			return nil, ctxErr
		}
		return nil, fmt.Errorf("WeasyPrint PDF render failed: %w: %s", err, strings.TrimSpace(stderr.String()))
	}

	pdf, err := os.ReadFile(outputPath)
	if err != nil {
		return nil, fmt.Errorf("read WeasyPrint PDF output: %w", err)
	}
	if len(pdf) < 8 || !bytes.HasPrefix(pdf, []byte("%PDF-")) {
		return nil, errors.New("WeasyPrint produced invalid PDF output")
	}
	return pdf, nil
}
