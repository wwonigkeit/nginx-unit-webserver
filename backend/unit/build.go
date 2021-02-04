package unit

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"

	"github.com/wwonigkeit/nginx-unit-webserver/backend/websocket"
)

//Build executes the build for all pplications
//It takes the external unmarshalled JSON object and the client
//object to relay the build progress
func Build(baseJSONObj *BaseStruct, c *websocket.Client) {

	fmt.Println("Build the " + baseJSONObj.Lang + " platform")

	c.Pool.Broadcast <- websocket.Message{Type: 1, Body: ("Executing the " + baseJSONObj.Lang + " build\n")}
	c.Pool.Broadcast <- websocket.Message{Type: 1, Body: ("Pull the NGINX Unit " + baseJSONObj.Lang + " build\n")}

	//Delete the directory if it exists
	cmd := exec.Command("rm", "-rf", BUILDDIR+"/machines/builds/nginx-unit-"+baseJSONObj.Lang)
	err := cmd.Run()
	if err != nil {
		log.Println(err)
		return
	}

	//Convert the container for first use
	cmd = exec.Command("vorteil", "projects", "convert-container", DockerPackages[baseJSONObj.Lang], BUILDDIR+"/machines/builds/nginx-unit-"+baseJSONObj.Lang)
	c.Pool.Broadcast <- websocket.Message{Type: 1, Body: PrintWaitLine(cmd)}
	c.Pool.Broadcast <- websocket.Message{Type: 1, Body: ("Succesfully retrieved repository " + DockerPackages[baseJSONObj.Lang])}

	//Copy the necessary build files: build-<language>.sh
	cmd = exec.Command("cp", BUILDDIR+"/templates/builds/"+baseJSONObj.Lang+"/build-"+baseJSONObj.Lang+".sh", BUILDDIR+"/machines/builds/nginx-unit-"+baseJSONObj.Lang+"/")
	err = cmd.Run()
	if err != nil {
		log.Println(err)
		return
	}
	c.Pool.Broadcast <- websocket.Message{Type: 1, Body: ("Succesfully copied the build-" + baseJSONObj.Lang + ".sh file")}

	//Copy the necessary build files: default.vcfg
	cmd = exec.Command("cp", BUILDDIR+"/templates/builds/"+baseJSONObj.Lang+"/default.vcfg", BUILDDIR+"/machines/builds/nginx-unit-"+baseJSONObj.Lang+"/")
	err = cmd.Run()
	if err != nil {
		log.Println(err)
		return
	}
	c.Pool.Broadcast <- websocket.Message{Type: 1, Body: ("Succesfully copied the default.vcfg file")}

	//Copy the necessary build files: nginx.png
	cmd = exec.Command("cp", BUILDDIR+"/templates/builds/"+baseJSONObj.Lang+"/nginx.png", BUILDDIR+"/machines/builds/nginx-unit-"+baseJSONObj.Lang+"/")
	err = cmd.Run()
	if err != nil {
		log.Println(err)
		return
	}
	c.Pool.Broadcast <- websocket.Message{Type: 1, Body: ("Succesfully copied the nginx.png file")}

	//Copy the necessary build files: /usr/local/bin/docker-entrypoint.sh
	cmd = exec.Command("cp", BUILDDIR+"/templates/builds/"+baseJSONObj.Lang+"/usr/local/bin/docker-entrypoint.sh", BUILDDIR+"/machines/builds/nginx-unit-"+baseJSONObj.Lang+"/usr/local/bin/")
	err = cmd.Run()
	if err != nil {
		log.Println(err)
		return
	}
	c.Pool.Broadcast <- websocket.Message{Type: 1, Body: ("Succesfully copied the /usr/local/bin/docker-entrypoint.sh file")}

	c.Pool.Broadcast <- websocket.Message{Type: 1, Body: ("Finished build for " + baseJSONObj.Lang + "\n")}
	/*
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
	*/
}

//PrintWaitLine waits for the completion of the command using the "Run" argument
//and takes the exec.Cmd struct as the input. It passes back the combined error
//and output (stdout & stderr)
func PrintWaitLine(cmd *exec.Cmd) string {

	//out, err := cmd.CombinedOutput()
	var b bytes.Buffer
	cmd.Stderr = &b
	err := cmd.Run()
	if err != nil {
		log.Println(err)
	}

	return string(b.Bytes())

}
