package mapper

import (
	"testing"
	"time"

	bdspropb "pb/types/bdspro"
)

func TestPostSaveWithProductToListing_ExpiredFromNumDate(t *testing.T) {
	post := &bdspropb.PostSaveWithProduct{
		Title:           "Tin A",
		Content:         "Nội dung",
		PostPrice:       1.5e9,
		TransactionType: 10,
		Visibility:      30,
		NumDate:         7,
	}
	got := PostSaveWithProductToListing(post, 99)
	if got == nil {
		t.Fatal("nil listing")
	}
	if got.OwnerID != 99 {
		t.Fatalf("owner: %d", got.OwnerID)
	}
	if got.NumDate != 7 {
		t.Fatalf("numDate: %d", got.NumDate)
	}
	if got.ExpiredAt == nil {
		t.Fatal("expected expired_at")
	}
	if got.ExpiredAt.Before(time.Now().Add(6 * 24 * time.Hour)) {
		t.Fatalf("expired too soon: %v", got.ExpiredAt)
	}
}

func TestPostSaveWithProductToListing_ExpiredAtUnixSeconds(t *testing.T) {
	exp := time.Date(2026, 12, 1, 0, 0, 0, 0, time.UTC)
	post := &bdspropb.PostSaveWithProduct{
		Title:     "T",
		ExpiredAt: int32(exp.Unix()),
	}
	got := PostSaveWithProductToListing(post, 1)
	if got.ExpiredAt == nil || !got.ExpiredAt.Equal(exp) {
		t.Fatalf("expired: got %v want %v", got.ExpiredAt, exp)
	}
}
