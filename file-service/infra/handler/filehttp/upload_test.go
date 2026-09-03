package filehttp

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"file/dto"
	"file/internal/filecontent"
	"file/internal/serviceauth"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeFileOperations struct {
	uploadInput   filecontent.UploadInput
	uploadContent []byte
	saved         filecontent.SavedFile
	err           error
}

func (f *fakeFileOperations) UploadFile(
	_ context.Context,
	input filecontent.UploadInput,
) (filecontent.SavedFile, error) {
	content, err := io.ReadAll(input.Content)
	if err != nil {
		return filecontent.SavedFile{}, err
	}

	input.Content = nil
	f.uploadInput = input
	f.uploadContent = content

	return f.saved, f.err
}

func (f *fakeFileOperations) PutOwnedFile(
	context.Context,
	filecontent.PutOwnedInput,
) (filecontent.SavedFile, error) {
	return filecontent.SavedFile{}, nil
}

func TestUploadFilePreservesPublicHTTPContract(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	require.NoError(t, writer.WriteField("secure", "false"))

	// Swagger says formData, but the existing runtime contract reads
	// description from query. Preserve runtime behavior during cutover.
	require.NoError(
		t,
		writer.WriteField("description", "form-description"),
	)

	part, err := writer.CreateFormFile("file", "hello.txt")
	require.NoError(t, err)

	_, err = part.Write([]byte("hello"))
	require.NoError(t, err)

	require.NoError(t, writer.Close())

	files := &fakeFileOperations{
		saved: filecontent.SavedFile{
			ID:          7,
			Path:        "p-token",
			Filename:    "hello.txt",
			Extension:   ".txt",
			Size:        5,
			Hash:        "hash-1",
			ContentType: "application/octet-stream",
		},
	}

	handler := NewHandler(
		files,
		serviceauth.Verifier{},
	)

	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(
		http.MethodPost,
		"/v1/file/upload?description=query-description",
		&body,
	)
	c.Request.Header.Set(
		"Content-Type",
		writer.FormDataContentType(),
	)

	handler.UploadFile(c)

	require.Equal(
		t,
		http.StatusOK,
		recorder.Code,
		recorder.Body.String(),
	)

	assert.Equal(
		t,
		"hello",
		string(files.uploadContent),
	)
	assert.Equal(
		t,
		"query-description",
		files.uploadInput.Description,
	)
	assert.False(t, files.uploadInput.Secure)
	assert.Equal(t, int64(5), files.uploadInput.Size)
	assert.Equal(t, "hello.txt", files.uploadInput.Filename)

	var response struct {
		File dto.FileInfo `json:"file"`
	}

	require.NoError(
		t,
		json.Unmarshal(recorder.Body.Bytes(), &response),
	)

	assert.Equal(t, uint64(7), response.File.ID)
	assert.Equal(t, "hello.txt", response.File.FileName)
	assert.Equal(t, "p-token", response.File.RelativePath)
	// Physical deployment paths are an internal storage detail and must never
	// be projected through the public File HTTP contract.
	assert.Empty(t, response.File.AbsolutePath)
	assert.Equal(t, ".txt", response.File.Extension)
	assert.Equal(t, int64(5), response.File.Size)
	assert.Equal(t, "hash-1", response.File.Hash)
	assert.Equal(
		t,
		"application/octet-stream",
		response.File.ContentType,
	)
}
