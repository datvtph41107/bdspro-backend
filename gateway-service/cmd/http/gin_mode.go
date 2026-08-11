package cmd

import "github.com/gin-gonic/gin"

func desiredGinMode(
	configuredMode string,
) string {
	if configuredMode == "" {
		return gin.ReleaseMode
	}

	return configuredMode
}

func configureGinMode(
	configuredMode string,
) {
	gin.SetMode(
		desiredGinMode(
			configuredMode,
		),
	)
}
