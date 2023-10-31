//go:build !static
// +build !static

package server

import (
	"embed"
	"io/fs"
	"log"
)

//go:embed slim/*
var staticFiles embed.FS

func FS() fs.FS {
	static, err := fs.Sub(fs.FS(staticFiles), "slim")
	if err != nil {
		log.Fatalf("failed to get static files: %s", err.Error())
	}
	return static
}

//go:embed staticpages/Login429.html
var LoginFailed string
