package xlog

import (
	"os"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// newZapLogger creates a new ZapLogger instance with the provided configuration

func newZapLogger(config *LogConfig, sinks ...LogSink) ILogger {
	if config.TraceLevel == "" {
		config.TraceLevel = "error"
	}

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)
	encoder := zapcore.NewJSONEncoder(encoderCfg)

	consoleWriter := zapcore.AddSync(os.Stdout)
	writers := []zapcore.WriteSyncer{consoleWriter}

	if config.File != nil && config.File.Filename != "" {
		if config.File.MaxSize == 0 {
			config.File.MaxSize = 10
		}
		if config.File.MaxBackups == 0 {
			config.File.MaxBackups = 5
		}
		if config.File.MaxAge == 0 {
			config.File.MaxAge = 7
		}
		fileWriter := zapcore.AddSync(&lumberjack.Logger{
			Filename:   config.File.Filename,
			MaxSize:    config.File.MaxSize,
			MaxBackups: config.File.MaxBackups,
			MaxAge:     config.File.MaxAge,
			Compress:   config.File.Compress,
		})
		writers = append(writers, fileWriter)
	}

	logLevel, err := zapcore.ParseLevel(config.Level)
	if err != nil {
		panic(err)
	}
	traceLevel, err := zapcore.ParseLevel(config.TraceLevel)
	if err != nil {
		traceLevel = zapcore.ErrorLevel
	}

	baseCore := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(writers...), logLevel)

	var cores []zapcore.Core
	cores = append(cores, baseCore)

	for _, sink := range sinks {
		cores = append(cores, &sinkCore{
			sink:     sink,
			minLevel: logLevel,
		})
	}

	tee := zapcore.NewTee(cores...)
	zapLogger := zap.New(tee, zap.AddStacktrace(traceLevel))

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

type sinkCore struct {
	sink     LogSink
	minLevel zapcore.Level
}

func (c *sinkCore) Enabled(lvl zapcore.Level) bool {
	return lvl >= c.minLevel
}

func (c *sinkCore) With(fields []zapcore.Field) zapcore.Core {
	return c
}

func (c *sinkCore) Check(entry zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(entry.Level) {
		return ce.AddCore(entry, c)
	}
	return ce
}

func (c *sinkCore) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	// Convert zap fields to map
	// fieldMap := make(map[string]interface{})
	// for _, field := range fields {
	// 	switch field.Type {
	// 	case zapcore.StringType:
	// 		fieldMap[field.Key] = field.String
	// 	case zapcore.Int64Type, zapcore.Int32Type, zapcore.Int16Type, zapcore.Int8Type:
	// 		fieldMap[field.Key] = field.Integer
	// 	case zapcore.Float64Type:
	// 		fieldMap[field.Key] = field.Interface
	// 	case zapcore.BoolType:
	// 		fieldMap[field.Key] = field.Integer == 1
	// 	case zapcore.TimeType:
	// 		fieldMap[field.Key] = field.Interface.(time.Time)
	// 	default:
	// 		fieldMap[field.Key] = field.Interface
	// 	}
	// }

	// Create LogEntry
	logEntry := &LogEntry{
		Level:      int(entry.Level),
		Time:       entry.Time,
		LoggerName: entry.LoggerName,
		Message:    entry.Message,
		Caller:     entry.Caller.String(),
		Stack:      entry.Stack,
		// Fields:     fieldMap,
	}

	// Write to sink
	c.sink.WriteLog(logEntry)
	return nil
}

func (c *sinkCore) Sync() error {
	return nil
}
