// Package config is the single source of truth for environment-derived
// configuration. Every setting is read here so defaults live in one place.
package config

import (
	"os"
	"strings"
)

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// DataDir is the directory where persistent data (the SQLite database) is
// stored. Defaults to the current directory for local development; production
// sets DATA_DIR to a mounted volume.
func DataDir() string {
	return env("DATA_DIR", ".")
}

// Port is the TCP port the HTTP server listens on.
func Port() string {
	return env("PORT", "3000")
}

// WebStaticPath is the directory of the built web frontend to serve. When
// empty (the dev default) the server does not serve the frontend.
func WebStaticPath() string {
	return os.Getenv("WEB_STATIC_PATH")
}

// Go2RTCWSURL is the go2rtc WebSocket base URL the camera proxy relays to.
func Go2RTCWSURL() string {
	return env("GO2RTC_WS_URL", "ws://localhost:1984")
}

// Go2RTCAPIURL is the go2rtc HTTP API base URL used to manage streams.
func Go2RTCAPIURL() string {
	return env("GO2RTC_API_URL", "http://localhost:1984")
}

// VapidSubject is the "sub" claim (a mailto: or URL) sent with web push.
func VapidSubject() string {
	return env("VAPID_SUBJECT", "crosshatch@bwees.io")
}

// VapidPublicKey and VapidPrivateKey are the web-push VAPID keys. When either
// is empty they are loaded from, or generated and persisted to, the database.
func VapidPublicKey() string  { return os.Getenv("VAPID_PUBLIC_KEY") }
func VapidPrivateKey() string { return os.Getenv("VAPID_PRIVATE_KEY") }

// AllowedOrigins is the list of extra origins permitted for WebSocket upgrades,
// beyond same-origin requests.
func AllowedOrigins() []string {
	raw := os.Getenv("WS_ALLOWED_ORIGINS")
	if raw == "" {
		return nil
	}

	var origins []string
	for _, o := range strings.Split(raw, ",") {
		if o = strings.TrimSpace(o); o != "" {
			origins = append(origins, o)
		}
	}
	return origins
}
