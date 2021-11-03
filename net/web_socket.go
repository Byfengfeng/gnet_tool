package net

import (
	"golang.org/x/net/websocket"
	"net/http"
)

func Start()  {
	handler := websocket.Handler(func(conn *websocket.Conn) {

	})
	http.ListenAndServe("0.0.0.0:7777",handler)
}

func helloHandler(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("Hello World!"))
}