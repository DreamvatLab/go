package xlog

import (
	"fmt"
	"time"

	"github.com/DreamvatLab/go/xtask"
	"github.com/kataras/golog"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
)

// shouldShowCaller 判断是否应该显示调用者信息
func shouldShowCaller(detailLevel int, currentLevel string) bool {
	currentLevelNum, exists := LogLevelMap[currentLevel]
	if !exists {
		return false
	}

	return detailLevel <= currentLevelNum
}

// newGologLogger creates a new GologLogger instance with the provided configuration
func newGologLogger(config *LogConfig, sinks ...LogSink) ILogger {
	if config.TraceLevel == "" {
		config.TraceLevel = LogLevelError
	}

	logger := golog.New()
	if config.Level != "" {
		logger.SetLevel(config.Level)
	}
	logger.SetTimeFormat("2006/01/02 15:04:05")

	// 文件输出（控制台输出用golog默认的，不要AddOutput(os.Stdout)）
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
		rotationTime := time.Duration(24) * time.Hour
		writer, err := rotatelogs.New(
			config.File.Filename+".%Y%m%d%H%M%S",
			rotatelogs.WithRotationTime(rotationTime),
			rotatelogs.WithRotationCount(uint(config.File.MaxBackups)),
		)
		if err == nil {
			logger.AddOutput(writer)
		}
	}

	// sink逻辑
	if len(sinks) > 0 {
		logger.Handle(func(log *golog.Log) bool {
			logEntry := &LogEntry{
				Level:      convertGologLevel(log.Level),
				Time:       time.Now(),
				LoggerName: "golog",
				Message:    log.Message,
				Caller:     "",
				Stack:      "",
			}

			xtask.ParallelRunSlice(len(sinks), sinks, func(sink LogSink) (interface{}, error) {
				sink.WriteLog(logEntry)
				return nil, nil
			})

			return false // 让golog正常输出
		})
	}

	// 解析 TraceLevel
	_detailLevel := LogLevelMap[config.TraceLevel]

	return &GologLogger{
		config:      config,
		innerLogger: logger,
		sinks:       sinks,
		detailLevel: _detailLevel,
	}
}

// GologLogger implements the ILogger interface using the kataras/golog library
type GologLogger struct {
	config      *LogConfig
	innerLogger *golog.Logger
	sinks       []LogSink
	detailLevel int
}

func (o *GologLogger) SetConfig(config *LogConfig) {
	o.config = config
	if config.Level != "" {
		o.innerLogger.SetLevel(config.Level)
	}
	// 更新 detailLevel
	if config.TraceLevel != "" {
		o.detailLevel = LogLevelMap[config.TraceLevel]
	}
}

func (o *GologLogger) GetConfig() *LogConfig {
	return o.config
}

func (o *GologLogger) getCallerInfo() string {
	// _, file, line, ok := runtime.Caller(5)
	// if ok {
	// 	return fmt.Sprintf("<%s:%d>", file, line)
	// }
	// return ""

	// pcs := make([]uintptr, 32)   // 增大数量可获取更深层栈
	// n := runtime.Callers(3, pcs) // skip: runtime.Callers, getCallerInfo, Debug
	// frames := runtime.CallersFrames(pcs[:n])

	// var sb strings.Builder
	// for {
	// 	frame, more := frames.Next()
	// 	file := strings.ToLower(frame.File)

	// 	// ❌ 过滤日志库和标准库路径（根据实际情况调整关键词）
	// 	if strings.Contains(file, "/dreamvatlab/") ||
	// 		strings.Contains(file, "runtime/") ||
	// 		strings.Contains(file, "program files/go") {
	// 		if !more {
	// 			break
	// 		}
	// 		continue
	// 	}

	// 	sb.WriteString(fmt.Sprintf("%s:%d %s", frame.File, frame.Line, frame.Function))
	// 	if !more {
	// 		break
	// 	}
	// }
	// return sb.String()
	return "# " // 暂不显示调用者信息，因为 golog配合 erros.WithStack 时，%+v 会打印出堆栈信息，导致重复
}

func (o *GologLogger) Debug(v ...interface{}) {
	if shouldShowCaller(o.detailLevel, LogLevelDebug) {
		caller := o.getCallerInfo()
		if caller != "" {
			o.innerLogger.Debugf("%s%+v", caller, v)
			return
		}
	}
	o.innerLogger.Debug(v...)
}

func (o *GologLogger) Debugf(format string, args ...interface{}) {
	if shouldShowCaller(o.detailLevel, LogLevelDebug) {
		caller := o.getCallerInfo()
		if caller != "" {
			o.innerLogger.Debugf("%s %s", caller, fmt.Sprintf(format, args...))
			return
		}
	}
	o.innerLogger.Debugf(format, args...)
}

func (o *GologLogger) Info(v ...interface{}) {
	if shouldShowCaller(o.detailLevel, LogLevelInfo) {
		caller := o.getCallerInfo()
		if caller != "" {
			o.innerLogger.Infof("%s%+v", caller, v)
			return
		}
	}
	o.innerLogger.Info(v...)
}

func (o *GologLogger) Infof(format string, args ...interface{}) {
	if shouldShowCaller(o.detailLevel, LogLevelInfo) {
		caller := o.getCallerInfo()
		if caller != "" {
			o.innerLogger.Infof("%s %s", caller, fmt.Sprintf(format, args...))
			return
		}
	}
	o.innerLogger.Infof(format, args...)
}

func (o *GologLogger) Warn(v ...interface{}) {
	if shouldShowCaller(o.detailLevel, LogLevelWarn) {
		caller := o.getCallerInfo()
		if caller != "" {
			o.innerLogger.Warnf("%s%+v", caller, v)
			return
		}
	}
	o.innerLogger.Warn(v...)
}

func (o *GologLogger) Warnf(format string, args ...interface{}) {
	if shouldShowCaller(o.detailLevel, LogLevelWarn) {
		caller := o.getCallerInfo()
		if caller != "" {
			o.innerLogger.Warnf("%s %s", caller, fmt.Sprintf(format, args...))
			return
		}
	}
	o.innerLogger.Warnf(format, args...)
}

func (o *GologLogger) Error(v ...interface{}) {
	if shouldShowCaller(o.detailLevel, LogLevelError) {
		caller := o.getCallerInfo()
		if caller != "" {
			o.innerLogger.Errorf("%s%+v", caller, v)
			return
		}
	}
	o.innerLogger.Error(v...)
}

func (o *GologLogger) Errorf(format string, args ...interface{}) {
	if shouldShowCaller(o.detailLevel, LogLevelError) {
		caller := o.getCallerInfo()
		if caller != "" {
			o.innerLogger.Errorf("%s %s", caller, fmt.Sprintf(format, args...))
			return
		}
	}
	o.innerLogger.Errorf(format, args...)
}

func (o *GologLogger) Fatal(v ...interface{}) {
	if shouldShowCaller(o.detailLevel, LogLevelFatal) {
		caller := o.getCallerInfo()
		if caller != "" {
			o.innerLogger.Fatalf("%s%+v", caller, v)
			return
		}
	}
	o.innerLogger.Fatal(v...)
}

func (o *GologLogger) Fatalf(format string, args ...interface{}) {
	if shouldShowCaller(o.detailLevel, LogLevelFatal) {
		caller := o.getCallerInfo()
		if caller != "" {
			o.innerLogger.Fatalf("%s %s", caller, fmt.Sprintf(format, args...))
			return
		}
	}
	o.innerLogger.Fatalf(format, args...)
}

func (o *GologLogger) Finalize() {
	// golog 不需要特殊的清理操作
}

func convertGologLevel(level golog.Level) int {
	switch level {
	case golog.DebugLevel:
		return 1
	case golog.InfoLevel:
		return 2
	case golog.WarnLevel:
		return 3
	case golog.ErrorLevel:
		return 4
	case golog.FatalLevel:
		return 5
	}
	return 2 // Info
}
