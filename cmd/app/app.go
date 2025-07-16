package app

import (
	"1337b04rd/internal/adapters/handlers"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

func RunServer() {
	port := flag.String("port", "8081", "Port to run the web server on")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of triple-s\n")
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr, "\nExample:")
		fmt.Fprintln(os.Stderr, "  go run main.go -port=8081")
	}

	flag.Parse()

	addr := fmt.Sprintf(":%s", *port)
	log.Printf("Server is running on http://localhost%s\n", addr)

	if err := http.ListenAndServe(addr, handlers.RooterWays()); err != nil {
		log.Fatalf("Failed to start server: %v\n", err)
	}
}

// TODO:
