package net

import (
	"golang.org/x/net/websocket"
	"net/http"
)

type webSocketListen struct {
	address string
	channelHandel func(conn *websocket.Conn)
}

func NewWebSocket(addr string,channelHandel func(conn *websocket.Conn)) *webSocketListen {
	return &webSocketListen{address: addr,channelHandel: channelHandel}
}

func (w *webSocketListen) Start() error {
	handler := websocket.Handler(func(conn *websocket.Conn) {
		w.channelHandel(conn)
	})
	http.Handle("/add",handler)
	return http.ListenAndServe(w.address, nil)
}