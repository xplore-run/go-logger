// Package pkg provides a flexible logging framework with support for multiple output sinks
// and configurable log levels. It allows structured logging with context-based fields
// and rotating file output.
package pkg

import (
	"context"
	"fmt"
	"log"
	"time"
)

// LogLevel represents the severity of a log message
type LogLevel int

const (
	INFO  LogLevel = iota // Default level for general operational information
	DEBUG                 // Detailed information for debugging purposes
	WARN                  // Warning messages for potentially harmful situations
	ERROR                 // Error messages for serious problems
	FATAL                 // Fatal messages that will terminate the application
)

// SinkType defines the type of output destination for logs
type SinkType int

const (
	FILE SinkType = iota // File-based logging output
)

// LoggerConfig holds the configuration for initializing a logger
type LoggerConfig struct {
	SinkType       SinkType              // Type of sink to use (e.g., FILE)
	Level          LogLevel              // Minimum log level to output
	FileSinkConfig *LoggerFileSinkConfig // Configuration for file-based logging
	TimeFormat     string                // Format for timestamp in logs
	BatchSize      int                   // Number of messages to batch before writing
	FlushTimeout   time.Duration         // Maximum time to wait before flushing logs
}

// LoggerFileSinkConfig configures the behavior of file-based logging
type LoggerFileSinkConfig struct {
	FilePath   string // Path to the log file
	MaxSize    int    // Maximum size in megabytes before rotating
	MaxBackups int    // Maximum number of old log files to retain
	MaxAge     int    // Maximum number of days to retain old log files
	Compress   bool   // Whether to compress rotated log files
}

// LogMessage represents a single log entry
type LogMessage struct {
	Level     LogLevel // Severity level of the message
	Namespace string   // Namespace/category for the message
	Content   string   // Actual log message content
	Timestamp string   // Time when the message was created
}

// CustomLogger implements the main logging functionality
type CustomLogger struct {
	sink   Sink
	config LoggerConfig
}

// Sink defines the interface for log output destinations
type Sink interface {
	Init(LoggerConfig) error                              // Initialize the sink with configuration
	Close() error                                         // Clean up resources
	Info(msg string)                                      // Log an info message
	Warn(msg string)                                      // Log a warning message
	Debug(msg string)                                     // Log a debug message
	Error(msg string)                                     // Log an error message
	FormatMessage(ctx context.Context, msg string) string // Format a message with context
}

// NewCustomLogger creates and initializes a new logger with the provided configuration
func NewCustomLogger(config LoggerConfig) (*CustomLogger, error) {
	customLogger := &CustomLogger{
		config: config,
	}
	var err error

	if config.SinkType == FILE {
		var sink FileSink

		err = sink.Init(config)

		if err != nil {
			return nil, err
		}

		customLogger.sink = &sink

	} else {
		return nil, fmt.Errorf("invalid sink type")
	}

	return customLogger, nil
}

// Close cleanly shuts down the logger
func (c *CustomLogger) Close() error {
	return c.sink.Close()
}

// Error logs a message at ERROR level if the logger's level permits
func (c *CustomLogger) Error(ctx context.Context, msg string) {
	if c.config.Level == ERROR || c.config.Level == WARN || c.config.Level == INFO || c.config.Level == DEBUG {
		c.formatAndLog(ctx, ERROR, msg)
	}
}

// Warn logs a message at WARN level if the logger's level permits
func (c *CustomLogger) Warn(ctx context.Context, msg string) {
	if c.config.Level == WARN || c.config.Level == INFO || c.config.Level == DEBUG {
		c.formatAndLog(ctx, WARN, msg)
	}
}

// Info logs a message at INFO level if the logger's level permits
func (c *CustomLogger) Info(ctx context.Context, msg string) {
	if c.config.Level == INFO || c.config.Level == DEBUG {
		c.formatAndLog(ctx, INFO, msg)
	}
}

// Debug logs a message at DEBUG level if the logger's level permits
func (c *CustomLogger) Debug(ctx context.Context, msg string) {
	if c.config.Level == DEBUG {
		c.formatAndLog(ctx, DEBUG, msg)
	}
}

// formatAndLog handles message formatting and panic recovery
func (c *CustomLogger) formatAndLog(ctx context.Context, level LogLevel, msg string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("[go-logger] Panic: %v", r)
		}
	}()

	c.doLog(level, c.sink.FormatMessage(ctx, msg))
}

// doLog routes the message to the appropriate sink method based on level
func (c *CustomLogger) doLog(level LogLevel, msg string) {
	switch level {
	case DEBUG:
		c.sink.Debug(msg)
	case INFO:
		c.sink.Info(msg)
	case WARN:
		c.sink.Warn(msg)
	default:
		c.sink.Error(msg)
	}
}

// InitLoggerExample provides an example configuration for quick setup
// It creates a logger that writes to a file in the specified directory
func InitLoggerExample(logFolder, logFileName string) *CustomLogger {
	loggerConfig := LoggerConfig{
		TimeFormat: "02-01-2006 15:04:05",
		SinkType:   FILE,
		FileSinkConfig: &LoggerFileSinkConfig{
			FilePath:   fmt.Sprintf("%s/logs/%s", logFolder, logFileName),
			MaxSize:    1,
			MaxBackups: 2,
			MaxAge:     1,
			Compress:   true,
		},
		BatchSize:    100,
		FlushTimeout: 5 * time.Second,
	}

	customLogger, err := NewCustomLogger(loggerConfig)
	if err != nil {
		log.Fatalf("Failed to create custom logger: %s", err.Error())
	}
	return customLogger
}
