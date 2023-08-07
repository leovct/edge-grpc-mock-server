// logger/logger.go

package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

// LoggerConfig contains configurations for the logger.
type LoggerConfig struct {
	Level       zerolog.Level
	CallerField string
}

// NewLogger creates a new Zerolog logger with the given configuration.
func NewLogger(config LoggerConfig) zerolog.Logger {
	zerolog.TimeFieldFormat = time.RFC822
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.UnixDate,
	}
	return zerolog.New(output).
		Level(zerolog.InfoLevel).
		With().
		Caller().
		Timestamp().
		Logger()
}
