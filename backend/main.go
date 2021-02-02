package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/wwonigkeit/nginx-unit-webserver/backend/unit"
	"github.com/wwonigkeit/nginx-unit-webserver/backend/websocket"
)

var message websocket.Message

func serveWs(pool *websocket.Pool, w http.ResponseWriter, r *http.Request) {
	fmt.Println("WebSocket Endpoint Hit")

	/*
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
		baseJSONObj := json.RawMessage(message.Body)
	*/

	baseJSONObj := json.RawMessage(`{"port":80,"limits":{"timeout":10,"requests":100},"processes":{"max":1,"spare":1,"idle_timeout":20},"user":"root","group":"root","environment":{"newKey":"New Value","newKey-1":"New Value"},"arguments":["--tmp","/tmp/fello"],"lang":"go","appname":"gobin","repo":" https://github.com/username/repo.git","cloud":"aws","working_directory":"/workingdir/","executable":"gobinapp"}`)

	var p unit.BaseStruct

	if err := json.Unmarshal(baseJSONObj, &p); err != nil {
		panic(err)
	}

	/*
		mapEnvironment, _ := json.Marshal(p.Environment)
		mapBase := structs.Map(p)

		fmt.Printf("%+v\n", p)
		fmt.Println(mapEnvironment)

		client.Pool.Broadcast <- message
	*/

	switch lang := p.Lang; lang {
	case "go", "nodejs":
		fmt.Println("Executing the Golang/NodeJS build")

		//client.Pool.Broadcast <- websocket.Message{Type: 2, Body: "Executing the Golang/NodeJS build"}

		var externalStruct unit.External

		if err := json.Unmarshal(baseJSONObj, &externalStruct); err != nil {
			panic(err)
		}

		unit.ExternalBuild(&externalStruct)

	case "java":
		fmt.Println("Executing the Java build")
		//client.Pool <- "Executing the Java build"
	case "perl":
		fmt.Println("Executing the Perl build")
		//client.Pool <- "Executing the Perl build"
	case "php":
		fmt.Println("Executing the PHP build")
		//client.Pool <- "Executing the PHP build"
	case "python":
		fmt.Println("Executing the Python build")
		//client.Pool <- "Executing the Python build"
	case "ruby":
		fmt.Println("Executing the Ruby build")
		//client.Pool <- "Executing the Ruby build"
	}
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
