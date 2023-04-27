package site

import (
	"embed"
)

//go:embed strava-frontend/build
var StaticFiles embed.FS
