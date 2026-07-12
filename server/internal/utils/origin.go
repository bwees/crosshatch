package utils

import (
	"net/http"
	"net/url"
	"os"
	"strings"
)

func AllowedOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")
	if origin == "" {
		return true
	}

	u, err := url.Parse(origin)
	if err != nil {
		return false
	}
	if strings.EqualFold(u.Host, r.Host) {
		return true
	}

	for _, allowed := range strings.Split(os.Getenv("WS_ALLOWED_ORIGINS"), ",") {
		if allowed = strings.TrimSpace(allowed); allowed != "" && strings.EqualFold(allowed, origin) {
			return true
		}
	}
	return false
}
