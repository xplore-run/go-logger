package pkg

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLogger(t *testing.T) {
	tempDir := os.TempDir()

	t.Run("test logger initialization", func(t *testing.T) {
		config := LoggerConfig{
			TimeFormat: "02-01-2006 15:04:05",
			SinkType:   FILE,
			FileSinkConfig: &LoggerFileSinkConfig{
				FilePath:   tempDir + "/test.log",
				MaxSize:    1,
				MaxBackups: 2,
				MaxAge:     1,
				Compress:   true,
			},
			BatchSize:    10,
			FlushTimeout: 5 * time.Second,
		}

		logger, err := NewCustomLogger(config)
		assert.NoError(t, err)
		assert.NotNil(t, logger)
		defer logger.Close()
	})

	t.Run("test invalid sink type", func(t *testing.T) {
		config := LoggerConfig{
			SinkType: SinkType(999), // Invalid sink type
		}
		logger, err := NewCustomLogger(config)
		assert.Error(t, err)
		assert.Nil(t, logger)
	})

	t.Run("test log levels", func(t *testing.T) {
		config := LoggerConfig{
			TimeFormat: "02-01-2006 15:04:05",
			SinkType:   FILE,
			Level:      DEBUG,
			FileSinkConfig: &LoggerFileSinkConfig{
				FilePath:   tempDir + "/test_levels.log",
				MaxSize:    1,
				MaxBackups: 2,
				MaxAge:     1,
				Compress:   true,
			},
		}

		logger, err := NewCustomLogger(config)
		assert.NoError(t, err)
		defer logger.Close()

		ctx := context.Background()
		logger.Debug(ctx, "debug message")
		logger.Info(ctx, "info message")
		logger.Warn(ctx, "warn message")
		logger.Error(ctx, "error message")
	})

	t.Run("test log level filtering", func(t *testing.T) {
		config := LoggerConfig{
			TimeFormat: "02-01-2006 15:04:05",
			SinkType:   FILE,
			Level:      INFO, // Only INFO and above should be logged
			FileSinkConfig: &LoggerFileSinkConfig{
				FilePath:   tempDir + "/test_filtering.log",
				MaxSize:    1,
				MaxBackups: 2,
				MaxAge:     1,
				Compress:   true,
			},
		}

		logger, err := NewCustomLogger(config)
		assert.NoError(t, err)
		defer logger.Close()

		ctx := context.Background()
		logger.Debug(ctx, "debug message") // Shouldn't be logged
		logger.Info(ctx, "info message")   // Should be logged
	})

	t.Run("test context formatting", func(t *testing.T) {
		config := LoggerConfig{
			TimeFormat: "02-01-2006 15:04:05",
			SinkType:   FILE,
			Level:      INFO,
			FileSinkConfig: &LoggerFileSinkConfig{
				FilePath:   tempDir + "/test_context.log",
				MaxSize:    1,
				MaxBackups: 2,
				MaxAge:     1,
				Compress:   true,
			},
		}

		logger, err := NewCustomLogger(config)
		assert.NoError(t, err)
		defer logger.Close()

		ctx := context.WithValue(context.Background(), contextKey("requestID"), "123")
		logger.Info(ctx, "test message with context")
	})
}
