package main

import (
	_ "embed"
	"io/fs"
	"net/http"
)

//go:embed templates/index.html
var indexHTML []byte

// StaticFileServer returns a handler for serving static files
func StaticFileServer() http.Handler {
	// Return an empty file server since we don't have static files
	return http.FileServer(http.FS(&emptyFileSystem{}))
}

// emptyFileSystem is a minimal implementation of fs.FS that returns errors
type emptyFileSystem struct{}

func (e *emptyFileSystem) Open(name string) (fs.File, error) {
	return nil, fs.ErrNotExist
}