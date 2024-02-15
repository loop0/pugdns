package utils

import (
	"log/slog"
	"os"
)

func GetEnvOrExit(name string) string {
	value := os.Getenv(name)
	if value == "" {
		slog.Error("Missing environment variable", name, "")
		os.Exit(1)
	}
	return value
}
