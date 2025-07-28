package config

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strconv"
)

type Config struct {
	ServerConfig *ServerConfig
	DBConfig     *DBConfig
	S3Config	 *S3Config
}

type ServerConfig struct {
	Port string
}

type DBConfig struct {
	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string
}

type S3Config struct {
	BaseURL string
}

func NewConfig() (*Config, error) {
	dbConfig := &DBConfig{
		DBHost:     getEnv("DB_HOST", "db"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBName:     getEnv("DB_NAME", "1337b04rd"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
	}

	s3Config := &S3Config{
		BaseURL: getEnv("S3_BASE_URL", "http://localhost:8080"),
	}

	serverConfig := &ServerConfig{
		Port: getEnv("SERVER_PORT", "8081"),
	}
	err := parseFlags(serverConfig)
	if err != nil {
		return nil, err
	}

	return &Config{
		DBConfig:     dbConfig,
		ServerConfig: serverConfig,
		S3Config:     s3Config,
	}, nil
}

func getEnv(key, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		slog.Warn("Missing environment variable, using default values!", "keys", key, "default value", defaultVal)
		return defaultVal
	}

	return val
}

func parseFlags(serverConfig *ServerConfig) error {
	port := flag.Int("port", 0, "Port to serve on")
	flag.Usage = func() {
		printHelp()
	}

	flag.Parse()

	if *port == 0 {
		return nil
	}

	if *port < 1024 || *port > 65565 {
		return fmt.Errorf("port number should be between 1024 and 65565")
	}
	portConv := strconv.Itoa(*port)
	serverConfig.Port = portConv
	return nil
}

func printHelp() {
	fmt.Println(`hacker board

Usage:
  1337b04rd [--port <N>]  
  1337b04rd --help

Options:
  --help       Show this screen.
  --port N     Port number.`)
}
