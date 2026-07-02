package orchestration

import (
	"strings"
	"testing"
)

func TestParseDecision_JSON(t *testing.T) {
	d := ParseDecision(`{"type":"revise","summary":"需要改","action_items":[{"description":"改分页"}]}`)
	if d.Type != "revise" {
		t.Fatalf("type mismatch: %q", d.Type)
	}
	if d.Summary != "需要改" {
		t.Fatalf("summary mismatch: %q", d.Summary)
	}
	if len(d.ActionItems) != 1 || d.ActionItems[0].Description != "改分页" {
		t.Fatalf("action_items mismatch: %+v", d.ActionItems)
	}
}

func TestParseDecision_MarkdownBlock(t *testing.T) {
	raw := "一些废话\n```json\n{\"type\":\"freeze\",\"summary\":\"通过\"}\n```\n更多废话"
	d := ParseDecision(raw)
	if d.Type != "freeze" {
		t.Fatalf("type mismatch: %q", d.Type)
	}
}

func TestParseDecision_Embedded(t *testing.T) {
	d := ParseDecision(`好的，决策如下：{"type":"adopt","summary":"ok"}，完毕`)
	if d.Type != "adopt" {
		t.Fatalf("type mismatch: %q", d.Type)
	}
}

func TestParseDecision_Fallback(t *testing.T) {
	d := ParseDecision("[DECISION]\ntype: revise\nsummary: 需要改分页")
	if d.Type != "revise" {
		t.Fatalf("type mismatch: %q", d.Type)
	}
}

func TestParseDecision_Empty(t *testing.T) {
	d := ParseDecision("不知道说什么")
	if d.Type != "unknown" {
		t.Fatalf("expected unknown, got %q", d.Type)
	}
}

func TestActionItemsText(t *testing.T) {
	items := []ActionItem{
		{Description: "改用 cursor 分页"},
		{Description: "补充 429 错误码"},
	}
	result := ActionItemsText(items)
	if !strings.Contains(result, "1. 改用 cursor") {
		t.Error("missing item 1")
	}
	if !strings.Contains(result, "2. 补充 429") {
		t.Error("missing item 2")
	}
}
