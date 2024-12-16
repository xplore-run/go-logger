package pkg

import (
	"context"
	"fmt"
	"strings"

	"github.com/rs/zerolog"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

type ContextFieldConfig struct {
	Keys []string
}

type FileSink struct {
	Config        LoggerConfig
	ContextFields ContextFieldConfig
}

var logr *zerolog.Logger

func (fs *FileSink) Init(lc LoggerConfig) error {
	config := lc.FileSinkConfig
	lumberjackLogger := &lumberjack.Logger{
		Filename:   config.FilePath,
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
		Compress:   config.Compress,
	}
	logrr := zerolog.New(lumberjackLogger).With().Timestamp().Logger()
	logr = &logrr
	return nil
}

func (fs *FileSink) Close() error {
	return nil
}

func (fs *FileSink) Info(msg string) {
	logr.Info().Msg(msg)
}

func (fs *FileSink) Warn(msg string) {
	logr.Warn().Msg(msg)
}

func (fs *FileSink) Debug(msg string) {
	logr.Debug().Msg(msg)
}

func (fs *FileSink) Error(msg string) {
	logr.Error().Msg(msg)
}

func (fs *FileSink) FormatMessage(ctx context.Context, message string) string {
	return fs.formatMessageWithDynamicFields(ctx, message)
}

func (fs *FileSink) formatMessageWithDynamicFields(ctx context.Context, msg string) string {
	if ctx == nil {
		return msg
	}

	// Extract dynamic context fields
	var extractedFields []string
	for _, key := range fs.ContextFields.Keys {
		value, ok := ctx.Value(contextKey(key)).(string)
		if ok && value != "" {
			extractedFields = append(extractedFields, fmt.Sprintf("%s:%s", key, value))
		}
	}

	// Construct the formatted message
	if len(extractedFields) > 0 {
		prefix := fmt.Sprintf("[%s]", strings.Join(extractedFields, "]["))
		return fmt.Sprintf("%s - %s", prefix, msg)
	}

	return msg
}

// Utility type for context keys
type contextKey string

// Example usage
func ExampleUsage() {
	// Create a FileSink with specific context field keys
	fileSink := &FileSink{
		ContextFields: ContextFieldConfig{
			Keys: []string{"requestID", "userID", "userType"},
		},
	}

	// Create a context with some values
	ctx := context.Background()
	ctx = context.WithValue(ctx, contextKey("requestID"), "req-123")
	ctx = context.WithValue(ctx, contextKey("userID"), "user-789")
	ctx = context.WithValue(ctx, contextKey("userType"), "admin")

	// Format a message
	formattedMsg := fileSink.FormatMessage(ctx, "Test log message")
	fmt.Println(formattedMsg)
}
