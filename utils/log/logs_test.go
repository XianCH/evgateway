package utils

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"
)

// mockWriter 用来替代 os.Stdout / 文件输出，捕获日志
type mockWriter struct {
	buf *bytes.Buffer
}

func (m *mockWriter) Write(p []byte) (int, error) {
	return m.buf.Write(p)
}

func TestLogger_LevelFiltering(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := &Logger{
		level:     INFO,
		stdLogger: NewStdLogger(&mockWriter{buf}),
		fileOut:   nil,
	}

	logger.Debug("this should be hidden")
	logger.Info("info message")

	out := buf.String()
	if strings.Contains(out, "DEBUG") {
		t.Errorf("DEBUG log should not appear, got: %s", out)
	}
	if !strings.Contains(out, "INFO") {
		t.Errorf("INFO log missing, got: %s", out)
	}
}

func TestLogger_WriteToFile(t *testing.T) {
	tmpFile := "test.log"
	defer os.Remove(tmpFile)

	logger := NewLogger(DEBUG, tmpFile)
	logger.Info("file log test")

	data, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("failed to read log file: %v", err)
	}

	if !strings.Contains(string(data), "file log test") {
		t.Errorf("log file does not contain expected message, got: %s", string(data))
	}
}

func TestLogger_FormatContainsTimestampAndCaller(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := &Logger{
		level:     DEBUG,
		stdLogger: NewStdLogger(&mockWriter{buf}),
		fileOut:   nil,
	}

	logger.Warn("format test")

	out := buf.String()
	if !strings.Contains(out, "WARN") {
		t.Errorf("expected WARN level in log, got: %s", out)
	}
	if !strings.Contains(out, ".go") {
		t.Errorf("expected caller file info in log, got: %s", out)
	}
	if !strings.Contains(out, "20") { // 粗略检查是否有年份时间戳
		t.Errorf("expected timestamp in log, got: %s", out)
	}
}

// 辅助：创建只输出到自定义 writer 的 log.Logger
func NewStdLogger(w *mockWriter) *log.Logger {
	return log.New(w, "", 0)
}
