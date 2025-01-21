package utils

import (
	"errors"
	"net/http"
)

func GetSessionCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			// No cookie, return no session
			return "", nil
		}
		return "", err
	}
	return cookie.Value, nil
}
