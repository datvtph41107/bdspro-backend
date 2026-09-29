package middlewares

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"chat/internal/constants"
	commonhttp "common/httpresponse"
	commonjwt "common/jwt"
)

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isWhitelisted(r.URL.Path, r.Method) {
			next.ServeHTTP(w, r)
			return
		}

		var tokenString string

		if strings.HasPrefix(r.URL.Path, "/ws") {
			tokenString = r.URL.Query().Get("token")
			if tokenString == "" {
				tokenString = r.Header.Get("Authorization")
			}
		} else {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				writeErrorResponse(r.Context(), w, commonjwt.AuthorizationHeaderMissingFault())
				return
			}

			tokenString = authHeader
			// tokenParts := strings.Split(authHeader, " ")
			// if len(tokenParts) == 2 && strings.ToLower(tokenParts[0]) == "bearer" {
			// 	tokenString = tokenParts[1]
			// } else {
			// 	writeErrorResponse(w, http.StatusUnauthorized, "Invalid Authorization header format")
			// 	return
			// }
		}

		if tokenString == "" {
			writeErrorResponse(r.Context(), w, commonjwt.TokenMissingFault())
			return
		}

		claims, err := ExtractPayload(tokenString)
		if err != nil {
			writeErrorResponse(r.Context(), w, commonjwt.InvalidOrExpiredTokenFault(err))
			return
		}

		ctx := context.WithValue(r.Context(), constants.CONTEXT_USER_KEY, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ExtractPayload(tokenString string) (map[string]any, error) {
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid JWT format")
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, errors.New("failed to decode payload")
	}

	var payload map[string]any
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		return nil, errors.New("failed to parse payload")
	}

	return payload, nil
}

func writeErrorResponse(
	ctx context.Context,
	w http.ResponseWriter,
	err error,
) {
	commonhttp.WriteProblem(
		ctx,
		w,
		commonhttp.ProblemFromError(err),
		commonhttp.WithLegacyDataJSONEnvelope(),
	)
}
