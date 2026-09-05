package filemedia

import "testing"

func TestParsePathPreservesCurrentVideoLayout(t *testing.T) {
	info, err := ParsePath(
		"files/public/video/2026-08-25/1752485034177081000/619.mp4",
	)
	if err != nil {
		t.Fatalf("ParsePath() error = %v", err)
	}

	if info.Scope != "public" {
		t.Fatalf("Scope = %q, want public", info.Scope)
	}
	if info.Type != "video" {
		t.Fatalf("Type = %q, want video", info.Type)
	}
	if info.Date != "2026-08-25" {
		t.Fatalf("Date = %q", info.Date)
	}
	if info.UUID != "1752485034177081000" {
		t.Fatalf("UUID = %q", info.UUID)
	}
	if info.FileID != "619" {
		t.Fatalf("FileID = %q, want 619", info.FileID)
	}
	if info.Ext != "mp4" {
		t.Fatalf("Ext = %q, want mp4", info.Ext)
	}
}

func TestParsePathRejectsMalformedLayout(t *testing.T) {
	if _, err := ParsePath("public/video/42"); err == nil {
		t.Fatal("ParsePath() error = nil, want invalid path")
	}
}
