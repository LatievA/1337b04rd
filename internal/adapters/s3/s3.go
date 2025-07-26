package s3

import (
	"log/slog"
	"net/http"
	"strings"
	"time"
)

type HTTPClient struct {
	client  *http.Client
	baseURL string
}

func NewHTTPClient(baseURL string) *HTTPClient {
	return &HTTPClient{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: baseURL,
	}
}

func (c *HTTPClient) CreateBucket(bucketName string) error {
	url := c.baseURL + "/" + bucketName
	slog.Info("Creating bucket", "url", url)

	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		slog.Error("Failed to create bucket request", "err", err)
		return err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		slog.Error("Failed to execute bucket creation request", "err", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		slog.Error("Failed to create bucket", "statusCode", resp.StatusCode)
		return err
	}

	return nil
}

func (c *HTTPClient) CreateObject(bucketName, objectKey, contentType string, data string) error {
	url := c.baseURL + "/" + bucketName + "/" + objectKey
	slog.Info("Creating object", "url", url, "contentType", contentType)

	req, err := http.NewRequest("PUT", url, strings.NewReader(data))
	if err != nil {
		slog.Error("Failed to create object request", "err", err)
		return err
	}
	req.Header.Set("Content-Type", contentType)

	resp, err := c.client.Do(req)
	if err != nil {
		slog.Error("Failed to execute object creation request", "err", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		slog.Error("Failed to create object", "statusCode", resp.StatusCode)
		return err
	}

	return nil
}