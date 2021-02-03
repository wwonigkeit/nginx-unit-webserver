package unit

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/wwonigkeit/nginx-unit-webserver/backend/websocket"
)

//EchoPHPJSONConfig echos the configuration for the PHP
//Unit startup configurations
/* UNCOMMENT FOR PRODUCTION */
func EchoPHPJSONConfig(phpJSONObj *PHP, c *websocket.Client) {
	// func EchoPHPJSONConfig(phpJSONObj *PHP) {

	var environmentstr string
	for key, value := range phpJSONObj.Environment {
		fmt.Println("Reading Value for Key :", key)
		fmt.Println("Reading Value for Value :", value)
		environmentstr = environmentstr + `"` + string(key) + `" : "` + fmt.Sprint(value) + `",`
	}

	var adminoptionstr string
	for key, value := range phpJSONObj.Options.Admin {
		//fmt.Println("Reading Value for Key :", key)
		//fmt.Println("Reading Value for Value :", value)
		adminoptionstr = adminoptionstr + `"` + string(key) + `" : "` + fmt.Sprint(value) + `",`
	}

	var useroptionstr string
	for key, value := range phpJSONObj.Options.User {
		//fmt.Println("Reading Value for Key :", key)
		//fmt.Println("Reading Value for Value :", value)
		useroptionstr = useroptionstr + `"` + string(key) + `" : "` + fmt.Sprint(value) + `",`
	}

	m := make(map[string]map[string]string)

	for _, target := range phpJSONObj.Targets {
		//fmt.Println("Reading Value for Reference :", target.Reference, " Index :", target.Index, " Root :", target.Root, " Script :", target.Script)
		m[target.Reference] = make(map[string]string)
		if target.Index != "" {
			m[target.Reference]["root"] = target.Root
			m[target.Reference]["index"] = target.Index
		} else {
			m[target.Reference]["root"] = target.Root
			m[target.Reference]["script"] = target.Script
		}

	}

	targetstr, _ := json.Marshal(m)

	jsonString := `{
		"applications" : {
			"` + phpJSONObj.Appname + `" : {
				"type" : "php",
				"limits" : {
					"timeout" : ` + strconv.Itoa(phpJSONObj.Limits.Timeout) + `,
					"requests" : ` + strconv.Itoa(phpJSONObj.Limits.Requests) + `
				},
				"processes" : {
					"max" : ` + strconv.Itoa(phpJSONObj.Processes.Max) + `,
					"spare" : ` + strconv.Itoa(phpJSONObj.Processes.Spare) + `,
					"idle_timeout" : ` + strconv.Itoa(phpJSONObj.Processes.IdleTimeout) + `
				},
				"working_directory" : "` + phpJSONObj.WorkingDirectory + `",
				"user" : "` + phpJSONObj.User + `",
				"group" : "` + phpJSONObj.Group + `",
				"environment" : {
					` + fmt.Sprint(environmentstr) + `
				},
				"options" : {
					"file" : "` + phpJSONObj.Options.File + `",
					"admin" : {
						` + fmt.Sprint(adminoptionstr) + `
					},
					"user" : {
						` + fmt.Sprint(useroptionstr) + `
					}
				},
				"targets" : {
					` + string(targetstr) + `
				}
			}
		},
		"listeners" : {
			"*:` + strconv.Itoa(phpJSONObj.Port) + `" : {
				"pass" : "applications/` + phpJSONObj.Targets[0].Reference + `/` + phpJSONObj.Appname + `"
			}
		}
	}`
	/* UNCOMMENT FOR PRODUCTION */
	c.Pool.Broadcast <- websocket.Message{Type: 1, Body: ("Finished config for " + phpJSONObj.Lang + "\n")}
	c.Pool.Broadcast <- websocket.Message{Type: 1, Body: jsonString}

	_ = ioutil.WriteFile((BUILDDIR + "/machines/builds/nginx-unit-" + phpJSONObj.Lang + "/docker-entrypoint.d/config.json"), []byte(jsonString), 0644)
}
