package unit

import (
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/wwonigkeit/nginx-unit-webserver/backend/websocket"
)

//EchoPythonJSONConfig echos the configuration for the Python
//Unit startup configurations
/* UNCOMMENT FOR PRODUCTION */
func EchoPythonJSONConfig(pythonJSONObj *Python, c *websocket.Client) {
	// func EchoPythonJSONConfig(pythonJSONObj *Python) {
	var environmentstr string

	for key, value := range pythonJSONObj.Environment {
		fmt.Println("Reading Value for Key :", key)
		fmt.Println("Reading Value for Value :", value)
		environmentstr = environmentstr + `"` + string(key) + `" : "` + fmt.Sprint(value) + `",`
	}

	jsonString := `{
		"applications" : {
			"` + pythonJSONObj.Appname + `" : {
				"type" : "perl",
				"limits" : {
					"timeout" : ` + strconv.Itoa(pythonJSONObj.Limits.Timeout) + `,
					"requests" : ` + strconv.Itoa(pythonJSONObj.Limits.Requests) + `
				},
				"processes" : {
					"max" : ` + strconv.Itoa(pythonJSONObj.Processes.Max) + `,
					"spare" : ` + strconv.Itoa(pythonJSONObj.Processes.Spare) + `,
					"idle_timeout" : ` + strconv.Itoa(pythonJSONObj.Processes.IdleTimeout) + `
				},
				"working_directory" : "` + pythonJSONObj.WorkingDirectory + `",
				"user" : "` + pythonJSONObj.User + `",
				"group" : "` + pythonJSONObj.Group + `",
				"environment" : {
					` + fmt.Sprint(environmentstr) + `
				},
				"module" : "` + pythonJSONObj.Module + `",
				"callable" : "` + pythonJSONObj.Callable + `",
				"home" : "` + pythonJSONObj.Home + `",
				"path" : "` + pythonJSONObj.Path + `",
				"protocol" : "` + pythonJSONObj.Protocol + `",
				"threads" : ` + strconv.Itoa(pythonJSONObj.Threads) + `,
				"thread_stack_size" : ` + strconv.Itoa(pythonJSONObj.ThreadStackSize) + `,
			}
		},
		"listeners" : {
			"*:` + strconv.Itoa(pythonJSONObj.Port) + `" : {
				"pass" : "applications/` + pythonJSONObj.Appname + `"
			}
		}
	}`
	/* UNCOMMENT FOR PRODUCTION */
	c.Pool.Broadcast <- websocket.Message{Type: 1, Body: ("Finished config for " + pythonJSONObj.Lang + "\n")}
	c.Pool.Broadcast <- websocket.Message{Type: 1, Body: jsonString}

	_ = ioutil.WriteFile((BUILDDIR + "/machines/builds/nginx-unit-" + pythonJSONObj.Lang + "/docker-entrypoint.d/config.json"), []byte(jsonString), 0644)
}
