package utils

import "triple-s/info"

func RecordsToBuckets(records [][]string) info.Buckets {
	var buckets info.Buckets
	for _, record := range records {
		bucket := info.Bucket{
			Name:             record[0],
			CreationTime:     record[1],
			LastModifiedTime: record[2],
			Status:           record[3],
		}
		buckets.Buckets = append(buckets.Buckets, bucket)
	}
	return buckets
}

func BucketsToRecords(buckets info.Buckets) [][]string {
	var records [][]string
	for _, bucket := range buckets.Buckets {
		record := []string{
			bucket.Name,
			bucket.CreationTime,
			bucket.LastModifiedTime,
			bucket.Status,
		}
		records = append(records, record)
	}
	return records
}

func RecordsToObjects(records [][]string) info.Objects {
	var objects info.Objects

	for _, record := range records {
		object := info.Object{
			ObjectKey:    record[0],
			ContentType:  record[1],
			Size:         record[2],
			LastModified: record[3],
		}
		objects.Objects = append(objects.Objects, object)
	}
	return objects
}

func ObjectsToRecords(objects info.Objects) [][]string {
	var records [][]string

	for _, object := range objects.Objects {
		record := []string{
			object.ObjectKey,
			object.ContentType,
			object.Size,
			object.LastModified,
		}
		records = append(records, record)
	}
	return records
}
