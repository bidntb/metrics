package nconfig

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

func GetServerAddress() string {
	defaultAddress := "localhost:8080"

	if envAddress := os.Getenv("ADDRESS"); envAddress != "" {
		defaultAddress = envAddress
	}

	serverAddress := flag.String("a", defaultAddress, "HTTP server endpoint address")
	flag.Parse()

	finalAddress := *serverAddress
	if envAddress := os.Getenv("ADDRESS"); envAddress != "" {
		finalAddress = envAddress
	}

	return finalAddress
}

func getEnvOrDefault(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func ParseFlags() (string, int, int) {
	serverAddr := flag.String("a", "localhost:8080", "HTTP server address")
	reportInterval := flag.Int("r", 10, "Report interval in seconds")
	pollInterval := flag.Int("p", 2, "Poll interval in seconds")

	flag.Parse()

	if flag.NArg() > 0 {
		fmt.Printf("Unknown arguments: %v\n", flag.Args())
		flag.Usage()
		return "", 0, 0
	}

	address := getEnvOrDefault("ADDRESS", *serverAddr)

	reportInt := *reportInterval
	if envReportInterval := getEnvOrDefault("REPORT_INTERVAL", ""); envReportInterval != "" {
		if value, err := strconv.Atoi(envReportInterval); err == nil {
			reportInt = value
		}
	}

	pollInt := *pollInterval
	if envPollInterval := getEnvOrDefault("POLL_INTERVAL", ""); envPollInterval != "" {
		if value, err := strconv.Atoi(envPollInterval); err == nil {
			pollInt = value
		}
	}

	return address, reportInt, pollInt
}
