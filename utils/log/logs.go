package utils

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

var L = NewLogger(INFO, "")

// 日志级别
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

var levelNames = map[LogLevel]string{
	DEBUG: "DEBUG",
	INFO:  "INFO",
	WARN:  "WARN",
	ERROR: "ERROR",
}

var levelColors = map[LogLevel]string{
	DEBUG: "\033[36m", // 青色
	INFO:  "\033[32m", // 绿色
	WARN:  "\033[33m", // 黄色
	ERROR: "\033[31m", // 红色
}

const resetColor = "\033[0m"

type Logger struct {
	level     LogLevel
	stdLogger *log.Logger
	fileOut   io.Writer
}

// NewLogger 创建一个新的日志实例
func NewLogger(level LogLevel, logFile string) *Logger {
	var fileWriter io.Writer
	if logFile != "" {
		// 确保目录存在
		_ = os.MkdirAll(filepath.Dir(logFile), 0755)

		f, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err == nil {
			fileWriter = f
		} else {
			fmt.Printf("failed to open log file: %v\n", err)
		}
	}

	return &Logger{
		level:     level,
		stdLogger: log.New(os.Stdout, "", 0), // 自定义格式
		fileOut:   fileWriter,
	}
}

// SetLevel 设置全局日志级别
func (l *Logger) SetLevel(level LogLevel) {
	l.level = level
}

// logf 内部日志输出
func (l *Logger) logf(lv LogLevel, format string, v ...any) {
	if lv < l.level {
		return
	}

	// 时间
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	// 调用位置
	_, file, line, ok := runtime.Caller(2)
	caller := "???"
	if ok {
		caller = fmt.Sprintf("%s:%d", filepath.Base(file), line)
	}

	// 格式化
	msg := fmt.Sprintf(format, v...)
	colored := fmt.Sprintf("%s[%s][%s][%s]%s %s",
		levelColors[lv], timestamp, levelNames[lv], caller, resetColor, msg)

	// 输出到终端
	l.stdLogger.Println(colored)

	// 如果配置了文件日志，也写入
	if l.fileOut != nil {
		plain := fmt.Sprintf("[%s][%s][%s] %s\n", timestamp, levelNames[lv], caller, msg)
		_, _ = l.fileOut.Write([]byte(plain))
	}
}

func (l *Logger) Debug(format string, v ...any) { l.logf(DEBUG, format, v...) }
func (l *Logger) Info(format string, v ...any)  { l.logf(INFO, format, v...) }
func (l *Logger) Warn(format string, v ...any)  { l.logf(WARN, format, v...) }
func (l *Logger) Error(format string, v ...any) { l.logf(ERROR, format, v...) }

func Debug(format string, v ...any) {
	L.Debug(format, v...)
}

func Info(format string, v ...any) {
	L.Info(format, v...)
}

func Warn(format string, v ...any) {
	L.Warn(format, v...)
}

func Error(format string, v ...any) {
	L.Error(format, v...)
}
