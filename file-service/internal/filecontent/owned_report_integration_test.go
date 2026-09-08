package filecontent

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"file/models"
	"file/repositories"
	"file/services"

	postgresdriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const reportLabOwnerNamespace = "g02_report_lab"

func newReportLabOwnedService(
	t *testing.T,
) (*Service, *gorm.DB, string) {
	t.Helper()

	if os.Getenv("QHPRO_REPORT_LAB_INTEGRATION") != "1" {
		t.Skip("report-lab integration test")
	}

	dsn := os.Getenv("QHPRO_TEST_FILE_DSN")
	if dsn == "" {
		t.Fatal("QHPRO_TEST_FILE_DSN is required")
	}

	db, err := gorm.Open(
		postgresdriver.Open(dsn),
		&gorm.Config{},
	)
	if err != nil {
		t.Fatalf("connect file PostgreSQL: %v", err)
	}

	if err := db.
		Unscoped().
		Where(
			"owner_namespace = ?",
			reportLabOwnerNamespace,
		).
		Delete(&models.FileEntity{}).
		Error; err != nil {
		t.Fatalf("clean owned file fixtures: %v", err)
	}

	root := t.TempDir()

	storage, err := services.NewStorageService(root)
	if err != nil {
		t.Fatalf("construct report-lab storage: %v", err)
	}

	service := NewService(
		repositories.NewFileRepository(db),
		storage,
		"report-lab-xor-key",
	)

	return service, db, root
}

func ownedInput(
	key string,
	content string,
) PutOwnedInput {
	return PutOwnedInput{
		Content:        []byte(content),
		Filename:       "qhpro-report-42.pdf",
		ContentType:    "application/pdf",
		Description:    "G02 owned File proof",
		OwnerNamespace: reportLabOwnerNamespace,
		OwnerKey:       key,
	}
}

func countOwnerRows(
	t *testing.T,
	db *gorm.DB,
	key string,
) int64 {
	t.Helper()

	var count int64

	if err := db.
		Model(&models.FileEntity{}).
		Where(
			"owner_namespace = ? AND owner_key = ?",
			reportLabOwnerNamespace,
			key,
		).
		Count(&count).
		Error; err != nil {
		t.Fatalf("count owner rows: %v", err)
	}

	return count
}

func physicalFiles(
	t *testing.T,
	root string,
) []string {
	t.Helper()

	var files []string

	err := filepath.Walk(
		root,
		func(
			path string,
			info os.FileInfo,
			err error,
		) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			if len(info.Name()) >= 7 &&
				info.Name()[:7] == ".owned-" {
				return nil
			}

			files = append(files, path)
			return nil
		},
	)
	if err != nil {
		t.Fatalf("walk owned files: %v", err)
	}

	return files
}

func TestReportLabOwnedFileSequentialRetryConverges(
	t *testing.T,
) {
	service, db, root :=
		newReportLabOwnedService(t)

	key := fmt.Sprintf(
		"sequential_%d",
		time.Now().UnixNano(),
	)

	input := ownedInput(
		key,
		"%PDF-1.7\nsame-owned-effect\n%%EOF\n",
	)

	first, err :=
		service.PutOwnedFile(
			context.Background(),
			input,
		)
	if err != nil {
		t.Fatalf("first PutOwnedFile(): %v", err)
	}

	second, err :=
		service.PutOwnedFile(
			context.Background(),
			input,
		)
	if err != nil {
		t.Fatalf("second PutOwnedFile(): %v", err)
	}

	if first.ID == 0 {
		t.Fatal("owned file ID is zero")
	}

	if first.ID != second.ID {
		t.Fatalf(
			"file IDs differ: first=%d second=%d",
			first.ID,
			second.ID,
		)
	}

	if first.Path != second.Path {
		t.Fatalf(
			"paths differ: first=%q second=%q",
			first.Path,
			second.Path,
		)
	}

	if got := countOwnerRows(t, db, key); got != 1 {
		t.Fatalf(
			"owner rows=%d, want 1",
			got,
		)
	}

	if got := len(physicalFiles(t, root)); got != 1 {
		t.Fatalf(
			"physical files=%d, want 1",
			got,
		)
	}
}

func TestReportLabOwnedFileConcurrentRetryConverges(
	t *testing.T,
) {
	service, db, root :=
		newReportLabOwnedService(t)

	key := fmt.Sprintf(
		"concurrent_%d",
		time.Now().UnixNano(),
	)

	input := ownedInput(
		key,
		"%PDF-1.7\nconcurrent-owned-effect\n%%EOF\n",
	)

	const callers = 8

	results := make(
		[]SavedFile,
		callers,
	)
	errs := make(
		[]error,
		callers,
	)

	var wg sync.WaitGroup

	for i := 0; i < callers; i++ {
		wg.Add(1)

		go func(index int) {
			defer wg.Done()

			results[index], errs[index] =
				service.PutOwnedFile(
					context.Background(),
					input,
				)
		}(i)
	}

	wg.Wait()

	var winnerID uint64
	var winnerPath string

	for i := 0; i < callers; i++ {
		if errs[i] != nil {
			t.Fatalf(
				"caller %d error: %v",
				i,
				errs[i],
			)
		}

		if results[i].ID == 0 {
			t.Fatalf(
				"caller %d returned zero File ID",
				i,
			)
		}

		if winnerID == 0 {
			winnerID = results[i].ID
			winnerPath = results[i].Path
		}

		if results[i].ID != winnerID {
			t.Fatalf(
				"caller %d File ID=%d, want %d",
				i,
				results[i].ID,
				winnerID,
			)
		}

		if results[i].Path != winnerPath {
			t.Fatalf(
				"caller %d path=%q, want %q",
				i,
				results[i].Path,
				winnerPath,
			)
		}
	}

	if got := countOwnerRows(t, db, key); got != 1 {
		t.Fatalf(
			"owner rows=%d, want 1",
			got,
		)
	}

	if got := len(physicalFiles(t, root)); got != 1 {
		t.Fatalf(
			"physical files=%d, want 1",
			got,
		)
	}
}

func TestReportLabOwnedFileRejectsDifferentRetry(
	t *testing.T,
) {
	service, db, root :=
		newReportLabOwnedService(t)

	key := fmt.Sprintf(
		"conflict_%d",
		time.Now().UnixNano(),
	)

	first := ownedInput(
		key,
		"%PDF-1.7\noriginal\n%%EOF\n",
	)

	if _, err :=
		service.PutOwnedFile(
			context.Background(),
			first,
		); err != nil {
		t.Fatalf(
			"first PutOwnedFile(): %v",
			err,
		)
	}

	changed := ownedInput(
		key,
		"%PDF-1.7\nDIFFERENT\n%%EOF\n",
	)

	_, err :=
		service.PutOwnedFile(
			context.Background(),
			changed,
		)

	if !errors.Is(
		err,
		ErrOwnedFileConflict,
	) {
		t.Fatalf(
			"conflict error=%v, want ErrOwnedFileConflict",
			err,
		)
	}

	if got := countOwnerRows(t, db, key); got != 1 {
		t.Fatalf(
			"owner rows=%d, want 1",
			got,
		)
	}

	files := physicalFiles(t, root)

	if len(files) != 1 {
		t.Fatalf(
			"physical files=%d, want 1",
			len(files),
		)
	}

	content, err := os.ReadFile(files[0])
	if err != nil {
		t.Fatalf(
			"read original owned file: %v",
			err,
		)
	}

	if string(content) != string(first.Content) {
		t.Fatalf(
			"original physical effect was changed: %q",
			string(content),
		)
	}
}
