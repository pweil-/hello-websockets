package main

import (
	"golang.org/x/net/websocket"
	"fmt"
	"log"
	"crypto/tls"
)

func main() {
	wsConfig, _ := websocket.NewConfig("wss://www.example2.com/echo", "http://localhost/")
	wsConfig.TlsConfig = &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         "www.example2.com",
	}


	ws, err := websocket.DialConfig(wsConfig)
	if err != nil {
		log.Fatal(err)
	}

	message := []byte("hello, world!")
	_, err = ws.Write(message)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Send: %s\n", message)

	var msg = make([]byte, len(message))
	_, err = ws.Read(msg)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Receive: %s\n", msg)
}
