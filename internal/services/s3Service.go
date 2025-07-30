package services

import (
	"1337b04rd/internal/adapters/s3"
	"1337b04rd/internal/domain"
	"context"
	"fmt"
	"mime"
	"path/filepath"
)

// S3ServiceImpl implements domain.S3Service using the HTTPClient
type S3ServiceImpl struct {
	client *s3.HTTPClient
}

// NewS3Service creates a new S3ServiceImpl with given base URL and bucket name
func NewS3Service(baseURL string) domain.S3Service {
	httpClient := s3.NewHTTPClient(baseURL)
	return &S3ServiceImpl{client: httpClient}
}

// UploadImage reads raw bytes and uploads to S3, returning the public URL
func (s *S3ServiceImpl) UploadImage(ctx context.Context, fileData []byte, bucketName, objectKey string) (string, error) {
	// Detect content type from filename extension or data
	contentType := mime.TypeByExtension(
		"." + filepath.Ext(objectKey)[1:],
	)
	if contentType == "" {
		// fallback to generic binary
		contentType = "application/octet-stream"
	}

	// Use the HTTPClient to do the upload
	url, err := s.client.CreateObject(bucketName, objectKey, contentType, fileData)
	if err != nil {
		return "", fmt.Errorf("failed to upload image to S3: %w", err)
	}
	return url, nil
}
