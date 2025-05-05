package xlog

import (
	"sync"
	"testing"
	"time"

	"go.uber.org/zap/zapcore"
)

// memorySink stores logs in memory for testing and debugging purposes
type memorySink struct {
	mu    sync.RWMutex
	logs  []*LogEntry
	limit int // Maximum number of logs to store
}

// newMemorySink creates a new memorySink with the specified log limit
func newMemorySink(limit int) *memorySink {
	return &memorySink{
		logs:  make([]*LogEntry, 0, limit),
		limit: limit,
	}
}

// WriteLog implements the LogSink interface
func (s *memorySink) WriteLog(entry *LogEntry) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// If we've reached the limit, remove the oldest log
	if len(s.logs) >= s.limit {
		s.logs = s.logs[1:]
	}

	// Add the new log entry
	s.logs = append(s.logs, entry)
}

// getLogs returns all stored logs
func (s *memorySink) getLogs() []*LogEntry {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Return a copy of the logs
	logs := make([]*LogEntry, len(s.logs))
	copy(logs, s.logs)
	return logs
}

// clear removes all stored logs
func (s *memorySink) clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.logs = s.logs[:0]
}

// getLogCount returns the number of stored logs
func (s *memorySink) getLogCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.logs)
}

func TestMemorySink(t *testing.T) {
	// Create a new memory sink with a limit of 5 logs
	sink := newMemorySink(5)

	// Test initial state
	if sink.getLogCount() != 0 {
		t.Errorf("Expected 0 logs, got %d", sink.getLogCount())
	}

	// Create a test log entry
	entry := &LogEntry{
		Level:      int(zapcore.InfoLevel),
		Time:       time.Now(),
		LoggerName: "test",
		Message:    "Test message",
		// Fields:     map[string]interface{}{"key": "value"},
	}

	// Write a log entry
	sink.WriteLog(entry)

	// Test log count
	if sink.getLogCount() != 1 {
		t.Errorf("Expected 1 log, got %d", sink.getLogCount())
	}

	// Test retrieving logs
	logs := sink.getLogs()
	if len(logs) != 1 {
		t.Errorf("Expected 1 log, got %d", len(logs))
	}

	// Test log content
	if logs[0].Message != "Test message" {
		t.Errorf("Expected message 'Test message', got '%s'", logs[0].Message)
	}

	// Test log limit
	for i := 0; i < 10; i++ {
		sink.WriteLog(entry)
	}

	if sink.getLogCount() != 5 {
		t.Errorf("Expected 5 logs after limit, got %d", sink.getLogCount())
	}

	// Test clearing logs
	sink.clear()
	if sink.getLogCount() != 0 {
		t.Errorf("Expected 0 logs after clear, got %d", sink.getLogCount())
	}
}

func TestMemorySinkWithLogger(t *testing.T) {
	// Create a memory sink
	sink := newMemorySink(10)

	// Initialize logger with the memory sink
	config := &LogConfig{
		Level: "debug",
	}
	Init(config, sink)

	// Write some test logs
	Debug("Debug message")
	Info("Info message")
	Warn("Warning message")
	Error("Error message")

	// Test log count
	if sink.getLogCount() != 4 {
		t.Errorf("Expected 4 logs, got %d", sink.getLogCount())
	}

	// Test log levels
	logs := sink.getLogs()
	expectedLevels := []int{
		int(zapcore.DebugLevel),
		int(zapcore.InfoLevel),
		int(zapcore.WarnLevel),
		int(zapcore.ErrorLevel),
	}

	for i, log := range logs {
		if log.Level != expectedLevels[i] {
			t.Errorf("Expected level %d, got %d", expectedLevels[i], log.Level)
		}
	}

	// Clean up
	Finalize()
}

// func TestLogFormat(t *testing.T) {
// 	// Create a memory sink
// 	sink := newMemorySink(10)

// 	// Initialize logger with the memory sink
// 	config := &LogConfig{
// 		Level: "debug",
// 	}
// 	Init(config, sink)

// 	// Write a test log with fields
// 	logger := _logger.(*ZapLogger).innerLogger
// 	logger.Info("Test log message", zap.String("key1", "value1"), zap.Int("key2", 123))

// 	// Get the log entry
// 	logs := sink.getLogs()
// 	if len(logs) != 1 {
// 		t.Fatalf("Expected 1 log entry, got %d", len(logs))
// 	}

// 	entry := logs[0]

// 	// Verify log entry fields
// 	if entry.Level != int(zapcore.InfoLevel) {
// 		t.Errorf("Expected level %d, got %d", zapcore.InfoLevel, entry.Level)
// 	}

// 	if entry.Message != "Test log message" {
// 		t.Errorf("Expected message 'Test log message', got '%s'", entry.Message)
// 	}

// 	// Verify fields
// 	if len(entry.Fields) != 2 {
// 		t.Errorf("Expected 2 fields, got %d", len(entry.Fields))
// 	}

// 	if entry.Fields["key1"] != "value1" {
// 		t.Errorf("Expected field key1='value1', got '%v'", entry.Fields["key1"])
// 	}

// 	if v, ok := entry.Fields["key2"].(int64); !ok || v != 123 {
// 		t.Errorf("Expected field key2=123 (int64), got '%v' (%T)", entry.Fields["key2"], entry.Fields["key2"])
// 	}

// 	// Clean up
// 	Finalize()
// }

func TestFormattedLogOutput(t *testing.T) {
	// Create a memory sink
	sink := newMemorySink(10)

	// Initialize logger with the memory sink
	config := &LogConfig{
		Level: "debug",
	}
	Init(config, sink)

	// Test different format strings
	Debugf("Number: %d, String: %s, Value: %v", 42, "test", true)
	Infof("Float: %.2f, Multiple: %d-%s", 3.14159, 100, "abc")
	Warnf("Complex: %+v, Quoted: %q", struct{ Name string }{"test"}, "quote me")

	// Get the log entries
	logs := sink.getLogs()
	if len(logs) != 3 {
		t.Fatalf("Expected 3 log entries, got %d", len(logs))
	}

	// Verify debug message
	if logs[0].Level != int(zapcore.DebugLevel) {
		t.Errorf("Expected level %d, got %d", zapcore.DebugLevel, logs[0].Level)
	}
	expectedDebug := "Number: 42, String: test, Value: true"
	if logs[0].Message != expectedDebug {
		t.Errorf("Expected message '%s', got '%s'", expectedDebug, logs[0].Message)
	}

	// Verify info message
	if logs[1].Level != int(zapcore.InfoLevel) {
		t.Errorf("Expected level %d, got %d", zapcore.InfoLevel, logs[1].Level)
	}
	expectedInfo := "Float: 3.14, Multiple: 100-abc"
	if logs[1].Message != expectedInfo {
		t.Errorf("Expected message '%s', got '%s'", expectedInfo, logs[1].Message)
	}

	// Verify warn message
	if logs[2].Level != int(zapcore.WarnLevel) {
		t.Errorf("Expected level %d, got %d", zapcore.WarnLevel, logs[2].Level)
	}
	expectedWarn := `Complex: {Name:test}, Quoted: "quote me"`
	if logs[2].Message != expectedWarn {
		t.Errorf("Expected message '%s', got '%s'", expectedWarn, logs[2].Message)
	}

	// Clean up
	Finalize()
}
