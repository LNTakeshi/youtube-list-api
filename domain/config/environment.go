package config

import "os"

const ProjectID = "youtube-list-app-276208"

func IsLocal() bool {
	return os.Getenv("local") == "true"
}
