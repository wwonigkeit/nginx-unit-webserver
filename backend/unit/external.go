package unit

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os/exec"
	"strconv"

	"github.com/wwonigkeit/nginx-unit-webserver/backend/websocket"
)

//EchoExternalJSONConfig echos the configuration for the Golang and NodeJs
//unit startup configurations
/* UNCOMMENT FOR PRODUCTION */
func EchoExternalJSONConfig(extJSONObj *External, c *websocket.Client) {
	//func EchoExternalJSONConfig(extJSONObj *External) {
	var environmentstr string

	for key, value := range extJSONObj.Environment {
		environmentstr = environmentstr + `"` + string(key) + `" : "` + fmt.Sprint(value) + `",`
	}

	args, _ := json.Marshal(extJSONObj.Arguments)

	jsonString := `{
		"applications" : {
			"` + extJSONObj.Appname + `" : {
				"type" : "external",
				"limits" : {
					"timeout" : ` + strconv.Itoa(extJSONObj.Limits.Timeout) + `,
					"requests" : ` + strconv.Itoa(extJSONObj.Limits.Requests) + `
				},
				"processes" : {
					"max" : ` + strconv.Itoa(extJSONObj.Processes.Max) + `,
					"spare" : ` + strconv.Itoa(extJSONObj.Processes.Spare) + `,
					"idle_timeout" : ` + strconv.Itoa(extJSONObj.Processes.IdleTimeout) + `
				},
				"working_directory" : "` + extJSONObj.WorkingDirectory + `",
				"user" : "` + extJSONObj.User + `",
				"group" : "` + extJSONObj.Group + `",
				"environment" : {
					` + fmt.Sprint(environmentstr) + `
				},
				"executable" : "` + extJSONObj.Executable + `",
				"arguments" : ` + string(args) + `
			}
		},
		"listeners" : {
			"*:` + strconv.Itoa(extJSONObj.Port) + `" : {
				"pass" : "applications/` + extJSONObj.Appname + `"
			}
		}
	}`
	/* UNCOMMENT FOR PRODUCTION */
	c.Pool.Broadcast <- websocket.Message{Type: 1, Body: ("Finished config for " + extJSONObj.Lang + "\n")}
	c.Pool.Broadcast <- websocket.Message{Type: 1, Body: jsonString}

	_ = ioutil.WriteFile((BUILDDIR + "/machines/builds/nginx-unit-" + extJSONObj.Lang + "/docker-entrypoint.d/config.json"), []byte(jsonString), 0644)
}

//BuildExternalMachine pushes the machine to the appropriate cloud platform as an image
func BuildExternalMachine(extJSONObj *External, c *websocket.Client) {
	//Provision the machine image to the appropriate cloud platform
	cmd := exec.Command("vorteil", "images", "provision", BUILDDIR+"/machines/builds/nginx-unit-"+extJSONObj.Lang+"/", BUILDDIR+"/templates/provisioners/"+extJSONObj.Cloud+".provisioner", "--program[0].env", `"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin","UNIT_VERSION=1.21.0-1~buster","GITHUB_REPO=`+extJSONObj.Repo+`","WORKINGDIR=`+extJSONObj.WorkingDirectory+`","EXECUTABLE=`+extJSONObj.Executable+`"`, "--name", "nginx-unit-go", "--force")

	stdout, err := cmd.StdoutPipe()

	if err != nil {
		log.Println(err)
		return
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Println(err)
		return
	}
	if err = cmd.Start(); err != nil {
		log.Println(err)
		return
	}

	s := bufio.NewScanner(io.MultiReader(stdout, stderr))
	for s.Scan() {
		c.Pool.Broadcast <- websocket.Message{Type: 1, Body: s.Text()}
	}

	if err := cmd.Wait(); err != nil {
		log.Println(err)
		return
	}

	c.Pool.Broadcast <- websocket.Message{Type: 1, Body: ("Succesfully provisioned the machine image")}
}
