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

	/* UNCOMMENT FROM PRODUCTION */
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
	jsonmessage := json.RawMessage(message.Body)

	// jsonmessage := json.RawMessage(`{"port":80,"limits":{"timeout":10,"requests":100},"processes":{"max":1,"spare":1,"idle_timeout":20},"user":"root","group":"root","environment":{"newKey":"New Value","newKey-1":"New Value"},"arguments":["--tmp","/tmp/fello"],"lang":"go","appname":"gobin","repo":" https://github.com/username/repo.git","cloud":"aws","working_directory":"/workingdir/","executable":"gobinapp"}`)

	var p unit.BaseStruct
	if err := json.Unmarshal(jsonmessage, &p); err != nil {
		panic(err)
	}

	//Build the base configuration for the machine
	/* UNCOMMENT FROM PRODUCTION */
	unit.GenericPrepBuild(&p, client)

	// Switch used to select the appropriate language modules to build
	switch lang := p.Lang; lang {
	case "go", "nodejs":
		var externalStruct unit.External
		if err := json.Unmarshal(jsonmessage, &externalStruct); err != nil {
			panic(err)
		}
		// Create the config.json file for the language
		/* UNCOMMENT FOR PRODUCTION */
		unit.EchoExternalJSONConfig(&externalStruct, client)
		// unit.EchoExternalJSONConfig(&externalStruct)

		//Provision to the cloud platform
		unit.BuildExternalMachine(&externalStruct, client)

	case "java":
		var javaStruct unit.Java
		if err := json.Unmarshal(jsonmessage, &javaStruct); err != nil {
			panic(err)
		}
		// Create the config.json file for the language
		/* UNCOMMENT FOR PRODUCTION */
		unit.EchoJavaJSONConfig(&javaStruct, client)
		// unit.EchoJavaJSONConfig(&javaStruct)

	case "perl":
		var perlStruct unit.Perl

		if err := json.Unmarshal(jsonmessage, &perlStruct); err != nil {
			panic(err)
		}
		// Create the config.json file for the language
		/* UNCOMMENT FOR PRODUCTION */
		unit.EchoPerlJSONConfig(&perlStruct, client)
		// unit.EchoPerlJSONConfig(&perlStruct)

	case "php":
		var phpStruct unit.PHP

		if err := json.Unmarshal(jsonmessage, &phpStruct); err != nil {
			panic(err)
		}
		// Create the config.json file for the language
		/* UNCOMMENT FOR PRODUCTION */
		unit.EchoPHPJSONConfig(&phpStruct, client)
		// unit.EchoPHPJSONConfig(&phpStruct)

	case "python":
		var pythonStruct unit.Python

		if err := json.Unmarshal(jsonmessage, &pythonStruct); err != nil {
			panic(err)
		}
		// Create the config.json file for the language
		/* UNCOMMENT FOR PRODUCTION */
		unit.EchoPythonJSONConfig(&pythonStruct, client)
		// unit.EchoPythonJSONConfig(&pythonStruct)

	case "ruby":
		var rubyStruct unit.Ruby

		if err := json.Unmarshal(jsonmessage, &rubyStruct); err != nil {
			panic(err)
		}
		// Create the config.json file for the language
		/* UNCOMMENT FOR PRODUCTION */
		unit.EchoRubyJSONConfig(&rubyStruct, client)
		// unit.EchoRubyJSONConfig(&rubyStruct)
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
