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

//PythonConfig echos the configuration for the Python
//Unit startup configurations
func PythonConfig(pythonJSONObj *Python, c *websocket.Client) {

	empJSON, _ := json.MarshalIndent(pythonJSONObj, "", "  ")

	c.Pool.Broadcast <- websocket.Message{Type: 1, Body: ("Original config for " + pythonJSONObj.Lang + "\n")}
	c.Pool.Broadcast <- websocket.Message{Type: 1, Body: string(empJSON)}

	var environmentstr string

	for key, value := range pythonJSONObj.Environment {
		fmt.Println("Reading Value for Key :", key)
		fmt.Println("Reading Value for Value :", value)
		environmentstr = environmentstr + `"` + string(key) + `" : "` + fmt.Sprint(value) + `",`
	}

	jsonString := `{
		"applications" : {
			"` + pythonJSONObj.Appname + `" : {
				"type" : "python",
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

	pliantString := `
{
	"bodyData": {
		"lang" : "` + pythonJSONObj.Lang + `",
		"repo" : "` + pythonJSONObj.Repo + `",
		"package" : "` + DockerPackages[pythonJSONObj.Lang] + `",
		"cloud" : {
			"platform" : "` + pythonJSONObj.Cloud.Platform + `",
			"machinetype" : "` + pythonJSONObj.Cloud.MachineType + `"
		},
		"initialconfig" : 
			` + jsonString + `
	}
}
`

	c.Pool.Broadcast <- websocket.Message{Type: 1, Body: ("Finished Pliant config for " + pythonJSONObj.Lang + "\n")}
	c.Pool.Broadcast <- websocket.Message{Type: 1, Body: pliantString}
	c.Pool.Broadcast <- websocket.Message{Type: 1, Body: ("Sending to Pliant server\n")}

	resp, err := pliant.Connect(pliantString)

	if err != nil {
		c.Pool.Broadcast <- websocket.Message{Type: 1, Body: ("Error sending to Pliant server " + err.Error() + "\n")}
	} else {
		c.Pool.Broadcast <- websocket.Message{Type: 1, Body: ("Error sending to Pliant server " + *resp + "\n")}
	}

	c.Pool.Broadcast <- websocket.Message{Type: 1, Body: ("Finished config for " + pythonJSONObj.Lang + "\n")}
	c.Pool.Broadcast <- websocket.Message{Type: 1, Body: jsonString}

	_ = ioutil.WriteFile((BUILDDIR + "/machines/builds/nginx-unit-" + pythonJSONObj.Lang + "/docker-entrypoint.d/config.json"), []byte(jsonString), 0644)
}

//BuildPythonImage pushes the machine to the appropriate cloud platform as an image
func BuildPythonImage(pythonJSONObj *Python, c *websocket.Client) (imagename string) {

	// Create a image name which is unique
	now := time.Now()
	timestamp := now.Unix()
	image := "nginx-unit-" + pythonJSONObj.Lang + "-" + strconv.FormatInt(timestamp, 10)

	// Provision the machine image to the appropriate cloud platform
	cmd := exec.Command("vorteil", "images", "provision", BUILDDIR+"/machines/builds/nginx-unit-"+pythonJSONObj.Lang+"/", BUILDDIR+"/templates/provisioners/"+pythonJSONObj.Cloud.Platform+".provisioner", "--program[0].env", `"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin","UNIT_VERSION=1.21.0-1~buster","GITHUB_REPO=`+pythonJSONObj.Repo+`","WORKINGDIR=`+pythonJSONObj.WorkingDirectory+`","MODULE=`+pythonJSONObj.Module+`"`, "--name", image, "--force")

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
