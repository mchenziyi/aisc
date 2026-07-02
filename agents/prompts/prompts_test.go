package prompts

import (
	"strings"
	"testing"
)

func TestLoadReviewer(t *testing.T) {
	// go test 的 cwd 在包目录下，需要调整路径
	SetDir(".")
	
	prompt, err := LoadReviewer("tech-lead", "exhaustive")
	if err != nil {
		t.Fatal("LoadReviewer:", err)
	}
	if !strings.Contains(prompt, "Tech Lead") {
		t.Error("missing role hint")
	}
	if !strings.Contains(prompt, "全量深度评审") {
		t.Error("missing exhaustive template")
	}
}

func TestLoadReviewerVerification(t *testing.T) {
	SetDir(".")
	prompt, err := LoadReviewer("qa", "verification")
	if err != nil {
		t.Fatal("LoadReviewer:", err)
	}
	if !strings.Contains(prompt, "QA Tester") {
		t.Error("missing QA role hint")
	}
	if !strings.Contains(prompt, "复核") {
		t.Error("missing verification template")
	}
}

func TestLoadPMMissing(t *testing.T) {
	SetDir(".")
	_, err := Load("pm", "nonexistent")
	if err == nil {
		t.Error("expected error for missing file")
	}
}
