package net

import (
	"fmt"
	"golang.org/x/net/websocket"
	"net/http"
)

func Start()  {
	handler := websocket.Handler(func(conn *websocket.Conn) {
		bytes := make([]byte,1024)
		fmt.Println(conn.Request().Header)
		conn.Read(bytes)
		fmt.Println(string(bytes))
		conn.Write([]byte("Hello World!"))
	})
	http.ListenAndServe("0.0.0.0:7777",handler)
}

func helloHandler(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("Hello World!"))
}