package reactddr

import (
	"embed"
	"io/fs"
)

//go:embed dist/*
var f embed.FS

func ServeDDR() fs.FS {
	file, err := fs.Sub(f, "dist")
	if err != nil {
		panic(err)
	}

	return file
}
