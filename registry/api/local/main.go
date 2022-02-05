package main

import (
	"os"
	"youtubelist/registry/api"
)

func main() {
	os.Setenv("local", "true")
	api.Start()
}
