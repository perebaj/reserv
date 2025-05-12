// Package reserv slog.go gather all important functions to create a slog logger.
package reserv

import (
	"fmt"
	"log/slog"
	"math"
	"os"
	"strings"
)

// ConfigSlog gather all the configuration options for the slog logger.
type ConfigSlog struct {
	// Level represents the available log levels. Available values are: debug, info, warn, error, none
	Level string
	// Format represents the available log formats. Available values are: logfmt, json, gcp
	Format string
}

const (
	// LevelInfo is the log level for info logs.
	LevelInfo = "info"
	// LevelDebug is the log level for debug logs.
	LevelDebug = "debug"
	// LevelWarn is the log level for warn logs.
	LevelWarn = "warn"
	// LevelError is the log level for error logs.
	LevelError = "error"
	// LevelNone is the log level for no logs.
	LevelNone = "none"
)

const (
	// FormatLogFmt is the default log format. It is a human-readable format that is easy to read.
	FormatLogFmt = "logfmt"
	// FormatJSON is a machine-readable format that is easy to parse. It is the default format for the slog logger.
	FormatJSON = "json"
	// FormatGCP is the Google Cloud Platform log format. It is a machine-readable format that is easy to parse.
	FormatGCP = "gcp"
)

// AvailableLogLevels is a list of supported logging levels
var AvailableLogLevels = []string{
	LevelDebug,
	LevelInfo,
	LevelWarn,
	LevelError,
	LevelNone,
}

// AvailableLogFormats is a list of supported log formats
var AvailableLogFormats = []string{
	FormatLogFmt,
	FormatJSON,
	FormatGCP,
}

// NewLoggerSlog returns a *slog.Logger that prints in the provided format at the
// provided level.
func NewLoggerSlog(c ConfigSlog) (*slog.Logger, error) {
	lvlOption, err := parseLevel(c.Level)
	if err != nil {
		return nil, err
	}

	handler, err := getHandlerFromFormat(c.Format, slog.HandlerOptions{
		Level:     lvlOption,
		AddSource: true,
	})
	if err != nil {
		return nil, err
	}

	return slog.New(handler), nil
}

// getHandlerFromFormat returns a slog.Handler based on the provided format and slog options.
func getHandlerFromFormat(format string, opts slog.HandlerOptions) (slog.Handler, error) {
	var handler slog.Handler
	switch strings.ToLower(format) {
	case FormatLogFmt:
		handler = slog.NewTextHandler(os.Stdout, &opts)
		return handler, nil
	case FormatJSON:
		handler = slog.NewJSONHandler(os.Stdout, &opts)
		return handler, nil
	case FormatGCP:
		opts.ReplaceAttr = replaceSlogAttributesGCP
		handler = slog.NewJSONHandler(os.Stdout, &opts)
		return handler, nil
	default:
		return nil, fmt.Errorf("log format %s unknown, %v are possible values", format, AvailableLogFormats)
	}
}

// replaceSlogAttributesGCP returns a slog.Handler that formats logs in a way that is
// compatible with Google Cloud Logging.
func replaceSlogAttributesGCP(_ []string, a slog.Attr) slog.Attr {
	// Customize the name of some fields to match Google Cloud expectations
	// More: https://cloud.google.com/logging/docs/agent/logging/configuration#process-payload
	if a.Key == slog.LevelKey {
		return slog.Attr{
			Key:   "severity",
			Value: a.Value,
		}
	}

	if a.Key == slog.MessageKey {
		return slog.Attr{
			Key:   "message",
			Value: a.Value,
		}
	}

	return a
}

// parseLevel returns the slog.Level based on the provided string.
func parseLevel(lvl string) (slog.Level, error) {
	switch strings.ToLower(lvl) {
	case LevelDebug:
		return slog.LevelDebug, nil
	case LevelInfo:
		return slog.LevelInfo, nil
	case LevelWarn:
		return slog.LevelWarn, nil
	case LevelError:
		return slog.LevelError, nil
	case LevelNone:
		return slog.Level(math.MaxInt), nil
	default:
		return -1, fmt.Errorf("log log_level %s unknown, %v are possible values", lvl, AvailableLogLevels)
	}
}
