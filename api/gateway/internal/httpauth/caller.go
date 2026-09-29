package httpauth

import (
	"common/identity"
	"common/jwtverify"
	"context"
	"net/http"
	"strings"
)

type CallerInput struct {
	Principal      *jwtverify.Principal
	APIKey         *VerifiedAPIKey
	AllowAnonymous bool
}

type CallerResolution struct {
	Caller    identity.Caller
	APIKey    VerifiedAPIKey
	HasAPIKey bool
}

func ResolveHTTPCaller(
	ctx context.Context,
	r *http.Request,
	principal *jwtverify.Principal,
	allowAnonymous bool,
	verifier APIKeyVerifier,
) (CallerResolution, error) {
	raw, hasKey, err := APIKeyFromRequest(r)
	if err != nil {
		return CallerResolution{}, err
	}
	var verified *VerifiedAPIKey
	if hasKey {
		if verifier == nil {
			return CallerResolution{}, ErrAPIKeyUnavailable
		}
		value, err := verifier.VerifyAPIKey(ctx, raw)
		if err != nil {
			return CallerResolution{}, err
		}
		verified = &value
	}

	caller, err := ClassifyCaller(CallerInput{Principal: principal, APIKey: verified, AllowAnonymous: allowAnonymous})
	if err != nil {
		return CallerResolution{}, err
	}

	result := CallerResolution{
		Caller:    caller,
		HasAPIKey: hasKey,
	}

	if verified != nil {
		result.APIKey = VerifiedAPIKey{ID: verified.ID, AppName: strings.TrimSpace(verified.AppName)}
		if r != nil {
			r.Header.Del(LegacyAPIKeyHeader)
			r.Header.Set(APIKeyHeader, raw)
		}
	}

	return result, nil
}

func ClassifyCaller(input CallerInput) (identity.Caller, error) {
	if input.Principal != nil {
		switch jwtverify.TokenType(input.Principal.Type) {
		case jwtverify.AccessToken, jwtverify.TempToken:
		default:
			return identity.Caller{}, ErrTokenTypeNotAllowed
		}
		if !ActorFromPrincipal(input.Principal).IsValid() {
			return identity.Caller{}, ErrCallerClassification
		}
	}

	verified, hasKey, err := normalizeVerifiedAPIKey(input.APIKey)
	if err != nil {
		return identity.Caller{}, err
	}
	caller := identity.Caller{}

	switch {
	case input.Principal != nil:
		caller.Kind = identity.CallerUser
	case hasKey:
		caller.Kind = identity.CallerAPIKey
	case input.AllowAnonymous:
		caller.Kind = identity.CallerAnonymous
	default:
		return identity.Caller{}, ErrAuthenticationRequired
	}

	if hasKey {
		caller.APIKeyID = verified.ID
		caller.APIKeyApp = verified.AppName
	}
	if !caller.IsValid() {
		return identity.Caller{}, ErrCallerClassification
	}

	return caller, nil
}

func normalizeVerifiedAPIKey(key *VerifiedAPIKey) (VerifiedAPIKey, bool, error) {
	if key == nil {
		return VerifiedAPIKey{}, false, nil
	}
	app := strings.TrimSpace(key.AppName)
	if key.ID == 0 || app == "" || len(app) > 128 || containsControl(app) {
		return VerifiedAPIKey{}, false, ErrCallerClassification
	}
	return VerifiedAPIKey{ID: key.ID, AppName: app}, true, nil
}

func containsControl(value string) bool {
	for _, r := range value {
		if r < 0x20 || r == 0x7f {
			return true
		}
	}
	return false
}
