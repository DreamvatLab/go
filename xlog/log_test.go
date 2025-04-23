package xlog

import (
	"testing"
)

func init() {
	logConfig := &LogConfig{
		Level:      "debug",
		TraceLevel: "debug",
		File: &FileLogConfig{
			Filename: "test.log",
		},
	}

	Init(logConfig)
}

func TestDebug(t *testing.T) {
	Debug("test debug message")
}

func TestDebugf(t *testing.T) {
	Debugf("test debug message: %s", "formatted")
}

func TestInfo(t *testing.T) {
	Info("test info message")
}

func TestInfof(t *testing.T) {
	Infof("test info message: %s", "formatted")
}

func TestWarn(t *testing.T) {
	Warn("test warn message")
}

func TestWarnf(t *testing.T) {
	Warnf("test warn message: %s", "formatted")
}

func TestError(t *testing.T) {
	Error("test error message")
}

func TestErrorf(t *testing.T) {
	Errorf("test error message: %s", "formatted")
}

func TestFatal(t *testing.T) {
	Fatal("test fatal message")
}

func TestFatalf(t *testing.T) {
	Fatalf("test fatal message: %s", "formatted")
}
