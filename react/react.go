package react

import (
	"embed"
	"io/fs"
)

//go:embed build/*
var f embed.FS

func Serve() fs.FS {
	file, err := fs.Sub(f, "build")
	if err != nil {
		panic(err)
	}

	return file
}
