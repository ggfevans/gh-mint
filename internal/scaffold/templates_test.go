package scaffold

import (
	"os"
	"testing"
)

func TestResolveTemplate_Embedded(t *testing.T) {
	content, err := ResolveTemplate("contributing.md", "")
	if err != nil {
		t.Fatalf("ResolveTemplate: %v", err)
	}
	if len(content) == 0 {
		t.Error("expected non-empty content")
	}
}

func TestResolveTemplate_MissingEmbedded(t *testing.T) {
	_, err := ResolveTemplate("nonexistent.txt", "")
	if err == nil {
		t.Error("expected error for missing template")
	}
}

func TestResolveTemplate_UserOverride(t *testing.T) {
	dir := t.TempDir()
	userFile := dir + "/contributing.md"
	if err := os.WriteFile(userFile, []byte("custom content"), 0644); err != nil {
		t.Fatal(err)
	}
	content, err := ResolveTemplate("contributing.md", dir)
	if err != nil {
		t.Fatalf("ResolveTemplate: %v", err)
	}
	if string(content) != "custom content" {
		t.Errorf("expected user override, got: %s", content)
	}
}

func TestResolveTemplate_PathTraversal(t *testing.T) {
	dir := t.TempDir()
	_, err := ResolveTemplate("../../etc/passwd", dir)
	if err == nil {
		t.Error("expected error for path traversal")
	}
}
