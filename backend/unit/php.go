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

	"github.com/wwonigkeit/nginx-unit-webserver/backend/websocket"
)

//PhpConfig echos the configuration for the PHP
//Unit startup configurations
func PhpConfig(phpJSONObj *PHP, c *websocket.Client) {

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

	c.Pool.Broadcast <- websocket.Message{Type: 1, Body: ("Finished config for " + phpJSONObj.Lang + "\n")}
	c.Pool.Broadcast <- websocket.Message{Type: 1, Body: jsonString}

	_ = ioutil.WriteFile((BUILDDIR + "/machines/builds/nginx-unit-" + phpJSONObj.Lang + "/docker-entrypoint.d/config.json"), []byte(jsonString), 0644)
}

//BuildPhpImage pushes the machine to the appropriate cloud platform as an image
func BuildPhpImage(phpJSONObj *PHP, c *websocket.Client) (imagename string) {

	// Create a image name which is unique
	now := time.Now()
	timestamp := now.Unix()
	image := "nginx-unit-" + phpJSONObj.Lang + "-" + strconv.FormatInt(timestamp, 10)

	// Provision the machine image to the appropriate cloud platform
	cmd := exec.Command("vorteil", "images", "provision", BUILDDIR+"/machines/builds/nginx-unit-"+phpJSONObj.Lang+"/", BUILDDIR+"/templates/provisioners/"+phpJSONObj.Cloud.Platform+".provisioner", "--program[0].env", `"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin","UNIT_VERSION=1.21.0-1~buster","GITHUB_REPO=`+phpJSONObj.Repo+`","WORKINGDIR=`+phpJSONObj.WorkingDirectory+`","UNKNOWN=`+phpJSONObj.Targets[0].Root+`"`, "--name", image, "--force")

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
