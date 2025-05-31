package middlewares

import (
	"context"
	"net/http"
	"regexp"
	"strings"

	"github.com/tuandq2112/organization-service/appconfig"
)

func isWhitelisted(path, method string) bool {
	if strings.HasPrefix(path, "/swagger/") {
		return true
	}

	var whitelist []string

	switch method {
	case "GET":
		whitelist = appconfig.WHITELIST_METHODS_GET_PATH
	case "POST":
		whitelist = appconfig.WHITELIST_METHODS_POST_PATH
	case "PUT":
		whitelist = appconfig.WHITELIST_METHODS_PUT_PATH
	case "DELETE":
		whitelist = appconfig.WHITELIST_METHODS_DELETE_PATH
	default:
		return false
	}

	for _, route := range whitelist {
		if route == path {
			return true
		}
		if strings.Contains(route, "/*") {
			pattern := "^" + strings.ReplaceAll(regexp.QuoteMeta(route), `\*`, `([^/]+)`) + `$`
			matched, _ := regexp.MatchString(pattern, path)
			if matched {
				return true
			}
		}
	}
	return false
}

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isWhitelisted(r.URL.Path, r.Method) {
			next.ServeHTTP(w, r)
			return
		}
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusForbidden)
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		tokenString := tokenParts[1]
		claims, err := ParseJWT(tokenString)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), appconfig.USER_CONTEXT, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
