package assets

import (
	"embed"
	_ "embed"
)

//go:embed assets/*
var assets embed.FS
