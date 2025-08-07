package main

import (
	"fmt"
	"log"
	"net/http"
	"triple-s/flags"
	"triple-s/routes"
	"triple-s/storage"
)

func main() {
	flags.ParseFlags()
	err := storage.InitDir()
	if err != nil {
		log.Fatal(err)
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", flags.Port),
		Handler: routes.Routes(),
	}

	log.Printf("Starting the server on http://localhost:%d", flags.Port)
	log.Printf("Data directory: %s", flags.Dir)

	err = server.ListenAndServe()
	log.Fatal(err)
}
