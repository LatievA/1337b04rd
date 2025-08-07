package storage

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"triple-s/info"
	"triple-s/utils"
)

func CreateObject(w http.ResponseWriter, r *http.Request) {
	bucketName, objectKey := ValidatePath(r.URL.Path)
	if strings.TrimSpace(bucketName) == "" || strings.TrimSpace(objectKey) == "" {
		log.Println("Invalid  bucket or object key")
		ErrXMLResponse(w, http.StatusBadRequest, "Invalid bucket or object key")
		return
	}

	bucketPath := GetBucketPath(bucketName)
	if _, err := os.Stat(bucketPath); os.IsNotExist(err) {
		log.Printf("Bucket not found: %s\n", bucketName)
		ErrXMLResponse(w, http.StatusNotFound, "Bucket not found")
		return
	}

	objects, err := utils.ReadObjects(bucketPath)
	fmt.Println(objects)
	if err != nil {
		log.Printf("Failed to read objects file for bucket %s: %v\n", bucketName, err)
		ErrXMLResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	objectIDX := SearchObjectIDX(objects.Objects, objectKey)
	if objectIDX != -1 {
		objects.Objects = append(objects.Objects[:objectIDX], objects.Objects[objectIDX+1:]...)
	}

	objectPath := filepath.Join(bucketPath, objectKey)
	file, err := os.Create(objectPath)
	if err != nil {
		log.Printf("Failed to create object %s in bucket %s: %v\n", objectKey, bucketName, err)
		ErrXMLResponse(w, http.StatusInternalServerError, "Failed to create object")
		return
	}
	defer file.Close()

	_, err = io.Copy(file, r.Body)
	if err != nil {
		log.Printf("Failed to write object data for %s in bucket %s: %v\n", objectKey, bucketName, err)
		ErrXMLResponse(w, http.StatusInternalServerError, "Failed to write object data")
		return
	}

	newObject := info.Object{
		ObjectKey:    objectKey,
		ContentType:  r.Header.Get("Content-Type"),
		Size:         strconv.FormatInt(r.ContentLength, 10),
		LastModified: time.Now().Format(time.RFC3339Nano),
	}

	if newObject.ContentType == "" {
		newObject.ContentType = "application/octet-stream"
	}

	objects.Objects = append(objects.Objects, newObject)
	fmt.Println(objects)
	err = utils.WriteObject(bucketName, objects)
	if err != nil {
		log.Printf("Failed to update objects file for bucket %s: %v\n", bucketName, err)
		ErrXMLResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	log.Printf("Object %s created successfully in bucket %s", objectKey, bucketName)
	ErrXMLResponse(w, http.StatusOK, "Object created successfully")
}

func ListObjects(w http.ResponseWriter, r *http.Request) {
	bucketName := strings.TrimPrefix(r.URL.Path, "/")
	bucketPath := GetBucketPath(bucketName)
	objectsData, err := utils.ReadObjects(bucketPath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("Bucket not found: %s\n", bucketName)
			ErrXMLResponse(w, http.StatusNotFound, "Bucket not found")
		} else {
			log.Printf("Error reading objects file for bucket %s: %v\n", bucketName, err)
			ErrXMLResponse(w, http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	log.Printf("Objects listed successfully for bucket %s", bucketName)
	WriteXMLResponse(w, http.StatusOK, objectsData)
}

func GetObject(w http.ResponseWriter, r *http.Request) {
	bucketName, objectKey := ValidatePath(r.URL.Path)
	if strings.TrimSpace(bucketName) == "" || strings.TrimSpace(objectKey) == "" {
		log.Println("Invalid  bucket or object key")
		ErrXMLResponse(w, http.StatusBadRequest, "Invalid bucket or object key")
		return
	}

	bucketPath := GetBucketPath(bucketName)
	if _, err := os.Stat(bucketPath); os.IsNotExist(err) {
		log.Printf("Bucket not found: %s\n", bucketName)
		ErrXMLResponse(w, http.StatusNotFound, "Bucket not found")
		return
	}

	objectPath := filepath.Join(bucketPath, objectKey)
	file, err := os.Open(objectPath)
	if err != nil {
		log.Printf("Failed to open object %s in bucket %s: %v\n", objectKey, bucketName, err)
		ErrXMLResponse(w, http.StatusInternalServerError, "Failed to open object")
		return
	}
	defer file.Close()

	metaDataType, err := GetMetaDataType(bucketPath, objectKey)
	if err != nil {
		log.Printf("Failed to get metadata type for object %s in bucket %s: %v\n", objectKey, bucketName, err)
		ErrXMLResponse(w, http.StatusInternalServerError, "Failed to get metadata type")
		return
	}

	w.Header().Set("Content-Type", metaDataType)
	_, err = io.Copy(w, file)
	if err != nil {
		log.Printf("Failed to read object %s in bucket %s: %v\n", objectKey, bucketName, err)
		ErrXMLResponse(w, http.StatusInternalServerError, "Failed to read object")
		return
	}

	log.Printf("Object %s read successfully in bucket %s", objectKey, bucketName)
}

func DeleteObject(w http.ResponseWriter, r *http.Request) {
	bucketName, objectKey := ValidatePath(r.URL.Path)
	if strings.TrimSpace(bucketName) == "" || strings.TrimSpace(objectKey) == "" {
		log.Println("Invalid  bucket or object key")
		ErrXMLResponse(w, http.StatusBadRequest, "Invalid bucket or object key")
		return
	}

	bucketPath := GetBucketPath(bucketName)
	if _, err := os.Stat(bucketPath); os.IsNotExist(err) {
		log.Printf("Bucket not found: %s\n", bucketName)
		ErrXMLResponse(w, http.StatusNotFound, "Bucket not found")
		return
	}

	objectPath := filepath.Join(bucketPath, objectKey)
	if _, err := os.Stat(objectPath); err != nil {
		if os.IsNotExist(err) {
			log.Printf("Object not found: %s\n", objectKey)
			ErrXMLResponse(w, http.StatusNotFound, "Object not found")
			return
		} else {
			log.Printf("Error with getting stats of object")
			ErrXMLResponse(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
	}

	objects, err := utils.ReadObjects(bucketPath)
	if err != nil {
		log.Printf("Failed to read objects file for bucket %s: %v\n", bucketName, err)
		ErrXMLResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	objectIDX := SearchObjectIDX(objects.Objects, objectKey)
	if objectIDX == -1 {
		log.Printf("Object %s not found in bucket %s\n", objectKey, bucketName)
		ErrXMLResponse(w, http.StatusNotFound, "Object not found")
		return
	}

	objects.Objects = append(objects.Objects[:objectIDX], objects.Objects[objectIDX+1:]...)
	err = utils.WriteObject(bucketName, objects)
	if err != nil {
		log.Printf("Failed to update objects file for bucket %s: %v\n", bucketName, err)
		ErrXMLResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	err = os.Remove(objectPath)
	if err != nil {
		log.Printf("Failed to delete object %s in bucket %s: %v\n", objectKey, bucketName, err)
		ErrXMLResponse(w, http.StatusInternalServerError, "Failed to delete object")
		return
	}

	log.Printf("Object %s deleted successfully from bucket %s\n", objectKey, bucketName)
	w.WriteHeader(http.StatusNoContent)
}
