package storage

import (
	"log"
	"net/http"
	"strings"
	"time"
	"triple-s/info"
	"triple-s/utils"
)

func CreateBucket(w http.ResponseWriter, r *http.Request) {
	bucketName := strings.TrimPrefix(r.URL.Path, "/")
	if err := validateBucketName(bucketName); err != nil {
		log.Printf("Error validating bucket name %s: %v\n", bucketName, err)
		ErrXMLResponse(w, http.StatusBadRequest, "Invalid bucket name")
		return
	}

	bucketsData, err := utils.ReadBucket()
	if err != nil {
		log.Printf("Error reading bucket file: %v\n", err)
		ErrXMLResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	if SearchBucketIDX(bucketsData.Buckets, bucketName) != -1 {
		log.Printf("Bucket %s already exists\n", bucketName)
		ErrXMLResponse(w, http.StatusConflict, "Bucket already exists")
		return
	}

	newBucket := info.Bucket{
		Name:             bucketName,
		CreationTime:     time.Now().Format(time.RFC3339Nano),
		LastModifiedTime: time.Now().Format(time.RFC3339Nano),
		Status:           "Available",
	}
	bucketsData.Buckets = append(bucketsData.Buckets, newBucket)

	if err := utils.WriteBucket(bucketsData); err != nil {
		log.Printf("Error writing buckets file: %v\n", err)
		ErrXMLResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	if err := CreateBucketDir(bucketName); err != nil {
		log.Printf("Error creating bucket directory: %v\n", err)
		ErrXMLResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	if err := InitObjectFile(bucketName); err != nil {
		log.Printf("Error initializing object file: %v\n", err)
		ErrXMLResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	log.Printf("Bucket %s created successfully", bucketName)
	WriteXMLResponse(w, http.StatusOK, "Bucket created successfully")
}

func ListBuckets(w http.ResponseWriter, r *http.Request) {
	bucketData, err := utils.ReadBucket()
	if err != nil {
		log.Printf("Error reading bucket file: %v\n", err)
		ErrXMLResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	log.Printf("Buckets read successfully")
	WriteXMLResponse(w, http.StatusOK, bucketData)
}

func DeleteBucket(w http.ResponseWriter, r *http.Request) {
	bucketName := strings.TrimPrefix(r.URL.Path, "/")

	bucketsData, err := utils.ReadBucket()
	if err != nil {
		log.Printf("Error reading bucket file: %v\n", err)
		ErrXMLResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	bucketIDX := SearchBucketIDX(bucketsData.Buckets, bucketName)
	if bucketIDX == -1 {
		log.Printf("Bucket doesn't exist)")
		ErrXMLResponse(w, http.StatusNotFound, "Bucket not found")
		return
	}

	if !isBucketEmpty(bucketName) {
		log.Printf("Bucket %s is not empty\n", bucketName)
		ErrXMLResponse(w, http.StatusConflict, "Bucket is not empty")
		return
	}

	if bucketsData.Buckets[bucketIDX].Status == "Marked for delete" {
		if err := RemoveBucket(bucketName); err != nil {
			log.Printf("Error removing bucket directory: %v\n", err)
			ErrXMLResponse(w, http.StatusInternalServerError, "Internal server error")
			return
		}

		bucketsData.Buckets = append(bucketsData.Buckets[:bucketIDX], bucketsData.Buckets[bucketIDX+1:]...)

		if err := utils.WriteBucket(bucketsData); err != nil {
			log.Printf("Error writing buckets file: %v\n", err)
			ErrXMLResponse(w, http.StatusInternalServerError, "Internal server error")
			return
		}

		log.Printf("Bucket %s permanently deleted", bucketName)
		w.WriteHeader(http.StatusNoContent)
		return
	}

	bucketsData.Buckets[bucketIDX].Status = "Marked for delete"
	bucketsData.Buckets[bucketIDX].LastModifiedTime = time.Now().Format(time.RFC3339Nano)

	if err := utils.WriteBucket(bucketsData); err != nil {
		log.Printf("Error writing buckets file: %v\n", err)
		ErrXMLResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	log.Printf("Bucket %s marked for delete (soft delete)", bucketName)
	WriteXMLResponse(w, http.StatusOK, "Bucket marked for delete")
}
