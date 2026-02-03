package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gvns/gh-repo-defaults/internal/config"
)

// PrepareBoilerplate resolves templates and writes them into targetDir.
// Returns the list of destination paths that were written.
func PrepareBoilerplate(cfg config.BoilerplateConfig, targetDir string, userTemplateDir string) ([]string, error) {
	absTarget, err := filepath.EvalSymlinks(targetDir)
	if err != nil {
		return nil, fmt.Errorf("resolving target dir: %w", err)
	}

	var written []string

	for _, f := range cfg.Files {
		// Validate dest path
		if strings.Contains(f.Dest, "..") {
			return nil, fmt.Errorf("invalid destination path: %q", f.Dest)
		}
		destPath := filepath.Join(absTarget, f.Dest)
		absDest, err := filepath.Abs(destPath)
		if err != nil {
			return nil, fmt.Errorf("resolving dest path: %w", err)
		}
		if !strings.HasPrefix(absDest, absTarget+string(filepath.Separator)) && absDest != absTarget {
			return nil, fmt.Errorf("destination %q escapes target directory", f.Dest)
		}

		// Resolve template content
		content, err := ResolveTemplate(f.Src, userTemplateDir)
		if err != nil {
			return nil, fmt.Errorf("resolving template for %q: %w", f.Dest, err)
		}

		// Create parent directories and write file
		if err := os.MkdirAll(filepath.Dir(absDest), 0755); err != nil {
			return nil, fmt.Errorf("creating directory for %q: %w", f.Dest, err)
		}
		if err := os.WriteFile(absDest, content, 0644); err != nil {
			return nil, fmt.Errorf("writing %q: %w", f.Dest, err)
		}
		written = append(written, f.Dest)
	}

	return written, nil
}
