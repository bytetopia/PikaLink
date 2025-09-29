package utils

import (
	"os"
	"path/filepath"
)

// GetDataPath returns the data path from environment variable or current directory as fallback
func GetDataPath() string {
	dataPath := os.Getenv("DATA_PATH")
	if dataPath == "" {
		dataPath = "."
	}
	return dataPath
}

// GetLogsDir returns the full path to the logs directory
func GetLogsDir() string {
	dataPath := GetDataPath()
	return filepath.Join(dataPath, "logs")
}

// GetMainDbPath returns the full path to the main pikalink.db database file
func GetMainDbPath() string {
	dataPath := GetDataPath()
	return filepath.Join(dataPath, "pikalink.db")
}

// GetLogDbPath returns the full path to a specific month's log database file
func GetLogDbPath(month string) string {
	return filepath.Join(GetLogsDir(), month+".db")
}