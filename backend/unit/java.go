package unit

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/wwonigkeit/nginx-unit-webserver/backend/websocket"
)

//EchoJavaJSONConfig echos the configuration for the Java
//Unit startup configurations
/* UNCOMMENT FOR PRODUCTION */
func EchoJavaJSONConfig(javaJSONObj *Java, c *websocket.Client) {
	// func EchoJavaJSONConfig(javaJSONObj *Java) {
	var environmentstr string

	for key, value := range javaJSONObj.Environment {
		fmt.Println("Reading Value for Key :", key)
		fmt.Println("Reading Value for Value :", value)
		environmentstr = environmentstr + `"` + string(key) + `" : "` + fmt.Sprint(value) + `",`
	}

	classpath, _ := json.Marshal(javaJSONObj.Classpath)
	options, _ := json.Marshal(javaJSONObj.Options)

	jsonString := `{
		"applications" : {
			"` + javaJSONObj.Appname + `" : {
				"type" : "java",
				"limits" : {
					"timeout" : ` + strconv.Itoa(javaJSONObj.Limits.Timeout) + `,
					"requests" : ` + strconv.Itoa(javaJSONObj.Limits.Requests) + `
				},
				"processes" : {
					"max" : ` + strconv.Itoa(javaJSONObj.Processes.Max) + `,
					"spare" : ` + strconv.Itoa(javaJSONObj.Processes.Spare) + `,
					"idle_timeout" : ` + strconv.Itoa(javaJSONObj.Processes.IdleTimeout) + `
				},
				"working_directory" : "` + javaJSONObj.WorkingDirectory + `",
				"user" : "` + javaJSONObj.User + `",
				"group" : "` + javaJSONObj.Group + `",
				"environment" : {
					` + fmt.Sprint(environmentstr) + `
				},
				"webapp" : "` + javaJSONObj.Webapp + `",
				"classpath" : ` + string(classpath) + `,
				"options" : ` + string(options) + `,
				"threads" : ` + strconv.Itoa(javaJSONObj.Threads) + `,
				"thread_stack_size" : ` + strconv.Itoa(javaJSONObj.ThreadStackSize) + `,
			}
		},
		"listeners" : {
			"*:` + strconv.Itoa(javaJSONObj.Port) + `" : {
				"pass" : "applications/` + javaJSONObj.Appname + `"
			}
		}
	}`
	/* UNCOMMENT FOR PRODUCTION */
	c.Pool.Broadcast <- websocket.Message{Type: 1, Body: ("Finished config for " + javaJSONObj.Lang + "\n")}
	c.Pool.Broadcast <- websocket.Message{Type: 1, Body: jsonString}

	_ = ioutil.WriteFile((BUILDDIR + "/machines/builds/nginx-unit-" + javaJSONObj.Lang + "/docker-entrypoint.d/config.json"), []byte(jsonString), 0644)
}
