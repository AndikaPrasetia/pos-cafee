package utils

import (
	"crypto/rand"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// Logger is a structured logger for the application
var Logger *logrus.Logger

// InitLogger initializes the application logger with environment-aware configuration
func InitLogger(environment string, logLevel string) {
	if Logger != nil {
		return // Already initialized
	}

	Logger = logrus.New()

	// Determine log level
	level, err := logrus.ParseLevel(strings.ToLower(logLevel))
	if err != nil {
		level = logrus.InfoLevel // Default fallback
	}
	Logger.SetLevel(level)

	// Set output to stdout
	Logger.SetOutput(os.Stdout)

	// Set formatter based on environment
	if strings.ToLower(environment) == "production" {
		// JSON formatter for production (better for log aggregation systems)
		Logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
		})
	} else {
		// Text formatter for development (more human-readable)
		Logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
			ForceColors:     true,
		})
	}
}

// LogInfo logs an informational message
func LogInfo(message string, fields map[string]any) {
	if Logger == nil {
		// Initialize with defaults if not already done
		InitLogger("development", "info")
	}

	entry := Logger.WithFields(logrus.Fields(fields))
	entry.Info(message)
}

// LogError logs an error message
func LogError(message string, fields map[string]any) {
	if Logger == nil {
		// Initialize with defaults if not already done
		InitLogger("development", "error")
	}

	entry := Logger.WithFields(logrus.Fields(fields))
	entry.Error(message)
}

// LogWarn logs a warning message
func LogWarn(message string, fields map[string]any) {
	if Logger == nil {
		// Initialize with defaults if not already done
		InitLogger("development", "warn")
	}

	entry := Logger.WithFields(logrus.Fields(fields))
	entry.Warn(message)
}

// LogDebug logs a debug message
func LogDebug(message string, fields map[string]any) {
	if Logger == nil {
		// Initialize with defaults if not already done
		InitLogger("development", "debug")
	}

	entry := Logger.WithFields(logrus.Fields(fields))
	entry.Debug(message)
}

// LogWithFields adds structured fields to a log entry
func LogWithFields(fields map[string]any) *logrus.Entry {
	if Logger == nil {
		// Initialize with defaults if not already done
		InitLogger("development", "info")
	}

	return Logger.WithFields(logrus.Fields(fields))
}

// LogFatal logs a fatal message and exits the application
func LogFatal(message string, fields map[string]any) {
	if Logger == nil {
		// Initialize with defaults if not already done
		InitLogger("development", "fatal")
	}

	entry := Logger.WithFields(logrus.Fields(fields))
	entry.Fatal(message)
}

// LogPanic logs a panic message and panics
func LogPanic(message string, fields map[string]any) {
	if Logger == nil {
		// Initialize with defaults if not already done
		InitLogger("development", "panic")
	}

	entry := Logger.WithFields(logrus.Fields(fields))
	entry.Panic(message)
}

// GenerateUUID generates a random UUID string
func GenerateUUID() string {
	uuid := make([]byte, 16)
	_, err := rand.Read(uuid)
	if err != nil {
		// Fallback to timestamp if random generation fails
		return fmt.Sprintf("uuid-%d", time.Now().UnixNano())
	}

	// Convert to standard UUID format
	uuid[6] = (uuid[6] & 0x0f) | 0x40 // Version 4
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // Variant 10

	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
}
