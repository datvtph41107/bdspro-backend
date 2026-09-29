package middlewares

import (
	"regexp"
	"strings"
)

var whiteListedRoutes = map[string]map[string]bool{
	"/launchpad/api/v1/user/login":             {"POST": true},
	"/launchpad/api/v1/user/login-with-google": {"POST": true},
	"/launchpad/api/v1/user/register":          {"POST": true},
	"/launchpad/api/v1/webhook-callback":       {"POST": true},
	"/launchpad/api/v1/project/public":         {"GET": true},
	"/launchpad/api/v1/project/public/*":       {"GET": true},
}

func isWhitelisted(path, method string) bool {
	// Kiểm tra khớp chính xác trước
	if methods, exists := whiteListedRoutes[path]; exists {
		return methods[method]
	}

	// Kiểm tra với wildcard (*)
	for route, methods := range whiteListedRoutes {
		if strings.Contains(route, "/*") {
			pattern := "^" + strings.ReplaceAll(regexp.QuoteMeta(route), `\*`, `([^/]+)`) + `$`
			matched, _ := regexp.MatchString(pattern, path)
			if matched {
				return methods[method]
			}
		}
	}
	return false
}