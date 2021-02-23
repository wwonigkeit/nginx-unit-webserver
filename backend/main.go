package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/wwonigkeit/nginx-unit-webserver/backend/cloud"
	"github.com/wwonigkeit/nginx-unit-webserver/backend/unit"
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
	jsonmessage := json.RawMessage(message.Body)

	// From the message pare the base JSON structure to get the language
	// we're preparing to deploy
	var p unit.BaseStruct
	if err := json.Unmarshal(jsonmessage, &p); err != nil {
		panic(err)
	}

	// Build the base configuration for the machine, which includes getting the
	// correct NGINX unit docker container
	unit.Build(&p, client)

	var imagename string
	// Switch used to select the appropriate language modules to build
	switch lang := p.Lang; lang {
	case "go", "nodejs":
		var externalStruct unit.External
		if err := json.Unmarshal(jsonmessage, &externalStruct); err != nil {
			panic(err)
		}
		// Create the config.json file for the language
		unit.ExternalConfig(&externalStruct, client)

		//Provision to the cloud platform
		imagename = unit.BuildExternalImage(&externalStruct, client)

	case "java":
		var javaStruct unit.Java
		if err := json.Unmarshal(jsonmessage, &javaStruct); err != nil {
			panic(err)
		}
		// Create the config.json file for the language
		unit.JavaConfig(&javaStruct, client)

		//Provision image to the cloud platform
		imagename = unit.BuildJavaImage(&javaStruct, client)

	case "perl":
		var perlStruct unit.Perl

		if err := json.Unmarshal(jsonmessage, &perlStruct); err != nil {
			panic(err)
		}
		// Create the config.json file for the language
		unit.PerlConfig(&perlStruct, client)

		//Provision image to the cloud platform
		imagename = unit.BuildPerlImage(&perlStruct, client)

	case "php":
		var phpStruct unit.PHP

		if err := json.Unmarshal(jsonmessage, &phpStruct); err != nil {
			panic(err)
		}
		// Create the config.json file for the language
		unit.PhpConfig(&phpStruct, client)

		//Provision image to the cloud platform
		imagename = unit.BuildPhpImage(&phpStruct, client)

	case "python":
		var pythonStruct unit.Python

		if err := json.Unmarshal(jsonmessage, &pythonStruct); err != nil {
			panic(err)
		}
		// Create the config.json file for the language
		unit.PythonConfig(&pythonStruct, client)

		//Provision image to the cloud platform
		imagename = unit.BuildPythonImage(&pythonStruct, client)

	case "ruby":
		var rubyStruct unit.Ruby

		if err := json.Unmarshal(jsonmessage, &rubyStruct); err != nil {
			panic(err)
		}
		// Create the config.json file for the language
		unit.RubyConfig(&rubyStruct, client)

		//Provision image to the cloud platform
		imagename = unit.BuildRubyImage(&rubyStruct, client)
	}

	// Switch used to select the appropriate platform to deploy to
	switch platform := p.Cloud.Platform; platform {
	case "gcp":
		// Create a machine in Google Cloud Platform
		fmt.Println("Provisioning to GCP")
		cloud.BuildGcpInstance(imagename, p.Cloud.MachineType, client, p.Port)
	case "aws":
		// Create a machine in Amazon Web Services
		fmt.Println("Provisioning to AWS")
		cloud.BuildAwsInstance(imagename, p.Cloud.MachineType, client, p.Port)
	case "azure":
		//do azure
	}

	client.Pool.Broadcast <- websocket.Message{Type: 1, Body: "Closing the client connection"}

	pool.Unregister <- client
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
