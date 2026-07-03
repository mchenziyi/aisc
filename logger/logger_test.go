package logger

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoggerRoundTrip(t *testing.T) {
	tmp := t.TempDir()

	log, err := New(tmp, "backend")
	if err != nil {
		t.Fatal("New:", err)
	}

	// 模拟一个 Stage 的完整生命周期
	log.Info("stage_start")

	log.Log(INFO, "round_start", 0, F{"round": 1, "type": "全量深度评审"})

	// draft 阶段 + tool calls
	draftLog := log.With("draft", "backend")
	draftLog.Debug("tool_call", F{"tool": "read_file", "path": "docs/prd-frozen.md"})
	draftLog.Debug("tool_result", F{"tool": "read_file", "result": "ok"})
	draftLog.Debug("tool_call", F{"tool": "write_file", "path": "backend/main.go"})
	draftLog.Debug("tool_result", F{"tool": "write_file", "result": "ok"})
	draftLog.Debug("tool_call", F{"tool": "run_shell", "cmd": "go build ./..."})
	draftLog.Debug("tool_result", F{"tool": "run_shell", "result": "exit 0"})
	log.Log(INFO, "draft", 45000, nil)

	// review round
	log.Log(INFO, "round_start", 0, F{"round": 1, "type": "全量深度评审"})
	log.Log(INFO, "decision", 0, F{"type": "revise", "action_items": 5})
	log.Log(INFO, "review_round", 32000, nil)

	// revise
	reviseLog := log.With("revise", "backend")
	reviseLog.Debug("tool_call", F{"tool": "read_file", "path": "backend/internal/todo/handler.go"})
	reviseLog.Debug("tool_call", F{"tool": "write_file", "path": "backend/internal/todo/handler.go"})
	log.Log(INFO, "revise", 12000, nil)

	// round 2 → freeze
	log.Log(INFO, "round_start", 0, F{"round": 2, "type": "定向复核"})
	log.Log(INFO, "decision", 0, F{"type": "freeze", "action_items": 0})
	log.Log(INFO, "review_round", 18000, nil)

	log.Info("stage_frozen")
	log.Close()

	// 读取并验证
	entries, err := os.ReadDir(filepath.Join(tmp, ".aisc", "logs"))
	if err != nil {
		t.Fatal("no log dir:", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 log file, got %d", len(entries))
	}

	data, err := os.ReadFile(filepath.Join(tmp, ".aisc", "logs", entries[0].Name()))
	if err != nil {
		t.Fatal("read log:", err)
	}

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) < 10 {
		t.Fatalf("expected >=10 log lines, got %d", len(lines))
	}

	// 验证每条都是有效 JSON
	var toolCallCount int
	for _, line := range lines {
		var e Entry
		if err := json.Unmarshal([]byte(line), &e); err != nil {
			t.Fatalf("invalid JSON: %s", line)
		}
		if e.TS == "" {
			t.Error("missing ts")
		}
		if e.Level == "" {
			t.Error("missing level")
		}
		if e.Msg == "" {
			t.Error("missing msg")
		}
		if e.Msg == "tool_call" {
			toolCallCount++
		}
	}

	if toolCallCount < 5 {
		t.Errorf("expected >=5 tool_call entries, got %d", toolCallCount)
	}

	t.Logf("✅ %d log lines, %d tool calls, all valid JSON", len(lines), toolCallCount)
}

func TestLoggerWithContext(t *testing.T) {
	tmp := t.TempDir()
	log, _ := New(tmp, "api-design")
	defer log.Close()

	sub := log.With("draft", "tech-lead")
	sub.Info("generating")

	log.Close()

	data, _ := os.ReadFile(filepath.Join(tmp, ".aisc", "logs", mustGlob(t, tmp)))
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	var e Entry
	json.Unmarshal([]byte(lines[0]), &e)

	if e.Stage != "api-design" {
		t.Errorf("stage mismatch: %s", e.Stage)
	}
	if e.Step != "draft" {
		t.Errorf("step mismatch: %s", e.Step)
	}
	if e.Agent != "tech-lead" {
		t.Errorf("agent mismatch: %s", e.Agent)
	}
	t.Logf("✅ stage=%s step=%s agent=%s", e.Stage, e.Step, e.Agent)
}

func TestLoggerTimed(t *testing.T) {
	tmp := t.TempDir()
	log, _ := New(tmp, "request")
	defer log.Close()

	func() {
		defer log.Timed("draft")()
	}()

	log.Close()

	data, _ := os.ReadFile(filepath.Join(tmp, ".aisc", "logs", mustGlob(t, tmp)))
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")

	var e Entry
	json.Unmarshal([]byte(lines[0]), &e)
	if e.Msg != "draft" {
		t.Errorf("msg mismatch: %s", e.Msg)
	}
	// dur_ms >= 0 is acceptable (may be 0 for sub-ms operations)
	if e.DurMs < 0 {
		t.Errorf("dur_ms should be >=0, got %d", e.DurMs)
	}
	t.Logf("✅ timed draft: %dms (0 is OK for fast ops)", e.DurMs)
}

func mustGlob(t *testing.T, root string) string {
	t.Helper()
	entries, _ := os.ReadDir(filepath.Join(root, ".aisc", "logs"))
	if len(entries) == 0 {
		t.Fatal("no log files")
	}
	return entries[0].Name()
}
