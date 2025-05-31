package scripts

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

func MergeSwaggerFiles(dir string, output string) error {
	merged := make(map[string]interface{})
	merged["swagger"] = "2.0"
	merged["paths"] = map[string]interface{}{}
	merged["definitions"] = map[string]interface{}{}
	merged["tags"] = []interface{}{}

	// Add Bearer token security definition
	merged["securityDefinitions"] = map[string]interface{}{
		"Bearer": map[string]interface{}{
			"type":        "apiKey",
			"name":        "Authorization",
			"in":          "header",
			"description": "Enter your bearer token in the format **Bearer <token>**",
		},
	}
	// Apply Bearer token security globally
	merged["security"] = []map[string][]string{
		{"Bearer": {}},
	}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || filepath.Ext(path) != ".json" {
			return nil
		}
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		var doc map[string]interface{}
		if err := json.Unmarshal(data, &doc); err != nil {
			return err
		}

		// Merge paths
		if paths, ok := doc["paths"].(map[string]interface{}); ok {
			for k, v := range paths {
				merged["paths"].(map[string]interface{})[k] = v
			}
		}
		// Merge definitions
		if defs, ok := doc["definitions"].(map[string]interface{}); ok {
			for k, v := range defs {
				merged["definitions"].(map[string]interface{})[k] = v
			}
		}
		// Merge tags
		if tags, ok := doc["tags"].([]interface{}); ok {
			merged["tags"] = append(merged["tags"].([]interface{}), tags...)
		}
		// Set info from the first file
		if _, ok := merged["info"]; !ok {
			merged["info"] = doc["info"]
		}
		return nil
	})
	if err != nil {
		return err
	}

	// Write merged swagger
	out, err := json.MarshalIndent(merged, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(output, out, 0o644)
}
