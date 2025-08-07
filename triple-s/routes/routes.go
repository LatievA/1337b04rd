package routes

import (
	"net/http"
	"triple-s/storage"
)

func Routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("PUT /{BucketName}", storage.CreateBucket)
	mux.HandleFunc("GET /", storage.ListBuckets)
	mux.HandleFunc("DELETE /{BucketName}", storage.DeleteBucket)

	mux.HandleFunc("PUT /{BucketName}/{ObjectKey}", storage.CreateObject)
	mux.HandleFunc("GET /{BucketName}", storage.ListObjects)
	mux.HandleFunc("GET /{BucketName}/{ObjectKey}", storage.GetObject)
	mux.HandleFunc("DELETE /{BucketName}/{ObjectKey}", storage.DeleteObject)
	mux.HandleFunc("GET /health", storage.HealthCheckHandler)

	return mux
}
