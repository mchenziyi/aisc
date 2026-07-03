// Package logger 提供 Stage 执行过程的 JSON 行日志。
// 日志文件输出到 .aisc/logs/{timestamp}-{stage}.log。
package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Level 日志级别
type Level string

const (
	INFO  Level = "INFO"
	DEBUG Level = "DEBUG"
	ERROR Level = "ERROR"
)

// Entry 一条日志记录
type Entry struct {
	TS    string `json:"ts"`
	Level Level  `json:"level"`
	Stage string `json:"stage,omitempty"`
	Step  string `json:"step,omitempty"`
	Agent string `json:"agent,omitempty"`
	Msg   string `json:"msg"`
	DurMs int64  `json:"dur_ms,omitempty"`
	Extra any    `json:"extra,omitempty"`
}

// Logger 写 JSON 行日志到文件。
// 通过 With() 创建子 Logger 时共享同一个 mutex，保证并发安全。
type Logger struct {
	mu   *sync.Mutex
	w    io.WriteCloser
	base Entry
}

// New 创建 Logger，日志写到 .aisc/logs/{stage}.log
func New(root, stage string) (*Logger, error) {
	dir := filepath.Join(root, ".aisc", "logs")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	ts := time.Now().UTC().Format("20060102-150405")
	name := fmt.Sprintf("%s-%s.log", ts, stage)
	path := filepath.Join(dir, name)
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	return &Logger{
		mu:   &sync.Mutex{},
		w:    f,
		base: Entry{Stage: stage},
	}, nil
}

// With 返回共享同一 writer 和 mutex 的子 Logger。
func (l *Logger) With(step, agent string) *Logger {
	return &Logger{
		mu:   l.mu,
		w:    l.w,
		base: Entry{Stage: l.base.Stage, Step: step, Agent: agent},
	}
}

// F 是 Logger 方法的便捷 extra 参数类型
type F map[string]any

// Log 写入一条日志。
func (l *Logger) Log(level Level, msg string, durMs int64, extra any) {
	e := l.base
	e.TS = time.Now().UTC().Format(time.RFC3339Nano)
	e.Level = level
	e.Msg = msg
	e.DurMs = durMs
	e.Extra = extra
	data, _ := json.Marshal(e)
	l.mu.Lock()
	l.w.Write(append(data, '\n'))
	l.mu.Unlock()
}

// Info 写入 INFO 级别日志
func (l *Logger) Info(msg string) { l.Log(INFO, msg, 0, nil) }

// Debug 写入 DEBUG 级别日志
func (l *Logger) Debug(msg string, extra any) { l.Log(DEBUG, msg, 0, extra) }

// Error 写入 ERROR 级别日志
func (l *Logger) Error(msg string, extra any) { l.Log(ERROR, msg, 0, extra) }

// Timed 返回一个函数，调用时记录耗时。
// 用法: defer l.Timed("draft")()
func (l *Logger) Timed(step string) func() {
	start := time.Now()
	return func() {
		l.Log(INFO, step, time.Since(start).Milliseconds(), nil)
	}
}

// TimeBlock 在 fn 执行前后自动记录耗时至日志。
func (l *Logger) TimeBlock(step string, fn func()) {
	defer l.Timed(step)()
	fn()
}

// Close 关闭日志文件
func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.w.Close()
}
