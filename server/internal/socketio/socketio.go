package socketio

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"crosshatch/internal/utils"

	"github.com/zishang520/engine.io/v2/types"
	"github.com/zishang520/socket.io/v2/socket"
	"go.uber.org/fx"
)

const path = "/api/ws"

type EmitFunc func(event string, payload any)

type SocketIO struct {
	io        *socket.Server
	onConnect []func(emit EmitFunc)
}

func NewSocketIO(lc fx.Lifecycle) *SocketIO {
	opts := socket.DefaultServerOptions()
	opts.SetPath(path)
	opts.SetAllowRequest(func(ctx *types.HttpContext) error {
		if !utils.AllowedOrigin(ctx.Request()) {
			return errors.New("origin not allowed")
		}
		return nil
	})

	s := &SocketIO{io: socket.NewServer(nil, opts)}

	s.io.On("connection", func(clients ...any) {
		client := clients[0].(*socket.Socket)
		emit := func(event string, payload any) {
			emitTo(client, event, payload)
		}
		for _, fn := range s.onConnect {
			fn(emit)
		}
	})

	lc.Append(fx.Hook{
		OnStop: func(context.Context) error {
			s.io.Close(nil)
			return nil
		},
	})

	return s
}

func (s *SocketIO) OnConnect(fn func(emit EmitFunc)) {
	s.onConnect = append(s.onConnect, fn)
}

func (s *SocketIO) Emit(event string, payload any) {
	data, ok := encode(event, payload)
	if !ok {
		return
	}
	s.io.Emit(event, data)
}

func (s *SocketIO) Handler() http.Handler {
	return s.io.ServeHandler(nil)
}

func emitTo(client *socket.Socket, event string, payload any) {
	data, ok := encode(event, payload)
	if !ok {
		return
	}
	if err := client.Emit(event, data); err != nil {
		slog.Error("socketio: client emit failed", "event", event, "error", err)
	}
}

func encode(event string, payload any) (string, bool) {
	data, err := json.Marshal(payload)
	if err != nil {
		slog.Error("socketio: failed to encode event", "event", event, "error", err)
		return "", false
	}
	return string(data), true
}
