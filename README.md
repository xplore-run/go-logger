# go-logger

A flexible and extensible logging framework for Go applications with support for multiple output sinks, log rotation, and context-based structured logging.

It uses github.com/rs/zerolog and gopkg.in/natefinch/lumberjack.v2 libraries.

## Features
- Multiple log levels (DEBUG, INFO, WARN, ERROR, FATAL)
- File-based logging with rotation support
- Context-based structured logging
- Configurable timestamp formats
- Message batching and flush timeouts
- Compressed log archives

## Installation
```bash
go get github.com/xplore-run/go-logger
```

## Usage
```go
package main

import (
    "context"
    logger "github.com/xplore-run/go-logger/pkg"
)

func main() {
    // Initialize logger with file sink
    config := logger.LoggerConfig{
        SinkType: logger.FILE,
        Level:    logger.INFO,
        FileSinkConfig: &logger.LoggerFileSinkConfig{
            FilePath:   "/var/log/app.log",
            MaxSize:    10,    // 10MB
            MaxBackups: 5,     // Keep 5 old files
            MaxAge:     7,     // 7 days
            Compress:   true,
        },
        TimeFormat:   "2006-01-02 15:04:05",
        BatchSize:    100,
        FlushTimeout: 5 * time.Second,
    }

    log, err := logger.NewCustomLogger(config)
    if err != nil {
        panic(err)
    }
    defer log.Close()

    // Basic logging
    log.Info(context.Background(), "Application started")

    // With context
    ctx := context.WithValue(context.Background(), "requestID", "123")
    log.Info(ctx, "Processing request")
}
```

## Configuration
### LoggerConfig
| Field          | Type                | Description                           |
|----------------|---------------------|---------------------------------------|
| SinkType       | SinkType           | Output destination (FILE)             |
| Level          | LogLevel           | Minimum log level to output           |
| FileSinkConfig | LoggerFileSinkConfig| File-based logging configuration     |
| TimeFormat     | string             | Timestamp format                      |
| BatchSize      | int                | Number of messages to batch           |
| FlushTimeout   | time.Duration      | Maximum time before forcing flush     |

### Log Levels
- DEBUG: Detailed debugging information
- INFO: General operational information
- WARN: Warning messages
- ERROR: Error conditions
- FATAL: Critical errors that stop the application

### Running test cases
go test -v ./pkg/...

## View docs
go doc -all pkg

## Contributing
1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License
This project is licensed under the MIT License - see the LICENSE file for details.
```