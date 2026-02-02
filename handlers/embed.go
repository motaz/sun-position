package handlers

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed templates/index.html
var indexHTML []byte

//go:embed all:static
var staticFiles embed.FS

// StaticFileServer returns a handler for serving static files
func StaticFileServer() http.Handler {
	subFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		// Fallback to empty filesystem if static dir doesn't exist
		return http.FileServer(http.FS(&emptyFileSystem{}))
	}
	return http.FileServer(http.FS(subFS))
}

// emptyFileSystem is a minimal implementation of fs.FS that returns errors
type emptyFileSystem struct{}

func (e *emptyFileSystem) Open(name string) (fs.File, error) {
	return nil, fs.ErrNotExist
}