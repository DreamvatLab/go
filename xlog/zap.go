package xlog

import (
	"io"
	"os"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// newZapLogger creates a new ZapLogger instance with the provided configuration
func newZapLogger(config *LogConfig, additionalWriters ...io.Writer) ILogger {
	if config.TraceLevel == "" {
		config.TraceLevel = "error"
	}

	// Configure JSON encoder for structured logging
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)
	encoder := zapcore.NewJSONEncoder(encoderCfg)

	// Create console output for standard output
	consoleWriter := zapcore.AddSync(os.Stdout)

	writers := []zapcore.WriteSyncer{consoleWriter}

	// File writer
	if config.File != nil && config.File.Filename != "" {
		// Set default file logging configuration
		if config.File.MaxSize == 0 {
			config.File.MaxSize = 10 // Max log file size: 10MB
		}
		if config.File.MaxBackups == 0 {
			config.File.MaxBackups = 5 // Max log file backups: 5
		}
		if config.File.MaxAge == 0 {
			config.File.MaxAge = 7 // Max log file age: 7 days
		}

		// Create file logging output with rotation support
		fileWriter := zapcore.AddSync(&lumberjack.Logger{
			Filename:   config.File.Filename,
			MaxSize:    config.File.MaxSize,
			MaxBackups: config.File.MaxBackups,
			MaxAge:     config.File.MaxAge,
			Compress:   config.File.Compress,
		})

		writers = append(writers, fileWriter)
	}

	// Add additional io.Writer to the writers
	for _, additionalWriter := range additionalWriters {
		writers = append(writers, zapcore.AddSync(additionalWriter))
	}

	// Combine multiple writers if needed
	var writer zapcore.WriteSyncer
	if len(writers) > 1 {
		writer = zapcore.NewMultiWriteSyncer(writers...)
	} else {
		writer = writers[0]
	}

	// Parse and set the main logging level
	logLevel, err := zapcore.ParseLevel(config.Level)
	if err != nil {
		panic(err)
	}

	// Create the core logger with the configured encoder, writer, and level
	core := zapcore.NewCore(encoder, writer, logLevel)

	// Parse and set the trace level for stack traces
	traceLevel, err := zapcore.ParseLevel(config.TraceLevel)
	if err != nil {
		traceLevel = zapcore.ErrorLevel
	}

	// Create the final logger with stack trace support
	zapLogger := zap.New(core, zap.AddStacktrace(traceLevel))

	return &ZapLogger{
		config:      config,
		innerLogger: zapLogger,
	}
}

// ZapLogger implements the ILogger interface using the zap logging library
type ZapLogger struct {
	config      *LogConfig
	innerLogger *zap.Logger
}

// SetConfig sets the configuration for the logger
func (o *ZapLogger) SetConfig(config *LogConfig) {
	o.config = config
}

// GetConfig returns the current configuration for the logger
func (o *ZapLogger) GetConfig() *LogConfig {
	return o.config
}

// Debug logs a message at debug level
func (o *ZapLogger) Debug(v ...interface{}) {
	o.innerLogger.Sugar().Debug(v...)
}

// Debugf logs a formatted message at debug level
func (o *ZapLogger) Debugf(format string, args ...interface{}) {
	o.innerLogger.Sugar().Debugf(format, args...)
}

// Info logs a message at info level
func (o *ZapLogger) Info(v ...interface{}) {
	o.innerLogger.Sugar().Info(v...)
}

// Infof logs a formatted message at info level
func (o *ZapLogger) Infof(format string, args ...interface{}) {
	o.innerLogger.Sugar().Infof(format, args...)
}

// Warn logs a message at warning level
func (o *ZapLogger) Warn(v ...interface{}) {
	o.innerLogger.Sugar().Warn(v...)
}

// Warnf logs a formatted message at warning level
func (o *ZapLogger) Warnf(format string, args ...interface{}) {
	o.innerLogger.Sugar().Warnf(format, args...)
}

// Error logs a message at error level
func (o *ZapLogger) Error(v ...interface{}) {
	o.innerLogger.Sugar().Error(v...)
}

// Errorf logs a formatted message at error level
func (o *ZapLogger) Errorf(format string, args ...interface{}) {
	o.innerLogger.Sugar().Errorf(format, args...)
}

// Fatal logs a message at fatal level and then calls os.Exit(1)
func (o *ZapLogger) Fatal(v ...interface{}) {
	o.innerLogger.Sugar().Fatal(v...)
}

// Fatalf logs a formatted message at fatal level and then calls os.Exit(1)
func (o *ZapLogger) Fatalf(format string, args ...interface{}) {
	o.innerLogger.Sugar().Fatalf(format, args...)
}

// Finalize performs cleanup operations for the logger
func (o *ZapLogger) Finalize() {
	if o.innerLogger != nil {
		o.innerLogger.Sync()
	}
}
