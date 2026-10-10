package logging

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/mattn/go-isatty"
)

const (
	ansiGreen   = "\033[97;42m"
	ansiWhite   = "\033[90;47m"
	ansiYellow  = "\033[90;43m"
	ansiRed     = "\033[97;41m"
	ansiBlue    = "\033[97;44m"
	ansiMagenta = "\033[97;45m"
	ansiCyan    = "\033[97;46m"
	ansiReset   = "\033[0m"
)

func devConsoleColorEnabled() bool {
	if !DevConsoleEnabled() {
		return false
	}
	return isatty.IsTerminal(os.Stdout.Fd())
}

func colorLevel(level string) string {
	if !devConsoleColorEnabled() {
		return level
	}
	switch level {
	case "ERROR":
		return ansiRed + level + ansiReset
	case "WARN":
		return ansiYellow + level + ansiReset
	case "DEBUG":
		return ansiCyan + level + ansiReset
	default:
		return ansiGreen + level + ansiReset
	}
}

func statusCodeColor(code int) string {
	switch {
	case code >= http.StatusContinue && code < http.StatusOK:
		return ansiWhite
	case code >= http.StatusOK && code < http.StatusMultipleChoices:
		return ansiGreen
	case code >= http.StatusMultipleChoices && code < http.StatusBadRequest:
		return ansiWhite
	case code >= http.StatusBadRequest && code < http.StatusInternalServerError:
		return ansiYellow
	default:
		return ansiRed
	}
}

func methodColor(method string) string {
	switch method {
	case http.MethodGet:
		return ansiBlue
	case http.MethodPost:
		return ansiCyan
	case http.MethodPut:
		return ansiYellow
	case http.MethodDelete:
		return ansiRed
	case http.MethodPatch:
		return ansiGreen
	case http.MethodHead:
		return ansiMagenta
	case http.MethodOptions:
		return ansiWhite
	default:
		return ansiReset
	}
}

func latencyColor(latency time.Duration) string {
	switch {
	case latency < 100*time.Millisecond:
		return ansiWhite
	case latency < 200*time.Millisecond:
		return ansiGreen
	case latency < 300*time.Millisecond:
		return ansiCyan
	case latency < 500*time.Millisecond:
		return ansiBlue
	case latency < time.Second:
		return ansiYellow
	case latency < 2*time.Second:
		return ansiMagenta
	default:
		return ansiRed
	}
}

// truncateLatency mirrors gin defaultLogFormatter latency rounding.
func truncateLatency(latency time.Duration) time.Duration {
	switch {
	case latency > time.Minute:
		return latency.Truncate(10 * time.Second)
	case latency > time.Second:
		return latency.Truncate(10 * time.Millisecond)
	case latency > time.Millisecond:
		return latency.Truncate(10 * time.Microsecond)
	default:
		return latency
	}
}

// formatLatencyDisplay always includes a millisecond readout for dev consoles.
func formatLatencyDisplay(latency time.Duration) string {
	latency = truncateLatency(latency)
	ms := float64(latency) / float64(time.Millisecond)
	if latency < 10*time.Millisecond {
		return fmt.Sprintf("%9.3fms", ms)
	}
	if latency < time.Second {
		return fmt.Sprintf("%8.2fms", ms)
	}
	return fmt.Sprintf("%9.0fms", ms)
}
