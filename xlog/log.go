package xlog

import (
	"time"
)

// Global logger instance
var (
	_logger ILogger = newZapLogger(&LogConfig{})
)

func Init(logConfig *LogConfig, sinks ...LogSink) {
	_logger = newZapLogger(logConfig, sinks...)
}

// ILogger defines the interface for logging operations
type ILogger interface {
	// Debug logs a message at debug level
	Debug(v ...interface{})
	// Debugf logs a formatted message at debug level
	Debugf(format string, args ...interface{})
	// Info logs a message at info level
	Info(v ...interface{})
	// Infof logs a formatted message at info level
	Infof(format string, args ...interface{})
	// Warn logs a message at warning level
	Warn(v ...interface{})
	// Warnf logs a formatted message at warning level
	Warnf(format string, args ...interface{})
	// Error logs a message at error level
	Error(v ...interface{})
	// Errorf logs a formatted message at error level
	Errorf(format string, args ...interface{})
	// Fatal logs a message at fatal level and then calls os.Exit(1)
	Fatal(v ...interface{})
	// Fatalf logs a formatted message at fatal level and then calls os.Exit(1)
	Fatalf(format string, args ...interface{})

	// Finalize performs any necessary cleanup operations for the logger
	Finalize()
}

// LogConfig holds the configuration for the logger
type LogConfig struct {
	// Level sets the minimum logging level
	Level string
	// TraceLevel sets the minimum level for trace logging
	TraceLevel string
	// File contains configuration for file-based logging
	File *FileLogConfig
}

// FileLogConfig holds the configuration for file-based logging
type FileLogConfig struct {
	// Filename specifies the location of the log file
	Filename string
	// MaxSize is the maximum size in megabytes of the log file before rotation
	MaxSize int
	// MaxBackups is the maximum number of old log files to retain
	MaxBackups int
	// MaxAge is the maximum number of days to retain old log files
	MaxAge int
	// Compress determines if rotated log files should be compressed
	Compress bool
}

type LogSink interface {
	// WriteLog writes a log entry with all fields and optional format arguments
	WriteLog(entry *LogEntry)
}

// LogEntry represents a complete log entry with all fields
type LogEntry struct {
	Level      int
	Time       time.Time
	LoggerName string
	Message    string
	Caller     string
	Stack      string
	// Fields     map[string]interface{}
}

// WriteLog writes a debug level log message
func WriteLog(f func(v ...interface{}), v ...interface{}) {
	if _logger == nil {
		panic("logger is not initialized")
	}

	f(v...)
}

// WriteLogf writes a formatted debug level log message
func WriteLogf(f func(format string, args ...interface{}), format string, args ...interface{}) {
	if _logger == nil {
		panic("logger is not initialized")
	}

	f(format, args...)
}

// Debug logs a message at debug level
func Debug(v ...interface{}) {
	WriteLog(_logger.Debug, v...)
}

// Debugf logs a formatted message at debug level
func Debugf(format string, args ...interface{}) {
	WriteLogf(_logger.Debugf, format, args...)
}

// Info logs a message at info level
func Info(v ...interface{}) {
	WriteLog(_logger.Info, v...)
}

// Infof logs a formatted message at info level
func Infof(format string, args ...interface{}) {
	WriteLogf(_logger.Infof, format, args...)
}

// Warn logs a message at warning level
func Warn(v ...interface{}) {
	WriteLog(_logger.Warn, v...)
}

// Warnf logs a formatted message at warning level
func Warnf(format string, args ...interface{}) {
	WriteLogf(_logger.Warnf, format, args...)
}

// Error logs a message at error level
func Error(v ...interface{}) {
	WriteLog(_logger.Error, v...)
}

// Errorf logs a formatted message at error level
func Errorf(format string, args ...interface{}) {
	WriteLogf(_logger.Errorf, format, args...)
}

// Fatal logs a message at fatal level and then calls os.Exit(1)
func Fatal(v ...interface{}) {
	WriteLog(_logger.Fatal, v...)
}

// Fatalf logs a formatted message at fatal level and then calls os.Exit(1)
func Fatalf(format string, args ...interface{}) {
	WriteLogf(_logger.Fatalf, format, args...)
}

// Finalize performs any necessary cleanup operations for the logger
func Finalize() {
	if _logger != nil {
		_logger.Finalize()
	}
}
