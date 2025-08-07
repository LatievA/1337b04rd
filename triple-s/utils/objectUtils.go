package utils

import (
	"fmt"
	"log"
	"path/filepath"
	"triple-s/flags"
	"triple-s/info"
)

func ReadObjects(bucketName string) (info.Objects, error) {
	objectPath := filepath.Join(bucketName, "objects.csv")
	log.Printf("Reading objects metadata from %s\n", objectPath)

	records, err := ReadCSV(objectPath)
	fmt.Println(records)
	if err != nil {
		return info.Objects{}, err
	}

	return RecordsToObjects(records), nil
}

func WriteObject(bucketName string, objects info.Objects) error {
	objectPath := filepath.Join(flags.Dir, bucketName, "objects.csv")
	records := ObjectsToRecords(objects)
	return WriteCSV(objectPath, info.ObjectsHeader, records)
}
