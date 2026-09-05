package filehttp

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var (
	ErrConflict = errors.New(
		"owned file conflicts with existing effect",
	)

	ErrUnauthorized = errors.New(
		"file service authentication rejected",
	)
)

type Config struct {
	BaseURL        string
	PublicBaseURL  string
	ServiceName    string
	ServiceAuthKey string
	Timeout        time.Duration
}

type Client struct {
	baseURL        string
	publicBaseURL  string
	serviceName    string
	serviceAuthKey string
	http           *http.Client
}

func New(cfg Config) (*Client, error) {
	cfg.BaseURL = strings.TrimRight(
		strings.TrimSpace(cfg.BaseURL),
		"/",
	)
	cfg.PublicBaseURL = strings.TrimRight(
		strings.TrimSpace(cfg.PublicBaseURL),
		"/",
	)
	cfg.ServiceName = strings.TrimSpace(cfg.ServiceName)
	cfg.ServiceAuthKey = strings.TrimSpace(cfg.ServiceAuthKey)

	if cfg.BaseURL == "" {
		return nil, errors.New(
			"file HTTP base URL is required",
		)
	}

	parsed, err := url.Parse(cfg.BaseURL)
	if err != nil ||
		parsed.Scheme == "" ||
		parsed.Host == "" {
		return nil, errors.New(
			"file HTTP base URL is invalid",
		)
	}

	if cfg.PublicBaseURL == "" {
		cfg.PublicBaseURL = cfg.BaseURL
	}

	publicParsed, err := url.Parse(cfg.PublicBaseURL)
	if err != nil ||
		publicParsed.Scheme == "" ||
		publicParsed.Host == "" {
		return nil, errors.New(
			"file public base URL is invalid",
		)
	}

	if cfg.ServiceName == "" {
		return nil, errors.New(
			"file HTTP service name is required",
		)
	}

	if cfg.ServiceAuthKey == "" {
		return nil, errors.New(
			"file HTTP service auth key is required",
		)
	}

	if cfg.Timeout <= 0 {
		return nil, errors.New(
			"file HTTP timeout must be positive",
		)
	}

	return &Client{
		baseURL:        cfg.BaseURL,
		publicBaseURL:  cfg.PublicBaseURL,
		serviceName:    cfg.ServiceName,
		serviceAuthKey: cfg.ServiceAuthKey,
		http: &http.Client{
			Timeout: cfg.Timeout,
		},
	}, nil
}

type ownedFileResponse struct {
	FileID      uint64 `json:"fileId"`
	Path        string `json:"path"`
	Filename    string `json:"filename"`
	ContentType string `json:"contentType"`
	Size        int64  `json:"size"`
}

func (c *Client) PutOwnedFile(
	ctx context.Context,
	ownerNamespace string,
	ownerKey string,
	file io.Reader,
	filename string,
	contentType string,
) (string, error) {
	if ctx == nil {
		return "", errors.New("context is nil")
	}
	if file == nil {
		return "", errors.New("file is nil")
	}

	ownerNamespace =
		strings.TrimSpace(ownerNamespace)
	ownerKey =
		strings.TrimSpace(ownerKey)
	filename =
		filepath.Base(
			strings.TrimSpace(filename),
		)
	contentType =
		strings.TrimSpace(contentType)

	if ownerNamespace == "" ||
		ownerKey == "" ||
		filename == "" {
		return "", errors.New(
			"owned file identity and filename are required",
		)
	}

	reader, writer :=
		io.Pipe()

	multipartWriter :=
		multipart.NewWriter(writer)

	req, err :=
		http.NewRequestWithContext(
			ctx,
			http.MethodPost,
			c.baseURL+"/internal/v1/file/owned",
			reader,
		)
	if err != nil {
		_ = reader.Close()
		_ = writer.Close()
		return "", err
	}

	req.Header.Set(
		"Content-Type",
		multipartWriter.FormDataContentType(),
	)
	req.Header.Set(
		"X-Service-Name",
		c.serviceName,
	)
	req.Header.Set(
		"X-Service-Auth",
		signServiceRequest(
			time.Now(),
			c.serviceName,
			c.serviceAuthKey,
		),
	)

	writeDone :=
		make(chan error, 1)

	go func() {
		err :=
			writeMultipart(
				multipartWriter,
				ownerNamespace,
				ownerKey,
				file,
				filename,
				contentType,
			)

		if closeErr :=
			multipartWriter.Close(); err == nil {
			err = closeErr
		}

		_ = writer.CloseWithError(err)
		writeDone <- err
	}()

	resp, err :=
		c.http.Do(req)
	if err != nil {
		_ = reader.CloseWithError(err)
		<-writeDone

		return "", fmt.Errorf(
			"call File owned upload: %w",
			err,
		)
	}
	defer resp.Body.Close()

	if writeErr := <-writeDone; writeErr != nil {
		return "", fmt.Errorf(
			"write File multipart request: %w",
			writeErr,
		)
	}

	switch resp.StatusCode {
	case http.StatusOK:
		var result ownedFileResponse

		if err :=
			json.NewDecoder(resp.Body).
				Decode(&result); err != nil {
			return "", fmt.Errorf(
				"decode File owned upload response: %w",
				err,
			)
		}

		if strings.TrimSpace(result.Path) == "" {
			return "", errors.New(
				"File owned upload returned empty path",
			)
		}

		return result.Path, nil

	case http.StatusConflict:
		return "", ErrConflict

	case http.StatusUnauthorized:
		return "", ErrUnauthorized

	default:
		body, _ :=
			io.ReadAll(
				io.LimitReader(
					resp.Body,
					4096,
				),
			)

		return "", fmt.Errorf(
			"File owned upload failed: status=%d body=%q",
			resp.StatusCode,
			strings.TrimSpace(string(body)),
		)
	}
}

func (c *Client) PublicURL(path string) (string, error) {
	if c == nil {
		return "", errors.New("File HTTP client is nil")
	}

	path = strings.TrimSpace(path)
	if path == "" {
		return "", errors.New("File path is required")
	}

	return c.publicBaseURL +
		"/v1/file/load/" +
		url.PathEscape(path) +
		"/local", nil
}

func (c *Client) Close() {
	if c == nil || c.http == nil {
		return
	}

	c.http.CloseIdleConnections()
}

func writeMultipart(
	writer *multipart.Writer,
	ownerNamespace string,
	ownerKey string,
	file io.Reader,
	filename string,
	contentType string,
) error {
	if err :=
		writer.WriteField(
			"owner_namespace",
			ownerNamespace,
		); err != nil {
		return err
	}

	if err :=
		writer.WriteField(
			"owner_key",
			ownerKey,
		); err != nil {
		return err
	}

	if err :=
		writer.WriteField(
			"content_type",
			contentType,
		); err != nil {
		return err
	}

	part, err :=
		writer.CreateFormFile(
			"file",
			filename,
		)
	if err != nil {
		return err
	}

	_, err =
		io.Copy(part, file)

	return err
}

func signServiceRequest(
	now time.Time,
	serviceName string,
	secret string,
) string {
	timestamp :=
		strconv.FormatInt(
			now.Unix(),
			10,
		)

	mac :=
		hmac.New(
			sha256.New,
			[]byte(secret),
		)

	_, _ =
		mac.Write(
			[]byte(
				timestamp +
					":" +
					serviceName,
			),
		)

	signature :=
		base64.RawURLEncoding.
			EncodeToString(
				mac.Sum(nil),
			)

	return timestamp + ":" + signature
}
