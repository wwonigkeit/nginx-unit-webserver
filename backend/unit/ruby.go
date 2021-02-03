package unit

import (
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/wwonigkeit/nginx-unit-webserver/backend/websocket"
)

//EchoRubyJSONConfig echos the configuration for the Ruby
//Unit startup configurations
/* UNCOMMENT FOR PRODUCTION */
func EchoRubyJSONConfig(rubyJSONObj *Ruby, c *websocket.Client) {
	// func EchoRubyJSONConfig(rubyJSONObj *Ruby) {
	var environmentstr string

	for key, value := range rubyJSONObj.Environment {
		fmt.Println("Reading Value for Key :", key)
		fmt.Println("Reading Value for Value :", value)
		environmentstr = environmentstr + `"` + string(key) + `" : "` + fmt.Sprint(value) + `",`
	}

	jsonString := `{
		"applications" : {
			"` + rubyJSONObj.Appname + `" : {
				"type" : "perl",
				"limits" : {
					"timeout" : ` + strconv.Itoa(rubyJSONObj.Limits.Timeout) + `,
					"requests" : ` + strconv.Itoa(rubyJSONObj.Limits.Requests) + `
				},
				"processes" : {
					"max" : ` + strconv.Itoa(rubyJSONObj.Processes.Max) + `,
					"spare" : ` + strconv.Itoa(rubyJSONObj.Processes.Spare) + `,
					"idle_timeout" : ` + strconv.Itoa(rubyJSONObj.Processes.IdleTimeout) + `
				},
				"working_directory" : "` + rubyJSONObj.WorkingDirectory + `",
				"user" : "` + rubyJSONObj.User + `",
				"group" : "` + rubyJSONObj.Group + `",
				"environment" : {
					` + fmt.Sprint(environmentstr) + `
				},
				"script" : "` + rubyJSONObj.Script + `",
				"threads" : ` + strconv.Itoa(rubyJSONObj.Threads) + `
			}
		},
		"listeners" : {
			"*:` + strconv.Itoa(rubyJSONObj.Port) + `" : {
				"pass" : "applications/` + rubyJSONObj.Appname + `"
			}
		}
	}`
	/* UNCOMMENT FOR PRODUCTION */
	c.Pool.Broadcast <- websocket.Message{Type: 1, Body: ("Finished config for " + rubyJSONObj.Lang + "\n")}
	c.Pool.Broadcast <- websocket.Message{Type: 1, Body: jsonString}

	_ = ioutil.WriteFile((BUILDDIR + "/machines/builds/nginx-unit-" + rubyJSONObj.Lang + "/docker-entrypoint.d/config.json"), []byte(jsonString), 0644)
}
