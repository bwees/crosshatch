// Package proxy provides a WebSocket reverse proxy to a go2rtc instance so the
// browser can reach live camera streams through the API origin.
package proxy

import (
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"time"

	"crosshatch/internal/utils"

	"github.com/gorilla/websocket"
	"go.uber.org/fx"
)

const (
	// defaultTarget is the go2rtc base URL when GO2RTC_WS_URL is unset.
	defaultTarget = "ws://localhost:1984"

	// upstreamPath is go2rtc's WebSocket API endpoint. The incoming query
	// (e.g. "?src=<serial>") is forwarded unchanged.
	upstreamPath = "/api/ws"

	closeGrace = time.Second
)

var Module = fx.Provide(NewGo2RTCProxy)

// Go2RTCProxy forwards an inbound WebSocket connection to go2rtc, piping frames
// in both directions while preserving their text/binary type.
type Go2RTCProxy struct {
	target   *url.URL
	upgrader websocket.Upgrader
	dialer   *websocket.Dialer
}

// NewGo2RTCProxy reads the upstream URL from GO2RTC_WS_URL, falling back to
// defaultTarget.
func NewGo2RTCProxy() *Go2RTCProxy {
	raw := os.Getenv("GO2RTC_WS_URL")
	if raw == "" {
		raw = defaultTarget
	}

	target, err := url.Parse(raw)
	if err != nil {
		slog.Error("go2rtc: invalid GO2RTC_WS_URL, using default", "value", raw, "error", err)
		target, _ = url.Parse(defaultTarget)
	}

	return &Go2RTCProxy{
		target: target,
		upgrader: websocket.Upgrader{
			CheckOrigin: utils.AllowedOrigin,
		},
		dialer: websocket.DefaultDialer,
	}
}

// ServeHTTP dials go2rtc first (so dial failures surface as a clean HTTP error)
// and then upgrades the client connection and relays between the two.
func (p *Go2RTCProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	upstreamURL := *p.target
	upstreamURL.Path = upstreamPath
	upstreamURL.RawQuery = r.URL.RawQuery

	// Dial without forwarding the Origin header; go2rtc rejects cross-origin
	// upgrades, and from its perspective this request originates server-side.
	upstream, resp, err := p.dialer.Dial(upstreamURL.String(), nil)
	if err != nil {
		status := http.StatusBadGateway
		if resp != nil {
			status = resp.StatusCode
		}
		slog.Error("go2rtc: upstream dial failed", "url", upstreamURL.String(), "status", status, "error", err)
		http.Error(w, "go2rtc upstream unavailable", status)
		return
	}
	defer upstream.Close()

	client, err := p.upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("go2rtc: client upgrade failed", "error", err)
		return
	}
	defer client.Close()

	// One goroutine per direction. Each connection has exactly one writer
	// goroutine, satisfying gorilla's single-concurrent-writer requirement.
	errc := make(chan error, 2)
	go relay(upstream, client, errc) // client -> upstream
	go relay(client, upstream, errc) // upstream -> client
	<-errc
}

// relay copies messages from src to dst until either side errors, then sends a
// best-effort close frame so the peer learns the stream ended.
func relay(dst, src *websocket.Conn, errc chan<- error) {
	for {
		mt, data, err := src.ReadMessage()
		if err != nil {
			_ = dst.WriteControl(websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
				time.Now().Add(closeGrace))
			errc <- err
			return
		}
		if err := dst.WriteMessage(mt, data); err != nil {
			errc <- err
			return
		}
	}
}
