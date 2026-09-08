package filecontent

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	_utils "common/utils"
	"file/models"
)

var (
	ErrInvalidOwnedFile = errors.New(
		"owned file input is invalid",
	)
	ErrOwnedFileConflict = errors.New(
		"owned file key was already used with different input",
	)
)

var (
	ownerNamespacePattern = regexp.MustCompile(`^[a-z0-9._-]+$`)

	ownerKeyPattern = regexp.MustCompile(`^[A-Za-z0-9._:-]+$`)
)

type PutOwnedInput struct {
	Content     []byte
	Filename    string
	ContentType string
	Description string

	OwnerNamespace string
	OwnerKey       string
}

func (s *Service) PutOwnedFile(
	ctx context.Context,
	input PutOwnedInput,
) (SavedFile, error) {
	if ctx == nil {
		return SavedFile{}, fmt.Errorf(
			"%w: context is nil",
			ErrInvalidOwnedFile,
		)
	}

	if len(input.Content) == 0 {
		return SavedFile{}, fmt.Errorf(
			"%w: content is empty",
			ErrInvalidOwnedFile,
		)
	}
	if len(input.Content) > MaxContentBytes {
		return SavedFile{}, fmt.Errorf(
			"%w: content exceeds 64 MiB",
			ErrInvalidOwnedFile,
		)
	}

	input.Filename =
		strings.TrimSpace(input.Filename)

	input.ContentType =
		strings.TrimSpace(input.ContentType)

	input.OwnerNamespace =
		strings.TrimSpace(input.OwnerNamespace)

	input.OwnerKey =
		strings.TrimSpace(input.OwnerKey)

	if input.Filename == "" {
		return SavedFile{}, fmt.Errorf(
			"%w: filename is empty",
			ErrInvalidOwnedFile,
		)
	}

	if len(input.OwnerNamespace) > 64 ||
		!ownerNamespacePattern.MatchString(
			input.OwnerNamespace,
		) {
		return SavedFile{}, fmt.Errorf(
			"%w: owner namespace is invalid",
			ErrInvalidOwnedFile,
		)
	}

	if len(input.OwnerKey) > 128 ||
		!ownerKeyPattern.MatchString(
			input.OwnerKey,
		) {
		return SavedFile{}, fmt.Errorf(
			"%w: owner key is invalid",
			ErrInvalidOwnedFile,
		)
	}

	requestHash :=
		ownedRequestHash(input)

	namespace := input.OwnerNamespace
	ownerKey := input.OwnerKey

	candidate := &models.FileEntity{
		Description:      input.Description,
		OwnerNamespace:   &namespace,
		OwnerKey:         &ownerKey,
		OwnerRequestHash: &requestHash,
	}

	created, err :=
		s.files.TryCreateOwned(
			ctx,
			candidate,
		)
	if err != nil {
		return SavedFile{}, fmt.Errorf(
			"create owned file authority: %w",
			err,
		)
	}

	fileEntity := candidate

	if !created {
		existing, found, err :=
			s.files.FindByOwner(
				ctx,
				namespace,
				ownerKey,
			)
		if err != nil {
			return SavedFile{}, fmt.Errorf(
				"load owned file authority: %w",
				err,
			)
		}
		if !found {
			return SavedFile{}, errors.New(
				"owned file authority disappeared after conflict",
			)
		}

		if existing.OwnerRequestHash == nil ||
			*existing.OwnerRequestHash != requestHash {
			return SavedFile{},
				ErrOwnedFileConflict
		}

		if existing.Path != "" &&
			existing.Hash != "" {
			return savedOwnedFile(existing), nil
		}

		fileEntity = existing
	}

	fileInfo, err :=
		s.storage.StoreOwnedContent(
			bytes.NewReader(input.Content),
			input.Filename,
			input.ContentType,
			int64(len(input.Content)),
			nil,
			fileEntity.ID,
			ordinaryDocumentPath,
			"owned",
			namespace,
		)
	if err != nil {
		// Keep the owner row. A later retry must adopt the same File ID
		// rather than creating a second logical effect.
		return SavedFile{}, fmt.Errorf(
			"store owned file content: %w",
			err,
		)
	}

	encryptedPath :=
		_utils.XorEncode(
			fileInfo.RelativePath,
			s.xorKey,
		)

	fileEntity.SetFields(fileInfo)
	fileEntity.Path = "p" + encryptedPath

	if err :=
		s.files.SaveWithContext(
			ctx,
			fileEntity,
		); err != nil {
		// Do not delete the stable file or owner row here.
		// Retry is allowed to converge on the same identity/path.
		return SavedFile{}, fmt.Errorf(
			"save owned file metadata: %w",
			err,
		)
	}

	return savedOwnedFile(fileEntity), nil
}

func savedOwnedFile(
	file *models.FileEntity,
) SavedFile {
	if file == nil {
		return SavedFile{}
	}

	return SavedFile{
		ID:            file.ID,
		Path:          file.Path,
		Filename:      file.Name,
		ThumbnailPath: file.ThumbnailPath,
		Extension:     file.Extension,
		Size:          file.Size,
		Hash:          file.Hash,
		ContentType:   file.Mine,
	}
}

func ownedRequestHash(
	input PutOwnedInput,
) string {
	contentSum :=
		sha256.Sum256(input.Content)

	fingerprint := struct {
		Filename      string `json:"filename"`
		ContentType   string `json:"contentType"`
		Description   string `json:"description"`
		ContentSHA256 string `json:"contentSha256"`
	}{
		Filename:    input.Filename,
		ContentType: input.ContentType,
		Description: input.Description,
		ContentSHA256: hex.EncodeToString(
			contentSum[:],
		),
	}

	encoded, _ := json.Marshal(fingerprint)
	sum := sha256.Sum256(encoded)

	return hex.EncodeToString(sum[:])
}
