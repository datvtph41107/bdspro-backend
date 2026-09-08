package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Swagger 2.0 structure
type Swagger struct {
	Swagger             string                 `json:"swagger"`
	Info                Info                   `json:"info"`
	Host                string                 `json:"host,omitempty"`
	BasePath            string                 `json:"basePath,omitempty"`
	Schemes             []string               `json:"schemes,omitempty"`
	Consumes            []string               `json:"consumes,omitempty"`
	Produces            []string               `json:"produces,omitempty"`
	Paths               map[string]PathItem    `json:"paths"`
	Definitions         map[string]interface{} `json:"definitions"`
	Parameters          map[string]interface{} `json:"parameters,omitempty"`
	Responses           map[string]interface{} `json:"responses,omitempty"`
	SecurityDefinitions map[string]Security    `json:"securityDefinitions,omitempty"`
	Security            []map[string][]string  `json:"security,omitempty"`
	Tags                []Tag                  `json:"tags,omitempty"`
	ExternalDocs        interface{}            `json:"externalDocs,omitempty"`
}

type Info struct {
	Title          string  `json:"title"`
	Description    string  `json:"description,omitempty"`
	Version        string  `json:"version"`
	TermsOfService string  `json:"termsOfService,omitempty"`
	Contact        Contact `json:"contact,omitempty"`
	License        License `json:"license,omitempty"`
}

type Contact struct {
	Name  string `json:"name,omitempty"`
	URL   string `json:"url,omitempty"`
	Email string `json:"email,omitempty"`
}

type License struct {
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}

type PathItem struct {
	Ref     string     `json:"$ref,omitempty"`
	Get     *Operation `json:"get,omitempty"`
	Put     *Operation `json:"put,omitempty"`
	Post    *Operation `json:"post,omitempty"`
	Delete  *Operation `json:"delete,omitempty"`
	Options *Operation `json:"options,omitempty"`
	Head    *Operation `json:"head,omitempty"`
	Patch   *Operation `json:"patch,omitempty"`
}

type Operation struct {
	Tags        []string               `json:"tags,omitempty"`
	Summary     string                 `json:"summary,omitempty"`
	Description string                 `json:"description,omitempty"`
	OperationID string                 `json:"operationId,omitempty"`
	Consumes    []string               `json:"consumes,omitempty"`
	Produces    []string               `json:"produces,omitempty"`
	Parameters  []Parameter            `json:"parameters,omitempty"`
	Responses   map[string]Response    `json:"responses"`
	Schemes     []string               `json:"schemes,omitempty"`
	Deprecated  bool                   `json:"deprecated,omitempty"`
	Security    []map[string][]string  `json:"security,omitempty"`
	Extensions  map[string]interface{} `json:"-"`
}

type Parameter struct {
	Name        string      `json:"name"`
	In          string      `json:"in"`
	Description string      `json:"description,omitempty"`
	Required    bool        `json:"required,omitempty"`
	Type        string      `json:"type,omitempty"`
	Format      string      `json:"format,omitempty"`
	Schema      interface{} `json:"schema,omitempty"`
	Items       interface{} `json:"items,omitempty"`
	Enum        []string    `json:"enum,omitempty"`
	Default     interface{} `json:"default,omitempty"`
}

type Response struct {
	Description string                 `json:"description"`
	Schema      interface{}            `json:"schema,omitempty"`
	Headers     map[string]interface{} `json:"headers,omitempty"`
	Examples    map[string]interface{} `json:"examples,omitempty"`
}

type Security struct {
	Type             string            `json:"type"`
	Description      string            `json:"description,omitempty"`
	Name             string            `json:"name,omitempty"`
	In               string            `json:"in,omitempty"`
	Flow             string            `json:"flow,omitempty"`
	AuthorizationURL string            `json:"authorizationUrl,omitempty"`
	TokenURL         string            `json:"tokenUrl,omitempty"`
	Scopes           map[string]string `json:"scopes,omitempty"`
}

type Tag struct {
	Name         string      `json:"name"`
	Description  string      `json:"description,omitempty"`
	ExternalDocs interface{} `json:"externalDocs,omitempty"`
}

func main() {
	docsDir := "docs"
	outputFile := "merged_swagger.json"

	// Initialize merged swagger
	mergedSwagger := &Swagger{
		Swagger: "2.0",
		Info: Info{
			Title:       "BDS Pro API Documentation",
			Description: "Complete API documentation for BDS Pro microservices",
			Version:     "1.0.0",
		},
		BasePath:            "",
		Consumes:            []string{"application/json"},
		Produces:            []string{"application/json"},
		Paths:               make(map[string]PathItem),
		Definitions:         make(map[string]interface{}),
		SecurityDefinitions: make(map[string]Security),
		Security: []map[string][]string{
			{"BearerAuth": []string{}},
		},
	}

	// Add Bearer authentication
	mergedSwagger.SecurityDefinitions["BearerAuth"] = Security{
		Type:        "apiKey",
		Name:        "Authorization",
		In:          "header",
		Description: "Bearer token for authentication",
	}

	// Walk through all subdirectories in docs
	err := filepath.WalkDir(docsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip if not a file or not a JSON file
		if d.IsDir() || !strings.HasSuffix(d.Name(), ".swagger.json") {
			return nil
		}

		fmt.Printf("Processing: %s\n", path)

		// Read and parse JSON file
		data, err := os.ReadFile(path)
		if err != nil {
			log.Printf("Error reading file %s: %v", path, err)
			return nil
		}

		var swagger Swagger
		if err := json.Unmarshal(data, &swagger); err != nil {
			log.Printf("Error parsing JSON in %s: %v", path, err)
			return nil
		}

		// Merge paths
		for path, pathItem := range swagger.Paths {
			// Add Bearer authentication to all operations
			if pathItem.Get != nil {
				addBearerAuthToOperation(pathItem.Get)
			}
			if pathItem.Post != nil {
				addBearerAuthToOperation(pathItem.Post)
			}
			if pathItem.Put != nil {
				addBearerAuthToOperation(pathItem.Put)
			}
			if pathItem.Delete != nil {
				addBearerAuthToOperation(pathItem.Delete)
			}
			if pathItem.Patch != nil {
				addBearerAuthToOperation(pathItem.Patch)
			}

			// Use the path from the file (don't modify it)
			mergedSwagger.Paths[path] = pathItem
		}

		// Merge definitions
		for name, definition := range swagger.Definitions {
			mergedSwagger.Definitions[name] = definition
		}

		// Merge tags
		for _, tag := range swagger.Tags {
			// Check if tag already exists
			exists := false
			for _, existingTag := range mergedSwagger.Tags {
				if existingTag.Name == tag.Name {
					exists = true
					break
				}
			}
			if !exists {
				mergedSwagger.Tags = append(mergedSwagger.Tags, tag)
			}
		}

		return nil
	})
	if err != nil {
		log.Fatalf("Error walking directory: %v", err)
	}

	// Write merged swagger to file
	outputData, err := json.MarshalIndent(mergedSwagger, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling merged swagger: %v", err)
	}

	if err := os.WriteFile(outputFile, outputData, 0o644); err != nil {
		log.Fatalf("Error writing output file: %v", err)
	}

	// Copy merged swagger to gateway service docs directory
	gatewayDocsDir := "../../gateway-service/docs"
	if err := os.MkdirAll(gatewayDocsDir, 0o755); err != nil {
		log.Printf("Warning: Could not create gateway docs directory: %v", err)
	} else {
		gatewayOutputFile := filepath.Join(gatewayDocsDir, "merged_swagger.json")
		if err := os.WriteFile(gatewayOutputFile, outputData, 0o644); err != nil {
			log.Printf("Warning: Could not copy to gateway docs: %v", err)
		} else {
			fmt.Printf("Copied merged swagger to %s\n", gatewayOutputFile)
		}
	}

	fmt.Printf("Successfully merged all swagger files into %s\n", outputFile)
	fmt.Printf("Total paths: %d\n", len(mergedSwagger.Paths))
	fmt.Printf("Total definitions: %d\n", len(mergedSwagger.Definitions))
	fmt.Printf("Total tags: %d\n", len(mergedSwagger.Tags))
}

func addBearerAuthToOperation(operation *Operation) {
	// Add Bearer authentication to the operation
	operation.Security = []map[string][]string{
		{"BearerAuth": []string{}},
	}

	// Add Authorization parameter if not already present
	// hasAuthParam := false
	// for _, param := range operation.Parameters {
	// 	if param.Name == "Authorization" && param.In == "header" {
	// 		hasAuthParam = true
	// 		break
	// 	}
	// }

	// if !hasAuthParam {
	// 	authParam := Parameter{
	// 		Name:        "Authorization",
	// 		In:          "header",
	// 		Description: "Bearer token for authentication",
	// 		Required:    true,
	// 		Type:        "string",
	// 	}
	// 	operation.Parameters = append(operation.Parameters, authParam)
	// }
}
