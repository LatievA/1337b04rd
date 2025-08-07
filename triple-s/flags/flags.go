package flags

import (
	"flag"
	"fmt"
	"log"
	"strings"
)

var (
	Port int
	Dir  string
)

func ParseFlags() error {
	flag.IntVar(&Port, "port", 8080, "Port to serve one")
	flag.StringVar(&Dir, "dir", "./data", "The directory of files to host")
	flag.Usage = PrintHelp

	flag.Parse()

	if Dir == "" {
		log.Fatal("Error: directory is required")
	}

	if Port <= 0 || Port > 65535 {
		log.Fatal("Invalid port number")
	}

	prohibitedDirs := []string{"home", "arch", "go.mod", "main.go", "README.md", "base", "handlers", "utils", "storage", "routes", "info", "flags"}
	for _, dir := range prohibitedDirs {
		if strings.Contains(Dir, dir) {
			log.Fatalf("%s directory is not allowed. Please provide a valid directory.", Dir)
		}
	}

	log.Printf("Serving files from directory: %s", Dir)
	log.Printf("Listening port on: %d", Port)

	return nil
}

func PrintHelp() {
	fmt.Println(`Simple Storage Service.

	**Usage:**
		triple-s [-port <N>] [-dir <S>]  
		triple-s --help
	
	**Options:**
	- --help     Show this screen.
	- --port N   Port number
	- --dir S    Path to the directory`)
}
