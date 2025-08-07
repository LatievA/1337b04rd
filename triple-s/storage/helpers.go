package storage

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"triple-s/flags"
	"triple-s/info"
	"triple-s/utils"
)

func validateBucketName(name string) error {
	if len(name) < 3 || len(name) > 63 {
		return fmt.Errorf("bucket name must be between 3 and 63 characters")
	}

	for i := 0; i < len(name); i++ {
		if (name[i] < 'a' || name[i] > 'z') && name[i] != '.' && name[i] != '-' && (name[i] < '0' || name[i] > '9') {
			return fmt.Errorf("only lowercase letters, numbers, hyphens (-), and dots (.) are allowed")
		}
	}

	if strings.Contains(name, "--") {
		return fmt.Errorf("bucket name cannot contain consecutive hyphens")
	}

	if strings.Contains(name, "..") {
		return fmt.Errorf("bucket name cannot contain consecutive periods")
	}

	if strings.HasPrefix(name, "-") || strings.HasSuffix(name, "-") {
		return fmt.Errorf("bucket name cannot begin or end with a hyphen")
	}

	ipRegex := `^\d+(\.\d+){3}$`
	if match, _ := regexp.MatchString(ipRegex, name); match {
		return fmt.Errorf("bucket name must not be an IP address")
	}

	return nil
}

func SearchBucketIDX(buckets []info.Bucket, name string) int {
	for i, bucket := range buckets {
		if bucket.Name == name {
			return i
		}
	}
	return -1
}

func SearchObjectIDX(objects []info.Object, key string) int {
	for i, object := range objects {
		if object.ObjectKey == key {
			return i
		}
	}
	return -1
}

func CreateBucketDir(bucketName string) error {
	dirPath := filepath.Join(flags.Dir, bucketName)
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create bucket directory: %w", err)
	}
	return nil
}

func isBucketEmpty(bucketName string) bool {
	bucketPath := GetBucketPath(bucketName)
	objectsMetaDataPath := filepath.Join(bucketPath, "objects.csv")
	if _, err := os.Stat(objectsMetaDataPath); os.IsNotExist(err) {
		return true
	}

	objects, err := utils.ReadCSV(objectsMetaDataPath)
	if err != nil {
		fmt.Printf("Error reading objects metadata: %v\n", err)
		return false
	}

	return len(objects) == 0
}

func RemoveBucket(bucketname string) error {
	dirPath := filepath.Join(flags.Dir, bucketname)
	if err := os.RemoveAll(dirPath); err != nil {
		return fmt.Errorf("error during removing diretory: %w", err)
	}
	return nil
}

func GetBucketPath(bucketName string) string {
	return filepath.Join(flags.Dir, bucketName)
}

func ValidatePath(path string) (bucketName, objectKey string) {
	path = strings.TrimPrefix(path, "/")
	parts := strings.SplitN(path, "/", 2)
	bucketName = parts[0]
	if len(parts) > 1 {
		objectKey = parts[1]
	}

	return
}

func GetMetaDataType(bucketPath, objectPath string) (string, error) {
	metaDataPath := filepath.Join(bucketPath, "objects.csv")
	records, err := utils.ReadCSV(metaDataPath)
	if err != nil {
		return "", fmt.Errorf("failed to read metadata: %w", err)
	}

	for _, record := range records {
		if record[0] == objectPath {
			return record[1], nil
		}
	}
	return "", fmt.Errorf("metadata type not found")
}

func WriteXMLResponse(w http.ResponseWriter, statusCode int, v interface{}) {
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(statusCode)

	if err := xml.NewEncoder(w).Encode(v); err != nil {
		log.Printf("Error enocoding response: %v", err)
	}
}

func ErrXMLResponse(w http.ResponseWriter, code int, message string) {
	WriteXMLResponse(w, code, info.ErrResp{Code: code, Message: message})
}
