package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	// Set up global logger
	zerolog.TimeFieldFormat = time.RFC3339
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})

	// Set log level based on environment
	if os.Getenv("GIN_MODE") == "release" {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
}

// Fields type for structured logging
type Fields map[string]interface{}

// Debug logs a debug message with fields
func Debug(message string, fields Fields) {
	log.Debug().Fields(fields).Msg(message)
}

// Info logs an info message with fields
func Info(message string, fields Fields) {
	log.Info().Fields(fields).Msg(message)
}

// Warn logs a warning message with fields
func Warn(message string, fields Fields) {
	log.Warn().Fields(fields).Msg(message)
}

// Error logs an error message with fields
func Error(message string, err error, fields Fields) {
	if fields == nil {
		fields = Fields{}
	}
	fields["error"] = err
	log.Error().Fields(fields).Msg(message)
}

// Fatal logs a fatal message with fields and exits
func Fatal(message string, err error, fields Fields) {
	if fields == nil {
		fields = Fields{}
	}
	fields["error"] = err
	log.Fatal().Fields(fields).Msg(message)
}
