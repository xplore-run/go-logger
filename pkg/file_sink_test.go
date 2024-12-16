package pkg

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileSink(t *testing.T) {
	tempDir := os.TempDir()

	t.Run("test initialization", func(t *testing.T) {
		sink := &FileSink{}
		config := LoggerConfig{
			FileSinkConfig: &LoggerFileSinkConfig{
				FilePath:   filepath.Join(tempDir, "test_init.log"),
				MaxSize:    1,
				MaxBackups: 1,
				MaxAge:     1,
				Compress:   true,
			},
		}

		err := sink.Init(config)
		assert.NoError(t, err)
		defer sink.Close()
	})

	t.Run("test message formatting with context", func(t *testing.T) {
		sink := &FileSink{
			ContextFields: ContextFieldConfig{
				Keys: []string{"requestID", "userID"},
			},
		}

		ctx := context.Background()
		ctx = context.WithValue(ctx, contextKey("requestID"), "req123")
		ctx = context.WithValue(ctx, contextKey("userID"), "user456")

		formatted := sink.FormatMessage(ctx, "test message")
		assert.Contains(t, formatted, "requestID:req123")
		assert.Contains(t, formatted, "userID:user456")
		assert.Contains(t, formatted, "test message")
	})

	t.Run("test message formatting with empty context", func(t *testing.T) {
		sink := &FileSink{
			ContextFields: ContextFieldConfig{
				Keys: []string{"requestID"},
			},
		}

		formatted := sink.FormatMessage(nil, "test message")
		assert.Equal(t, "test message", formatted)
	})

	t.Run("test logging functions", func(t *testing.T) {
		sink := &FileSink{}
		config := LoggerConfig{
			FileSinkConfig: &LoggerFileSinkConfig{
				FilePath: filepath.Join(tempDir, "test_logging.log"),
			},
		}

		err := sink.Init(config)
		assert.NoError(t, err)
		defer sink.Close()

		// Test all log levels
		sink.Info("info message")
		sink.Debug("debug message")
		sink.Warn("warn message")
		sink.Error("error message")

		// Verify file exists
		_, err = os.Stat(config.FileSinkConfig.FilePath)
		assert.NoError(t, err)
	})

	t.Run("test context field extraction", func(t *testing.T) {
		sink := &FileSink{
			ContextFields: ContextFieldConfig{
				Keys: []string{"key1", "key2", "key3"},
			},
		}

		ctx := context.Background()
		ctx = context.WithValue(ctx, contextKey("key1"), "value1")
		ctx = context.WithValue(ctx, contextKey("key2"), "value2")
		// key3 intentionally left empty

		formatted := sink.formatMessageWithDynamicFields(ctx, "test")
		assert.Contains(t, formatted, "key1:value1")
		assert.Contains(t, formatted, "key2:value2")
		assert.NotContains(t, formatted, "key3")
	})

	t.Run("test file rotation", func(t *testing.T) {
		sink := &FileSink{}
		config := LoggerConfig{
			FileSinkConfig: &LoggerFileSinkConfig{
				FilePath:   filepath.Join(tempDir, "test_rotation.log"),
				MaxSize:    1,
				MaxBackups: 1,
				MaxAge:     1,
				Compress:   true,
			},
		}

		err := sink.Init(config)
		assert.NoError(t, err)
		defer sink.Close()

		// Write enough logs to trigger rotation
		for i := 0; i < 10000; i++ {
			sink.Info("test message for rotation")
		}

		// Check if backup file exists
		backupPattern := config.FileSinkConfig.FilePath + ".*"
		matches, err := filepath.Glob(backupPattern)
		assert.NoError(t, err)
		assert.Greater(t, len(matches), 0)
	})
}
