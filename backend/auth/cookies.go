package auth

import (
	"net/http"
	"time"
)

func GenerateAccessCookie(value string) *http.Cookie {
	accessCookie := &http.Cookie{
		Name:     "access-token",
		HttpOnly: true,
		Expires:  time.Now().Add(time.Minute * 15),
		Value:    value,
		Path:     "/",
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	}

	return accessCookie
}

func GenerateRefreshCookie(value string) *http.Cookie {
	refreshCookie := &http.Cookie{
		Name:     "refresh-token",
		HttpOnly: true,
		Expires:  time.Now().Add(time.Hour * 24),
		Value:    value,
		Path:     "/",
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	}

	return refreshCookie
}

func GetCookieValue(request *http.Request, cookieString string) (string, error) {
	cookie, err := request.Cookie(cookieString)

	if err != nil {
		return "", err
	}

	return cookie.Value, nil
}
