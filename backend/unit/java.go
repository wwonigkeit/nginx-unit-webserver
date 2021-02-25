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

//JavaConfig echos the configuration for the Java
//Unit startup configurations
func JavaConfig(javaJSONObj *Java, c *websocket.Client) {

	empJSON, _ := json.MarshalIndent(javaJSONObj, "", "  ")

	c.Pool.Broadcast <- websocket.Message{Type: 1, Body: ("Original config for " + javaJSONObj.Lang + "\n")}
	c.Pool.Broadcast <- websocket.Message{Type: 1, Body: string(empJSON)}

	var environmentstr string

	for key, value := range javaJSONObj.Environment {
		fmt.Println("Reading Value for Key :", key)
		fmt.Println("Reading Value for Value :", value)
		environmentstr = environmentstr + `"` + string(key) + `" : "` + fmt.Sprint(value) + `",`
	}

	classpath, _ := json.Marshal(javaJSONObj.Classpath)
	options, _ := json.Marshal(javaJSONObj.Options)

	jsonString := `
{
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

	pliantString := `
{
	"bodyData": {
		"lang" : "` + javaJSONObj.Lang + `",
		"repo" : "` + javaJSONObj.Repo + `",
		"package" : "` + DockerPackages[javaJSONObj.Lang] + `",
		"cloud" : {
			"platform" : "` + javaJSONObj.Cloud.Platform + `",
			"machinetype" : "` + javaJSONObj.Cloud.MachineType + `"
		},
		"initialconfig" : 
			` + jsonString + `
	}
}
`

	c.Pool.Broadcast <- websocket.Message{Type: 1, Body: ("Finished Pliant config for " + javaJSONObj.Lang + "\n")}
	c.Pool.Broadcast <- websocket.Message{Type: 1, Body: pliantString}
	c.Pool.Broadcast <- websocket.Message{Type: 1, Body: ("Sending to Pliant server\n")}

	resp, err := pliant.Connect(pliantString)

	if err != nil {
		c.Pool.Broadcast <- websocket.Message{Type: 1, Body: ("Error sending to Pliant server " + err.Error() + "\n")}
	} else {
		c.Pool.Broadcast <- websocket.Message{Type: 1, Body: ("Error sending to Pliant server " + *resp + "\n")}
	}

	c.Pool.Broadcast <- websocket.Message{Type: 1, Body: ("Finished config for " + javaJSONObj.Lang + "\n")}
	c.Pool.Broadcast <- websocket.Message{Type: 1, Body: jsonString}

	_ = ioutil.WriteFile((BUILDDIR + "/machines/builds/nginx-unit-" + javaJSONObj.Lang + "/docker-entrypoint.d/config.json"), []byte(jsonString), 0644)
}

//BuildJavaImage pushes the machine to the appropriate cloud platform as an image
func BuildJavaImage(javaJSONObj *Java, c *websocket.Client) (imagename string) {

	// Create a image name which is unique
	now := time.Now()
	timestamp := now.Unix()
	image := "nginx-unit-" + javaJSONObj.Lang + "-" + strconv.FormatInt(timestamp, 10)

	// Provision the machine image to the appropriate cloud platform
	cmd := exec.Command("vorteil", "images", "provision", BUILDDIR+"/machines/builds/nginx-unit-"+javaJSONObj.Lang+"/", BUILDDIR+"/templates/provisioners/"+javaJSONObj.Cloud.Platform+".provisioner", "--program[0].env", `"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin","UNIT_VERSION=1.21.0-1~buster","GITHUB_REPO=`+javaJSONObj.Repo+`","WORKINGDIR=`+javaJSONObj.WorkingDirectory+`","WEBAPP=`+javaJSONObj.Webapp+`"`, "--name", image, "--force")

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
