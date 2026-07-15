package utils

import (
	"net/http"
	"net/url"
	"strings"

	"crosshatch/internal/config"
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

	for _, allowed := range config.AllowedOrigins() {
		if strings.EqualFold(allowed, origin) {
			return true
		}
	}
	return false
}
