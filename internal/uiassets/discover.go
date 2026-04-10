package uiassets

import (
	"os"
	"path/filepath"
)

const EnvName = "DAE_WEBUI_DIR"

func EnsureEnv() {
	if os.Getenv(EnvName) != "" {
		return
	}
	if dir := Discover(); dir != "" {
		_ = os.Setenv(EnvName, dir)
	}
}

func Discover() string {
	candidates := make([]string, 0, 8)
	addWithParents := func(base string) {
		if base == "" {
			return
		}
		current := filepath.Clean(base)
		for range 5 {
			candidates = append(candidates, filepath.Join(current, "ui"))
			parent := filepath.Dir(current)
			if parent == current {
				break
			}
			current = parent
		}
	}

	if cwd, err := os.Getwd(); err == nil {
		addWithParents(cwd)
	}
	if exePath, err := os.Executable(); err == nil {
		addWithParents(filepath.Dir(exePath))
	}

	seen := map[string]struct{}{}
	for _, dir := range candidates {
		clean := filepath.Clean(dir)
		if _, ok := seen[clean]; ok {
			continue
		}
		seen[clean] = struct{}{}
		if isWebUIDir(clean) {
			return clean
		}
	}
	return ""
}

func isWebUIDir(dir string) bool {
	required := []string{"index.html", "styles.css", "script.js"}
	for _, name := range required {
		info, err := os.Stat(filepath.Join(dir, name))
		if err != nil || info.IsDir() {
			return false
		}
	}
	return true
}
