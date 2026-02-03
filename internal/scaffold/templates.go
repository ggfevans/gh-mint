package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gvns/gh-repo-defaults/templates"
)

// ResolveTemplate returns the content of a template file.
// If userDir is non-empty and contains the file, it takes precedence.
// Otherwise falls back to embedded templates.
func ResolveTemplate(name string, userDir string) ([]byte, error) {
	// Block path traversal
	if strings.Contains(name, "..") {
		return nil, fmt.Errorf("invalid template path: %q", name)
	}

	// Try user override first
	if userDir != "" {
		userPath := filepath.Join(userDir, name)
		absUser, err := filepath.Abs(userPath)
		if err == nil {
			absDir, _ := filepath.Abs(userDir)
			if strings.HasPrefix(absUser, absDir+string(filepath.Separator)) {
				if data, err := os.ReadFile(absUser); err == nil {
					return data, nil
				}
			}
		}
	}

	// Fall back to embedded
	data, err := templates.FS.ReadFile(name)
	if err != nil {
		return nil, fmt.Errorf("template %q not found: %w", name, err)
	}
	return data, nil
}
