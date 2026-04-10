//go:build !noembedui

package ui

import (
	"embed"
	"io/fs"
)

var Enabled = true

//go:embed index.html styles.css script.js icons/*.svg
var embeddedFiles embed.FS

var Files fs.FS = embeddedFiles
