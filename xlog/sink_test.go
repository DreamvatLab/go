package xlog

import (
	"bytes"
	"testing"
)

// MemorySink is a custom sink that stores logs in memory
type MemorySink struct {
	buffer *bytes.Buffer
}

// NewMemorySink creates a new MemorySink instance
func NewMemorySink() *MemorySink {
	return &MemorySink{
		buffer: bytes.NewBuffer(nil),
	}
}

// WriteLog implements the LogSink interface
func (s *MemorySink) WriteLog(level string, v ...interface{}) {
	s.buffer.WriteString(level + ": ")
	for _, arg := range v {
		s.buffer.WriteString(arg.(string))
	}
	s.buffer.WriteString("\n")
}

// GetLogs returns the stored logs as a string
func (s *MemorySink) GetLogs() string {
	return s.buffer.String()
}

// Clear clears the stored logs
func (s *MemorySink) Clear() {
	s.buffer.Reset()
}

func TestCustomSink(t *testing.T) {
	// Create a new memory sink
	sink := NewMemorySink()

	// Initialize logger with the custom sink
	logConfig := &LogConfig{
		Level:      "debug",
		TraceLevel: "debug",
	}
	Init(logConfig, sink)

	// Test different log levels
	Debug("test debug message")
	Info("test info message")
	Warn("test warn message")
	Error("test error message")

	// Get and verify logs
	logs := sink.GetLogs()
	expectedLogs := []string{
		"debug: test debug message",
		"info: test info message",
		"warn: test warn message",
		"error: test error message",
	}

	for _, expected := range expectedLogs {
		if !bytes.Contains([]byte(logs), []byte(expected)) {
			t.Errorf("Expected log not found: %s", expected)
		}
	}

	// Clear logs and test again
	sink.Clear()
	Debug("new debug message")
	logs = sink.GetLogs()
	if !bytes.Contains([]byte(logs), []byte("debug: new debug message")) {
		t.Error("Log clearing and new logging failed")
	}
}

func TestFormattedLogging(t *testing.T) {
	// Create a new memory sink
	sink := NewMemorySink()

	// Initialize logger with the custom sink
	logConfig := &LogConfig{
		Level:      "debug",
		TraceLevel: "debug",
	}
	Init(logConfig, sink)

	// Test formatted logging with different types
	Debugf("debug message with number: %d", 42)
	Infof("info message with string: %s", "hello")
	Warnf("warn message with float: %.2f", 3.14159)
	Errorf("error message with multiple args: %s %d %.2f", "test", 123, 45.67)

	// Get and verify logs
	logs := sink.GetLogs()
	expectedLogs := []string{
		"debug: debug message with number: 42",
		"info: info message with string: hello",
		"warn: warn message with float: 3.14",
		"error: error message with multiple args: test 123 45.67",
	}

	for _, expected := range expectedLogs {
		if !bytes.Contains([]byte(logs), []byte(expected)) {
			t.Errorf("Expected formatted log not found: %s", expected)
		}
	}

	// Test complex formatting
	sink.Clear()
	Debugf("complex format: %v %#v %T", []int{1, 2, 3}, []int{1, 2, 3}, []int{})
	logs = sink.GetLogs()
	if !bytes.Contains([]byte(logs), []byte("complex format: [1 2 3] []int{1, 2, 3} []int")) {
		t.Error("Complex formatting test failed")
	}
}
