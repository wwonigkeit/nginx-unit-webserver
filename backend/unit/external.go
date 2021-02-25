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
	"time"

	"github.com/wwonigkeit/nginx-unit-webserver/backend/pliant"
	"github.com/wwonigkeit/nginx-unit-webserver/backend/websocket"
)

//ExternalConfig echos the configuration for the Golang and NodeJs
//unit startup configurations
func ExternalConfig(extJSONObj *External, c *websocket.Client) {

	empJSON, _ := json.MarshalIndent(extJSONObj, "", "  ")

	c.Pool.Broadcast <- websocket.Message{Type: 1, Body: ("Original config for " + extJSONObj.Lang + "\n")}
	c.Pool.Broadcast <- websocket.Message{Type: 1, Body: string(empJSON)}

	var environmentstr string

	for key, value := range extJSONObj.Environment {
		environmentstr = environmentstr + `"` + string(key) + `" : "` + fmt.Sprint(value) + `",`
	}

	args, _ := json.Marshal(extJSONObj.Arguments)

	jsonString := `
{
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

	pliantString := `
{
	"bodyData": {
		"lang" : "` + extJSONObj.Lang + `",
		"repo" : "` + extJSONObj.Repo + `",
		"package" : "` + DockerPackages[extJSONObj.Lang] + `",
		"cloud" : {
			"platform" : "` + extJSONObj.Cloud.Platform + `",
			"machinetype" : "` + extJSONObj.Cloud.MachineType + `"
		},
		"initialconfig" : 
			` + jsonString + `
	}
}
`

	c.Pool.Broadcast <- websocket.Message{Type: 1, Body: ("Finished Pliant config for " + extJSONObj.Lang + "\n")}
	c.Pool.Broadcast <- websocket.Message{Type: 1, Body: pliantString}
	c.Pool.Broadcast <- websocket.Message{Type: 1, Body: ("Sending to Pliant server\n")}

	resp, err := pliant.Connect(pliantString)

	if err != nil {
		c.Pool.Broadcast <- websocket.Message{Type: 1, Body: ("Error sending to Pliant server " + err.Error() + "\n")}
	} else {
		c.Pool.Broadcast <- websocket.Message{Type: 1, Body: ("Error sending to Pliant server " + *resp + "\n")}
	}

	c.Pool.Broadcast <- websocket.Message{Type: 1, Body: ("Finished config for " + extJSONObj.Lang + "\n")}
	c.Pool.Broadcast <- websocket.Message{Type: 1, Body: jsonString}

	_ = ioutil.WriteFile((BUILDDIR + "/machines/builds/nginx-unit-" + extJSONObj.Lang + "/docker-entrypoint.d/config.json"), []byte(jsonString), 0644)
}

//BuildExternalImage pushes the machine to the appropriate cloud platform as an image
func BuildExternalImage(extJSONObj *External, c *websocket.Client) (imagename string) {

	// Create a image name which is unique
	now := time.Now()
	timestamp := now.Unix()
	image := "nginx-unit-" + extJSONObj.Lang + "-" + strconv.FormatInt(timestamp, 10)

	// Provision the machine image to the appropriate cloud platform
	cmd := exec.Command("vorteil", "images", "provision", BUILDDIR+"/machines/builds/nginx-unit-"+extJSONObj.Lang+"/", BUILDDIR+"/templates/provisioners/"+extJSONObj.Cloud.Platform+".provisioner", "--program[0].env", `"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin","UNIT_VERSION=1.21.0-1~buster","GITHUB_REPO=`+extJSONObj.Repo+`","WORKINGDIR=`+extJSONObj.WorkingDirectory+`","EXECUTABLE=`+extJSONObj.Executable+`"`, "--name", image, "--force")

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
		//time.Sleep(5 * time.Second)
		c.Pool.Broadcast <- websocket.Message{Type: 1, Body: s.Text()}
	}

	if err := cmd.Wait(); err != nil {
		log.Println(err)
		return image
	}

	c.Pool.Broadcast <- websocket.Message{Type: 1, Body: ("Succesfully provisioned the machine image")}
	return image
}
