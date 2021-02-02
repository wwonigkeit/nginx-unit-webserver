package main

import (
	"fmt"
	"net/http"

	"github.com/wwonigkeit/nginx-unit-webserver/backend/websocket"
)

var message websocket.Message

func serveWs(pool *websocket.Pool, w http.ResponseWriter, r *http.Request) {
	fmt.Println("WebSocket Endpoint Hit")
	conn, err := websocket.Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+v\n", err)
	}

	client := &websocket.Client{
		Conn: conn,
		Pool: pool,
	}

	pool.Register <- client
	message, _ = client.Read()

	fmt.Printf("Message Received: %+v\n", message)

	//client.Pool.Broadcast <- message

	//switch lang := message.Body.lang; lang {
	//case condition:

	//}
}

func setupRoutes() {
	pool := websocket.NewPool()
	go pool.Start()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(pool, w, r)
	})
}

func main() {
	fmt.Println("NGINX Unit Deployment App v0.01")
	setupRoutes()
	http.ListenAndServe(":8080", nil)
}
