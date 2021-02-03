package unit

import (
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/wwonigkeit/nginx-unit-webserver/backend/websocket"
)

//EchoPerlJSONConfig echos the configuration for the Perl
//Unit startup configurations
/* UNCOMMENT FOR PRODUCTION */
func EchoPerlJSONConfig(perlJSONObj *Perl, c *websocket.Client) {
	// func EchoPerlJSONConfig(perlJSONObj *Perl) {
	var environmentstr string

	for key, value := range perlJSONObj.Environment {
		fmt.Println("Reading Value for Key :", key)
		fmt.Println("Reading Value for Value :", value)
		environmentstr = environmentstr + `"` + string(key) + `" : "` + fmt.Sprint(value) + `",`
	}

	jsonString := `{
		"applications" : {
			"` + perlJSONObj.Appname + `" : {
				"type" : "perl",
				"limits" : {
					"timeout" : ` + strconv.Itoa(perlJSONObj.Limits.Timeout) + `,
					"requests" : ` + strconv.Itoa(perlJSONObj.Limits.Requests) + `
				},
				"processes" : {
					"max" : ` + strconv.Itoa(perlJSONObj.Processes.Max) + `,
					"spare" : ` + strconv.Itoa(perlJSONObj.Processes.Spare) + `,
					"idle_timeout" : ` + strconv.Itoa(perlJSONObj.Processes.IdleTimeout) + `
				},
				"working_directory" : "` + perlJSONObj.WorkingDirectory + `",
				"user" : "` + perlJSONObj.User + `",
				"group" : "` + perlJSONObj.Group + `",
				"environment" : {
					` + fmt.Sprint(environmentstr) + `
				},
				"script" : "` + perlJSONObj.Script + `",
				"threads" : ` + strconv.Itoa(perlJSONObj.Threads) + `,
				"thread_stack_size" : ` + strconv.Itoa(perlJSONObj.ThreadStackSize) + `,
			}
		},
		"listeners" : {
			"*:` + strconv.Itoa(perlJSONObj.Port) + `" : {
				"pass" : "applications/` + perlJSONObj.Appname + `"
			}
		}
	}`
	/* UNCOMMENT FOR PRODUCTION */
	c.Pool.Broadcast <- websocket.Message{Type: 1, Body: ("Finished config for " + perlJSONObj.Lang + "\n")}
	c.Pool.Broadcast <- websocket.Message{Type: 1, Body: jsonString}

	_ = ioutil.WriteFile((BUILDDIR + "/machines/builds/nginx-unit-" + perlJSONObj.Lang + "/docker-entrypoint.d/config.json"), []byte(jsonString), 0644)
}
