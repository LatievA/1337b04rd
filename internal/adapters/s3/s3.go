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

func (c *HTTPClient) CreateObject(bucketName, objectKey, contentType string, data []byte) (string, error) {
	url := c.baseURL + "/" + bucketName + "/" + objectKey
	slog.Info("Creating object", "url", url, "contentType", contentType)

	req, err := http.NewRequest("PUT", url, strings.NewReader(string(data)))
	if err != nil {
		slog.Error("Failed to create object request", "err", err)
		return "", err
	}
	req.Header.Set("Content-Type", contentType)

	resp, err := c.client.Do(req)
	if err != nil {
		slog.Error("Failed to execute object creation request", "err", err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		slog.Error("Failed to create object", "statusCode", resp.StatusCode)
		return "", err
	}

	return url, nil
}

func (c *HTTPClient) GetObject(bucketName, objectKey string) (string, error) {
	url := c.baseURL + "/" + bucketName + "/" + objectKey
	slog.Info("Getting object", "url", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		slog.Error("Failed to create get object request", "err", err)
		return "", err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		slog.Error("Failed to execute get object request", "err", err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Error("Failed to get object", "statusCode", resp.StatusCode)
		return "", err
	}

	data := ""
	if _, err := resp.Body.Read([]byte(data)); err != nil {
		slog.Error("Failed to read object data", "err", err)
		return "", err
	}

	return data, nil
}