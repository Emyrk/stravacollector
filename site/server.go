package server

import (
	"io/fs"
	"net/http"
	"path"
	"strings"
)

type handler struct {
	fs  fs.FS
	mux *http.ServeMux
}

func Handler(siteFS fs.FS) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.FS(siteFS)))

	return &handler{
		fs:  siteFS,
		mux: mux,
	}
}

func (h *handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	// reqFile is the static file requested
	reqFile := filePath(req.URL.Path)

	// If the original file exists, serve it
	if h.exists(reqFile) {
		h.mux.ServeHTTP(resp, req)
		return
	}

	// Serve the file assuming it's an html file
	// This matches paths like `/app/terminal.html`
	req.URL.Path = strings.TrimSuffix(req.URL.Path, "/")
	req.URL.Path += ".html"
	reqFile = filePath(req.URL.Path)
	if h.exists(reqFile) {
		h.mux.ServeHTTP(resp, req)
		return
	}

	req.URL.Path = "/"

	// This will send a correct 404
	h.mux.ServeHTTP(resp, req)
}

func (h *handler) exists(filePath string) bool {
	f, err := h.fs.Open(filePath)
	if err == nil {
		_ = f.Close()
	}
	return err == nil
}

// filePath returns the filepath of the requested file.
func filePath(p string) string {
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}
	return strings.TrimPrefix(path.Clean(p), "/")
}
