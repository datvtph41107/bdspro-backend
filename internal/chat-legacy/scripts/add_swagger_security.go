package scripts

import (
	"common/logging"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"chat/config"
)

func AddSwaggerSecurity() {
	logger := logging.WithComponent(
		context.Background(),
		"scripts",
	)

	config.LoadConfig()

	logger.Info(
		fmt.Sprintf(
			"Swagger path: %s",
			config.AppProperties.Swagger.Path,
		),
	)

	data, err := ioutil.ReadFile(
		config.AppProperties.Swagger.Path,
	)
	if err != nil {
		logger.Error(
			fmt.Sprintf(
				"Error reading file: %v",
				err,
			),
		)
		return
	}

	var swagger map[string]interface{}

	if err := json.Unmarshal(
		data,
		&swagger,
	); err != nil {
		logger.Error(
			fmt.Sprintf(
				"Error parsing JSON: %v",
				err,
			),
		)
		return
	}

	swagger["securityDefinitions"] = map[string]interface{}{
		"BearerAuth": map[string]interface{}{
			"type":        "apiKey",
			"name":        "Authorization",
			"in":          "header",
			"description": "Enter your token in the format: Bearer <token>",
		},
	}

	swagger["host"] = "localhost:8081"

	swagger["security"] = []map[string]interface{}{
		{"BearerAuth": []interface{}{}},
	}

	updatedData, err := json.MarshalIndent(
		swagger,
		"",
		"  ",
	)
	if err != nil {
		logger.Error(
			fmt.Sprintf(
				"Error encoding JSON: %v",
				err,
			),
		)
		return
	}

	if err := ioutil.WriteFile(
		config.AppProperties.Swagger.Path,
		updatedData,
		0644,
	); err != nil {
		logger.Error(
			fmt.Sprintf(
				"Error writing file: %v",
				err,
			),
		)
		return
	}

	logger.Info(
		"Updated swagger.json successfully!",
	)
}
