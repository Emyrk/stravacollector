package offlineserver

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/go-chi/chi/v5"
)

//go:embed static*
var staticFiles embed.FS

type Server struct {
	Router chi.Router
}

func New() (*Server, error) {
	s := &Server{}

	r, err := s.Routes()
	if err != nil {
		return nil, err
	}

	s.Router = r
	return s, nil
}

func (s *Server) Routes() (chi.Router, error) {
	r := chi.NewRouter()

	dir, err := fs.Sub(staticFiles, "static")
	if err != nil {
		return nil, err
	}

	// Serve all static files as last resort
	r.NotFound(http.FileServer(http.FS(dir)).ServeHTTP)

	return r, nil
}
