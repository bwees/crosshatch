package socketio

import "go.uber.org/fx"

var Module = fx.Provide(NewSocketIO)
