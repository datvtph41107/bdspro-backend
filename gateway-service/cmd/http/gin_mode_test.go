package cmd

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestDesiredGinModeDefaultsToRelease(
	t *testing.T,
) {
	if got :=
		desiredGinMode(
			"",
		); got != gin.ReleaseMode {

		t.Fatalf(
			"mode = %q, want %q",
			got,
			gin.ReleaseMode,
		)
	}
}

func TestDesiredGinModePreservesExplicitMode(
	t *testing.T,
) {
	testCases :=
		[]string{
			gin.DebugMode,
			gin.ReleaseMode,
			gin.TestMode,
		}

	for _, mode := range testCases {

		t.Run(
			mode,
			func(
				t *testing.T,
			) {
				if got :=
					desiredGinMode(
						mode,
					); got != mode {

					t.Fatalf(
						"mode = %q, want %q",
						got,
						mode,
					)
				}
			},
		)
	}
}

func TestConfigureGinModeDefaultsToQuietReleaseMode(
	t *testing.T,
) {
	previousMode :=
		gin.Mode()

	previousWriter :=
		gin.DefaultWriter

	defer func() {
		gin.DefaultWriter =
			previousWriter

		gin.SetMode(
			previousMode,
		)
	}()

	var output bytes.Buffer

	gin.DefaultWriter =
		&output

	configureGinMode(
		"",
	)

	router :=
		gin.New()

	router.GET(
		"/health",
		func(
			ctx *gin.Context,
		) {
			ctx.Status(
				http.StatusNoContent,
			)
		},
	)

	if got :=
		gin.Mode(); got != gin.ReleaseMode {

		t.Fatalf(
			"mode = %q, want %q",
			got,
			gin.ReleaseMode,
		)
	}

	if output.Len() != 0 {
		t.Fatalf(
			"release mode emitted Gin debug output:\n%s",
			output.String(),
		)
	}
}

func TestConfigureGinModeHonorsExplicitDebug(
	t *testing.T,
) {
	previousMode :=
		gin.Mode()

	previousWriter :=
		gin.DefaultWriter

	defer func() {
		gin.DefaultWriter =
			previousWriter

		gin.SetMode(
			previousMode,
		)
	}()

	var output bytes.Buffer

	gin.DefaultWriter =
		&output

	configureGinMode(
		gin.DebugMode,
	)

	router :=
		gin.New()

	router.GET(
		"/health",
		func(
			ctx *gin.Context,
		) {
			ctx.Status(
				http.StatusNoContent,
			)
		},
	)

	if got :=
		gin.Mode(); got != gin.DebugMode {

		t.Fatalf(
			"mode = %q, want %q",
			got,
			gin.DebugMode,
		)
	}

	if output.Len() == 0 {
		t.Fatal(
			"explicit debug mode emitted no Gin diagnostics",
		)
	}
}
