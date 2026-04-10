package controllerapi

import (
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrg/xdg"
	"github.com/daeuniverse/dae/common/consts"
	embeddedui "github.com/qi-mooo/dae-next/ui"
)

const webUIDirEnv = "DAE_WEBUI_DIR"

type webUI struct {
	dir        string
	source     string
	fileServer http.Handler
}

func discoverWebUI() *webUI {
	return discoverWebUIWithDirs(webUIDirs())
}

func discoverWebUIWithDirs(dirs []string) *webUI {
	for _, dir := range dirs {
		if dir == "" {
			continue
		}
		if isWebUIDir(dir) {
			return newDiskWebUI(dir)
		}
	}
	return newEmbeddedWebUI()
}

func newDiskWebUI(dir string) *webUI {
	return &webUI{
		dir:        dir,
		source:     "disk",
		fileServer: http.FileServer(http.Dir(dir)),
	}
}

func newEmbeddedWebUI() *webUI {
	if !isWebUIFS(embeddedui.Files) {
		return nil
	}
	return &webUI{
		source:     "embedded",
		fileServer: http.FileServer(http.FS(embeddedui.Files)),
	}
}

func webUIDirs() []string {
	dirs := make([]string, 0, 8)
	seen := map[string]struct{}{}
	add := func(dir string) {
		if dir == "" {
			return
		}
		clean := filepath.Clean(dir)
		if _, ok := seen[clean]; ok {
			return
		}
		seen[clean] = struct{}{}
		dirs = append(dirs, clean)
	}
	addWithParents := func(base string) {
		if base == "" {
			return
		}
		current := filepath.Clean(base)
		for range 4 {
			add(filepath.Join(current, "ui"))
			parent := filepath.Dir(current)
			if parent == current {
				break
			}
			current = parent
		}
	}

	add(os.Getenv(webUIDirEnv))

	if cwd, err := os.Getwd(); err == nil {
		addWithParents(cwd)
	}

	if exePath, err := os.Executable(); err == nil {
		addWithParents(filepath.Dir(exePath))
	}

	appDir := consts.AppName
	add(filepath.Join(xdg.DataHome, appDir, "ui"))
	for _, base := range xdg.DataDirs {
		add(filepath.Join(base, appDir, "ui"))
	}
	add(filepath.Join("/usr/local/share", appDir, "ui"))
	add(filepath.Join("/usr/share", appDir, "ui"))

	return dirs
}

func isWebUIDir(dir string) bool {
	for _, name := range requiredWebUIFiles {
		info, err := os.Stat(filepath.Join(dir, name))
		if err != nil || info.IsDir() {
			return false
		}
	}
	return true
}

var requiredWebUIFiles = []string{"index.html", "styles.css", "script.js"}

func isWebUIFS(fsys fs.FS) bool {
	for _, name := range requiredWebUIFiles {
		info, err := fs.Stat(fsys, name)
		if err != nil || info.IsDir() {
			return false
		}
	}
	return true
}

func isWebUIPath(path string) bool {
	return path == "/ui" || path == "/ui/" || strings.HasPrefix(path, "/ui/")
}

func (u *webUI) serveHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store, max-age=0")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	if r.URL.Path == "/ui" {
		http.Redirect(w, r, "/ui/", http.StatusTemporaryRedirect)
		return
	}

	trimmed := strings.TrimPrefix(r.URL.Path, "/ui")
	if trimmed == "" {
		trimmed = "/"
	}

	clone := r.Clone(r.Context())
	clone.URL.Path = trimmed
	u.fileServer.ServeHTTP(w, clone)
}
