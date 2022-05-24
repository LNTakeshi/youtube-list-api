package config

import (
	"cloud.google.com/go/compute/metadata"
	"os"
)

var ProjectID = "local"

func init() {
	if IsLocal() {
		return
	}

	if os.Getenv("PROJECT_ID") != "" {
		ProjectID = os.Getenv("PROJECT_ID")
		println("Override ProjectID: " + ProjectID)
		return
	}

	id, err := metadata.ProjectID()
	if err != nil {
		panic(err)
	}
	ProjectID = id
}

func IsLocal() bool {
	return os.Getenv("local") == "true" && os.Getenv("PROJECT_ID") == ""
}
