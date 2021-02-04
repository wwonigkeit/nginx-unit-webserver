package unit

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os/exec"
	"strconv"
	"time"

	"github.com/wwonigkeit/nginx-unit-webserver/backend/websocket"
)

//RubyConfig echos the configuration for the Ruby
//Unit startup configurations
func RubyConfig(rubyJSONObj *Ruby, c *websocket.Client) {

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

	c.Pool.Broadcast <- websocket.Message{Type: 1, Body: ("Finished config for " + rubyJSONObj.Lang + "\n")}
	c.Pool.Broadcast <- websocket.Message{Type: 1, Body: jsonString}

	_ = ioutil.WriteFile((BUILDDIR + "/machines/builds/nginx-unit-" + rubyJSONObj.Lang + "/docker-entrypoint.d/config.json"), []byte(jsonString), 0644)
}

//BuildRubyImage pushes the machine to the appropriate cloud platform as an image
func BuildRubyImage(rubyJSONObj *Ruby, c *websocket.Client) (imagename string) {

	// Create a image name which is unique
	now := time.Now()
	timestamp := now.Unix()
	image := "nginx-unit-" + rubyJSONObj.Lang + "-" + strconv.FormatInt(timestamp, 10)

	// Provision the machine image to the appropriate cloud platform
	cmd := exec.Command("vorteil", "images", "provision", BUILDDIR+"/machines/builds/nginx-unit-"+rubyJSONObj.Lang+"/", BUILDDIR+"/templates/provisioners/"+rubyJSONObj.Cloud.Platform+".provisioner", "--program[0].env", `"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin","UNIT_VERSION=1.21.0-1~buster","GITHUB_REPO=`+rubyJSONObj.Repo+`","WORKINGDIR=`+rubyJSONObj.WorkingDirectory+`","SCRIPT=`+rubyJSONObj.Script+`"`, "--name", image, "--force")

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
		return image
	}

	c.Pool.Broadcast <- websocket.Message{Type: 1, Body: ("Succesfully provisioned the machine image")}
	return image
}
