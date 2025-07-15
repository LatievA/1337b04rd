package handlers

import "net/http"

func RooterWays() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", GetCatalog)

	return mux
}

// Test method
func GetCatalog(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Welcome to the Catalog!"))
}
