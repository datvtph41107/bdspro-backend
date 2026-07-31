package cookie

import (
	appcontext "common/context"
	"context"
	"net/http"
	"user/internal/interface/providers"
)

type BroswerCookieProvider struct {
	PathRefresh string
}

func NewCookieProvider() providers.CookieProvider {
	return &BroswerCookieProvider{
		PathRefresh: "/v2/auth/token/refresh",
	}
}

func (s *BroswerCookieProvider) ResetAll(c context.Context) {
	writer, ok := appcontext.GetResponseWriter(c)
	if !ok {
		return
	}

	request, ok := appcontext.GetRequest(c)
	if !ok {
		return
	}

	cookies := request.Cookies()
	for _, cookie := range cookies {
		http.SetCookie(writer, &http.Cookie{
			Name:     cookie.Name,
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
			Secure:   false,
		})
	}
}

func (s *BroswerCookieProvider) ResetCookie(c context.Context, cookie *http.Cookie) {
	writer, ok := appcontext.GetResponseWriter(c)
	if !ok {
		return
	}
	http.SetCookie(writer, &http.Cookie{
		Name:     cookie.Name,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   false,
	})
}

func (s *BroswerCookieProvider) SetRefreshToken(c context.Context, refreshToken string) {
	writer, ok := appcontext.GetResponseWriter(c)
	if !ok {
		return
	}
	http.SetCookie(writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		MaxAge:   60 * 60 * 24 * 30,
		HttpOnly: true,
		Secure:   false,
	})
}

func (s *BroswerCookieProvider) DeleteRefreshToken(c context.Context) {
	writer, ok := appcontext.GetResponseWriter(c)
	if !ok {
		return
	}
	http.SetCookie(writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   false,
	})
}

func (s *BroswerCookieProvider) Set(c context.Context, key string, value string) {
	writer, ok := appcontext.GetResponseWriter(c)
	if !ok {
		return
	}
	http.SetCookie(writer, &http.Cookie{
		Name:     key,
		Value:    value,
		Path:     "/",
		MaxAge:   60 * 60 * 24 * 30,
		HttpOnly: true,
		Secure:   false,
	})
}

func (s *BroswerCookieProvider) Expire(c context.Context, key string, seconds int) error {
	writer, ok := appcontext.GetResponseWriter(c)
	if !ok {
		return nil
	}
	http.SetCookie(writer, &http.Cookie{
		Name:     key,
		Value:    "",
		Path:     "/",
		MaxAge:   seconds,
		HttpOnly: true,
		Secure:   false,
	})
	return nil
}

func (s *BroswerCookieProvider) Get(c context.Context, key string) (string, error) {
	// request, ok := appcontext.GetRequest(c)
	// if !ok {
	// 	return "", nil
	// }
	// return request.Cookie(key).Value, nil
	return "", nil
}
