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

type recordSink struct{ got []*LogEntry }

func (s *recordSink) WriteLog(e *LogEntry) { s.got = append(s.got, e) }

func TestWriteEntry_LevelThresholdAndFields(t *testing.T) {
	t.Cleanup(func() { Init(&LogConfig{Level: "debug", TraceLevel: "debug"}) })
	s := &recordSink{}
	Init(&LogConfig{Level: "info"}, s)

	WriteEntry(&LogEntry{Level: LogLevelMap[LogLevelDebug], Message: "debug 被门槛挡掉", Flags: 1})
	WriteEntry(&LogEntry{Level: LogLevelMap[LogLevelWarn], Message: "warn 通过", Flags: 1, Payload: `{"id":"x"}`})

	if len(s.got) != 1 {
		t.Fatalf("sink 收到 %d 条,期望 1 条(Debug 须被 Level=info 挡掉)", len(s.got))
	}
	e := s.got[0]
	if e.Message != "warn 通过" || e.Flags != 1 || e.Payload != `{"id":"x"}` || e.Time.IsZero() {
		t.Errorf("条目不对: %+v", e)
	}
}

func TestWriteEntry_DefaultLevelIsInfo(t *testing.T) {
	t.Cleanup(func() { Init(&LogConfig{Level: "debug", TraceLevel: "debug"}) })
	s := &recordSink{}
	Init(&LogConfig{}, s)

	WriteEntry(&LogEntry{Level: LogLevelMap[LogLevelDebug], Message: "d"})
	WriteEntry(&LogEntry{Level: LogLevelMap[LogLevelInfo], Message: "i"})

	if len(s.got) != 1 || s.got[0].Message != "i" {
		t.Fatalf("没配 Level 时应同 golog 缺省(info): %+v", s.got)
	}
}

// func TestFatalf(t *testing.T) {
// 	Fatalf("test fatal message: %s", "formatted")
// }
