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

//PerlConfig echos the configuration for the Perl
//Unit startup configurations
func PerlConfig(perlJSONObj *Perl, c *websocket.Client) {

	empJSON, _ := json.MarshalIndent(perlJSONObj, "", "  ")

	c.Pool.Broadcast <- websocket.Message{Type: 1, Body: ("Original config for " + perlJSONObj.Lang + "\n")}
	c.Pool.Broadcast <- websocket.Message{Type: 1, Body: string(empJSON)}

	var environmentstr string

	for key, value := range perlJSONObj.Environment {
		fmt.Println("Reading Value for Key :", key)
		fmt.Println("Reading Value for Value :", value)
		environmentstr = environmentstr + `"` + string(key) + `" : "` + fmt.Sprint(value) + `",`
	}

	jsonString := `
{
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

	pliantString := `
{
	"bodyData": {
		"lang" : "` + perlJSONObj.Lang + `",
		"repo" : "` + perlJSONObj.Repo + `",
		"package" : "` + DockerPackages[perlJSONObj.Lang] + `",
		"cloud" : {
			"platform" : "` + perlJSONObj.Cloud.Platform + `",
			"machinetype" : "` + perlJSONObj.Cloud.MachineType + `"
		},
		"initialconfig" : 
			` + jsonString + `
	}
}
`

	c.Pool.Broadcast <- websocket.Message{Type: 1, Body: ("Finished Pliant config for " + perlJSONObj.Lang + "\n")}
	c.Pool.Broadcast <- websocket.Message{Type: 1, Body: pliantString}
	c.Pool.Broadcast <- websocket.Message{Type: 1, Body: ("Sending to Pliant server\n")}

	resp, err := pliant.Connect(pliantString)

	if err != nil {
		c.Pool.Broadcast <- websocket.Message{Type: 1, Body: ("Error sending to Pliant server " + err.Error() + "\n")}
	} else {
		c.Pool.Broadcast <- websocket.Message{Type: 1, Body: ("Error sending to Pliant server " + *resp + "\n")}
	}

	c.Pool.Broadcast <- websocket.Message{Type: 1, Body: ("Finished config for " + perlJSONObj.Lang + "\n")}
	c.Pool.Broadcast <- websocket.Message{Type: 1, Body: jsonString}

	_ = ioutil.WriteFile((BUILDDIR + "/machines/builds/nginx-unit-" + perlJSONObj.Lang + "/docker-entrypoint.d/config.json"), []byte(jsonString), 0644)
}

//BuildPerlImage pushes the machine to the appropriate cloud platform as an image
func BuildPerlImage(perlJSONObj *Perl, c *websocket.Client) (imagename string) {

	// Create a image name which is unique
	now := time.Now()
	timestamp := now.Unix()
	image := "nginx-unit-" + perlJSONObj.Lang + "-" + strconv.FormatInt(timestamp, 10)

	// Provision the machine image to the appropriate cloud platform
	cmd := exec.Command("vorteil", "images", "provision", BUILDDIR+"/machines/builds/nginx-unit-"+perlJSONObj.Lang+"/", BUILDDIR+"/templates/provisioners/"+perlJSONObj.Cloud.Platform+".provisioner", "--program[0].env", `"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin","UNIT_VERSION=1.21.0-1~buster","GITHUB_REPO=`+perlJSONObj.Repo+`","WORKINGDIR=`+perlJSONObj.WorkingDirectory+`","SCRIPT=`+perlJSONObj.Script+`"`, "--name", image, "--force")

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
