package utils

import (
	"log"
	"path/filepath"
	"triple-s/flags"
	"triple-s/info"
)

func ReadBucket() (info.Buckets, error) {
	bucketPath := filepath.Join(flags.Dir, "buckets.csv")
	log.Printf("Reading buckets metadata: %s\n", bucketPath)

	records, err := ReadCSV(bucketPath)
	if err != nil {
		return info.Buckets{}, err
	}
	return RecordsToBuckets(records), nil
}

func WriteBucket(bucketData info.Buckets) error {
	bucketPath := filepath.Join(flags.Dir, "buckets.csv")
	records := BucketsToRecords(bucketData)
	return WriteCSV(bucketPath, info.BucketsHeader, records)
}
