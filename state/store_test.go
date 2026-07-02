package state

import (
	"os"
	"path/filepath"
	"testing"
)

func TestStoreRoundTrip(t *testing.T) {
	root := t.TempDir()

	// 1. Project
	p := &Project{Name: "test", Description: "desc", Status: "active"}
	if err := SaveProject(root, p); err != nil {
		t.Fatal("SaveProject:", err)
	}
	p2, err := LoadProject(root)
	if err != nil {
		t.Fatal("LoadProject:", err)
	}
	if p2.Name != "test" {
		t.Fatalf("name mismatch: %q", p2.Name)
	}

	// 2. Stage
	s := &Stage{
		ID: "stage-requirement", Type: "Requirement", Status: "drafting",
		Order: 1, OwnerAgent: "pm-agent",
		ReviewerAgents: []string{"tech-lead", "qa"},
		ArtifactID: "prd", CurrentVersion: 1,
	}
	if err := SaveStage(root, s); err != nil {
		t.Fatal("SaveStage:", err)
	}
	s2, err := LoadStage(root, "stage-requirement")
	if err != nil {
		t.Fatal("LoadStage:", err)
	}
	if s2.Status != "drafting" || len(s2.ReviewerAgents) != 2 {
		t.Fatalf("stage mismatch: %+v", s2)
	}

	// 3. Artifact
	path, err := SaveArtifact(root, "stage-requirement", "prd", "# PRD Content", 1)
	if err != nil {
		t.Fatal("SaveArtifact:", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatal("artifact file missing:", path)
	}
	content, err := ReadArtifact(root, "prd")
	if err != nil {
		t.Fatal("ReadArtifact:", err)
	}
	if content != "# PRD Content" {
		t.Fatalf("content mismatch: %q", content)
	}

	// 4. Frozen PRD
	if err := SaveFrozenPRD(root, "frozen"); err != nil {
		t.Fatal("SaveFrozenPRD:", err)
	}

	// 5. Meeting
	m := &Meeting{
		ID: "meeting-001",
		Meta: MeetingMeta{
			ID: "meeting-001", Round: 1, ArtifactVersion: 1,
			Type: "requirement_review", Stage: "stage-requirement",
			TargetArtifact: "prd-v1.md", Moderator: "pm",
			Participants: "tl,qa", Status: "in_progress",
			CreatedAt: "2026-01-01T00:00:00Z",
		},
		Body:    "## Decision\n\nfreeze",
		Reviews: []Review{{AgentID: "qa", Content: "looks good"}},
	}
	if err := SaveMeeting(root, m); err != nil {
		t.Fatal("SaveMeeting:", err)
	}
	// 验证文件存在
	meetingPath := filepath.Join(root, DirMeetings, "01-requirement", "meeting-meeting-001.md")
	if _, err := os.Stat(meetingPath); err != nil {
		t.Fatal("meeting file missing:", meetingPath)
	}

	// 6. Memory
	mem := &Memory{
		Type:    "decision",
		Title:   "需求评审 meeting-001",
		Content: "review content",
		Relations: []Relation{{Type: "based_on", TargetType: "meeting", TargetID: "meeting-001"}},
		Tags:    []string{"需求评审", "PRD"},
	}
	if err := SaveMemory(root, "qa", "meeting-001-review", mem); err != nil {
		t.Fatal("SaveMemory:", err)
	}
	memPath := filepath.Join(root, DirMemory, "qa", "meeting-001-review.json")
	if _, err := os.Stat(memPath); err != nil {
		t.Fatal("memory file missing:", memPath)
	}

	// 7. Decision Memory
	if err := SaveDecisionMemory(root, map[string]any{"type": "revise"}); err != nil {
		t.Fatal("SaveDecisionMemory:", err)
	}
	dm, err := LoadDecisionMemory(root)
	if err != nil {
		t.Fatal("LoadDecisionMemory:", err)
	}
	if dm["type"] != "revise" {
		t.Fatalf("decision mismatch: %v", dm)
	}
	DeleteDecisionMemory(root)
	if _, err := LoadDecisionMemory(root); err == nil {
		t.Fatal("decision_memory should be deleted")
	}

	// 8. 验证目录结构
	for _, sub := range []string{
		filepath.Join(DirAISC, "stages", "01-requirement"),
		filepath.Join(DirAISC, "meetings", "01-requirement"),
		filepath.Join(DirAISC, "memory", "qa"),
		DirDocs,
	} {
		if _, err := os.Stat(filepath.Join(root, sub)); err != nil {
			t.Errorf("missing dir: %s", sub)
		}
	}
}
