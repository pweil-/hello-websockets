package main

import (
	"fmt"
	"io"
	"net/http"

	"golang.org/x/net/websocket"
)

func echoServer(ws *websocket.Conn) {
	io.Copy(ws, ws)
}

func echoServerSecure(ws *websocket.Conn) {
	io.Copy(ws, ws)
}

func main() {
	http.Handle("/echo", websocket.Handler(echoServer))
	err := http.ListenAndServe("9999", nil)

	if err != nil {
		panic(fmt.Sprintf("Error in ListenAndServe: %s", err.Error()))
	}
}
