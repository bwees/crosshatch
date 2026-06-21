package controllers

import (
	"crosshatch/internal/proxy"
	"crosshatch/internal/socketio"

	"github.com/go-fuego/fuego"
)

type WebsocketController struct {
	socket *socketio.SocketIO
	proxy  *proxy.Go2RTCProxy
}

func (c *WebsocketController) Register(api *fuego.Server) {
	fuego.Handle(api, "/ws/", c.socket.Handler())
	fuego.GetStd(api, "/go2rtc", c.proxy.ServeHTTP, fuego.OptionHide())
}

func NewWebsocketController(socket *socketio.SocketIO, p *proxy.Go2RTCProxy) *WebsocketController {
	return &WebsocketController{socket: socket, proxy: p}
}
