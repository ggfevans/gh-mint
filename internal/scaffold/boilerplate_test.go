package scaffold

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ggfevans/gh-mint/internal/config"
)

func TestPrepareBoilerplate(t *testing.T) {
	dir := t.TempDir()
	cfg := config.BoilerplateConfig{
		Files: []config.BoilerplateFile{
			{Src: "contributing.md", Dest: "CONTRIBUTING.md"},
			{Src: "ci.yml", Dest: ".github/workflows/ci.yml"},
		},
	}
	files, err := PrepareBoilerplate(cfg, dir, "")
	if err != nil {
		t.Fatalf("PrepareBoilerplate: %v", err)
	}
	if len(files) != 2 {
		t.Fatalf("expected 2 files, got %d", len(files))
	}

	contribPath := filepath.Join(dir, "CONTRIBUTING.md")
	if _, err := os.Stat(contribPath); os.IsNotExist(err) {
		t.Error("CONTRIBUTING.md not created")
	}

	ciPath := filepath.Join(dir, ".github", "workflows", "ci.yml")
	if _, err := os.Stat(ciPath); os.IsNotExist(err) {
		t.Error(".github/workflows/ci.yml not created")
	}
}

func TestPrepareBoilerplate_PathTraversalInDest(t *testing.T) {
	dir := t.TempDir()
	cfg := config.BoilerplateConfig{
		Files: []config.BoilerplateFile{
			{Src: "contributing.md", Dest: "../../etc/evil"},
		},
	}
	_, err := PrepareBoilerplate(cfg, dir, "")
	if err == nil {
		t.Error("expected error for path traversal in dest")
	}
}
