package config

import (
	"cloud.google.com/go/compute/metadata"
	"encoding/json"
	"os"
)

var ProjectID = "local"

func init() {
	if os.Getenv("PROJECT_ID") != "" {
		ProjectID = os.Getenv("PROJECT_ID")
		println("Override ProjectID: " + ProjectID)
		return
	}

	if IsLocal() {
		b, err := os.ReadFile("/home/application_default_credentials.json")
		if err != nil {
			return
		}
		var i any
		err = json.Unmarshal(b, &i)
		if err != nil {
			return
		}
		id, ok := i.(map[string]any)["quota_project_id"]
		if !ok {
			return
		}
		ProjectID = id.(string)
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
