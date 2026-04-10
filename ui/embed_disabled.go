//go:build noembedui

package ui

import "io/fs"

var Enabled = false

var Files fs.FS = emptyFS{}

type emptyFS struct{}

func (emptyFS) Open(string) (fs.File, error) {
	return nil, fs.ErrNotExist
}
