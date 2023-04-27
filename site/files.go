package server

import (
	"embed"
	"io/fs"
	"log"
)

//go:embed strava-frontend/build/*
var staticFiles embed.FS

func FS() fs.FS {
	static, err := fs.Sub(fs.FS(staticFiles), "strava-frontend/build")
	if err != nil {
		log.Fatalf("failed to get static files: %w", err)
	}
	return static
}
