package utils

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Logger is a structured logger for the application
var Logger *logrus.Logger

// InitLogger initializes the application logger
func InitLogger() {
	Logger = logrus.New()

	// Set log level based on environment
	Logger.SetLevel(logrus.DebugLevel)

	// Set output to stdout
	Logger.SetOutput(os.Stdout)

	// Set formatter to JSON for production or text for development
	Logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}

// LogInfo logs an informational message
func LogInfo(message string, fields map[string]any) {
	if Logger == nil {
		InitLogger()
	}

	entry := Logger.WithFields(logrus.Fields(fields))
	entry.Info(message)
}

// LogError logs an error message
func LogError(message string, fields map[string]any) {
	if Logger == nil {
		InitLogger()
	}

	entry := Logger.WithFields(logrus.Fields(fields))
	entry.Error(message)
}

// LogWarn logs a warning message
func LogWarn(message string, fields map[string]any) {
	if Logger == nil {
		InitLogger()
	}

	entry := Logger.WithFields(logrus.Fields(fields))
	entry.Warn(message)
}

// LogDebug logs a debug message
func LogDebug(message string, fields map[string]any) {
	if Logger == nil {
		InitLogger()
	}

	entry := Logger.WithFields(logrus.Fields(fields))
	entry.Debug(message)
}

// LogWithFields adds structured fields to a log entry
func LogWithFields(fields map[string]any) *logrus.Entry {
	if Logger == nil {
		InitLogger()
	}

	return Logger.WithFields(logrus.Fields(fields))
}
