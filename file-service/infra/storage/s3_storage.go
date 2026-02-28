package storage

// import (
// 	"context"
// 	"fmt"
// 	"io"
// 	"os"
// 	"time"

// 	"github.com/aws/aws-sdk-go/aws"
// 	"github.com/aws/aws-sdk-go/aws/credentials"
// 	"github.com/aws/aws-sdk-go/aws/session"
// 	"github.com/aws/aws-sdk-go/service/s3"
// 	"github.com/aws/aws-sdk-go/service/s3/s3manager"
// )

// type S3Storage struct {
// 	bucket   string
// 	region   string
// 	endpoint string
// 	s3Client *s3.S3
// 	uploader *s3manager.Uploader
// 	baseURL  string
// }

// func NewS3Storage(bucket, region, endpoint string) *S3Storage {
// 	config := &aws.Config{
// 		Region: aws.String(region),
// 	}

// 	if endpoint != "" {
// 		config.Endpoint = aws.String(endpoint)
// 		config.S3ForcePathStyle = aws.Bool(true)

// 		// For MinIO or other S3-compatible services
// 		config.Credentials = credentials.NewStaticCredentials(
// 			os.Getenv("AWS_ACCESS_KEY_ID"),
// 			os.Getenv("AWS_SECRET_ACCESS_KEY"),
// 			"",
// 		)
// 	}

// 	sess := session.Must(session.NewSession(config))
// 	s3Client := s3.New(sess)
// 	uploader := s3manager.NewUploader(sess)

// 	return &S3Storage{
// 		bucket:   bucket,
// 		region:   region,
// 		endpoint: endpoint,
// 		s3Client: s3Client,
// 		uploader: uploader,
// 	}
// }

// func (s *S3Storage) WithBaseURL(baseURL string) *S3Storage {
// 	s.baseURL = baseURL
// 	return s
// }

// func (s *S3Storage) Store(ctx context.Context, key string, reader io.Reader, size int64, metadata map[string]string) error {
// 	// Convert metadata to S3 format
// 	s3Metadata := make(map[string]*string)
// 	for k, v := range metadata {
// 		s3Metadata[k] = aws.String(v)
// 	}

// 	// Upload to S3
// 	_, err := s.uploader.UploadWithContext(ctx, &s3manager.UploadInput{
// 		Bucket:      aws.String(s.bucket),
// 		Key:         aws.String(key),
// 		Body:        reader,
// 		ContentType: aws.String(metadata["Content-Type"]),
// 		Metadata:    s3Metadata,
// 	})

// 	return err
// }

// func (s *S3Storage) Retrieve(ctx context.Context, key string) (io.ReadCloser, int64, error) {
// 	result, err := s.s3Client.GetObjectWithContext(ctx, &s3.GetObjectInput{
// 		Bucket: aws.String(s.bucket),
// 		Key:    aws.String(key),
// 	})
// 	if err != nil {
// 		return nil, 0, err
// 	}

// 	return result.Body, *result.ContentLength, nil
// }

// func (s *S3Storage) Delete(ctx context.Context, key string) error {
// 	_, err := s.s3Client.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
// 		Bucket: aws.String(s.bucket),
// 		Key:    aws.String(key),
// 	})
// 	return err
// }

// func (s *S3Storage) Copy(ctx context.Context, srcKey, dstKey string) error {
// 	_, err := s.s3Client.CopyObjectWithContext(ctx, &s3.CopyObjectInput{
// 		Bucket:     aws.String(s.bucket),
// 		CopySource: aws.String(fmt.Sprintf("%s/%s", s.bucket, srcKey)),
// 		Key:        aws.String(dstKey),
// 	})
// 	return err
// }

// func (s *S3Storage) Exists(ctx context.Context, key string) (bool, error) {
// 	_, err := s.s3Client.HeadObjectWithContext(ctx, &s3.HeadObjectInput{
// 		Bucket: aws.String(s.bucket),
// 		Key:    aws.String(key),
// 	})
// 	if err != nil {
// 		if isNotFoundError(err) {
// 			return false, nil
// 		}
// 		return false, err
// 	}
// 	return true, nil
// }

// func (s *S3Storage) GetPresignedURL(ctx context.Context, key string, expiry time.Duration) (string, error) {
// 	req, _ := s.s3Client.GetObjectRequest(&s3.GetObjectInput{
// 		Bucket: aws.String(s.bucket),
// 		Key:    aws.String(key),
// 	})

// 	url, err := req.Presign(expiry)
// 	if err != nil {
// 		return "", err
// 	}

// 	return url, nil
// }

// func isNotFoundError(err error) bool {
// 	// Check if error is "NoSuchKey" or "NotFound"
// 	return err != nil && (err.Error() == "NoSuchKey" ||
// 		err.Error() == "NotFound" ||
// 		err.Error() == "AccessDenied")
// }

// func (s *S3Storage) List(ctx context.Context, prefix string) ([]string, error) {
// 	var objects []string

// 	err := s.s3Client.ListObjectsV2PagesWithContext(ctx, &s3.ListObjectsV2Input{
// 		Bucket: aws.String(s.bucket),
// 		Prefix: aws.String(prefix),
// 	}, func(page *s3.ListObjectsV2Output, lastPage bool) bool {
// 		for _, obj := range page.Contents {
// 			objects = append(objects, *obj.Key)
// 		}
// 		return !lastPage
// 	})

// 	return objects, err
// }
