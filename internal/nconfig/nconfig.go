package nconfig

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	ServerAddr     string
	ReportInterval int
	PollInterval   int
	StoreInterval  int
	FilePath       string
	Restore        bool
}

func getEnvOrDefault(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvIntOrDefault(key string, defaultValue int) int {
	if envValue := getEnvOrDefault(key, ""); envValue != "" {
		if value, err := strconv.Atoi(envValue); err == nil {
			return value
		}
	}
	return defaultValue
}

func getEnvBoolOrDefault(key string, defaultValue bool) bool {
	if envValue := getEnvOrDefault(key, ""); envValue != "" {
		lower := strings.ToLower(envValue)
		return lower == "true" || lower == "1" || lower == "yes" || lower == "y"
	}
	return defaultValue
}

func ParseConfig() Config {
	serverAddr := flag.String("a", "localhost:8080", "HTTP server address")
	reportInterval := flag.Int("r", 10, "Report interval in seconds")
	pollInterval := flag.Int("p", 2, "Poll interval in seconds")
	storeInterval := flag.Int("i", 300, "Store interval in seconds (0 = synchronous write)")
	filePath := flag.String("f", "metrics_storage.json", "Metrics storage file path")
	restore := flag.Bool("R", false, "Restore metrics from file on startup")

	flag.Parse()

	if flag.NArg() > 0 {
		fmt.Printf("Unknown arguments: %v\n", flag.Args())
		flag.Usage()
		os.Exit(1)
	}

	addr := getEnvOrDefault("ADDRESS", *serverAddr)

	reportInt := getEnvIntOrDefault("REPORT_INTERVAL", *reportInterval)
	pollInt := getEnvIntOrDefault("POLL_INTERVAL", *pollInterval)
	storeInt := getEnvIntOrDefault("STORE_INTERVAL", *storeInterval)
	fileStoragePath := getEnvOrDefault("FILE_STORAGE_PATH", *filePath)
	restoreMetrics := getEnvBoolOrDefault("RESTORE", *restore)

	return Config{
		ServerAddr:     addr,
		ReportInterval: reportInt,
		PollInterval:   pollInt,
		StoreInterval:  storeInt,
		FilePath:       fileStoragePath,
		Restore:        restoreMetrics,
	}
}
