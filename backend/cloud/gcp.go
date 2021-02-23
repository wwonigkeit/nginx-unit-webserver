package cloud

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/wwonigkeit/nginx-unit-webserver/backend/websocket"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/option"
)

//BuildGcpInstance takes the machine image and creates the running virtual machine
func BuildGcpInstance(image string, machine string, c *websocket.Client, port int) {
	c.Pool.Broadcast <- websocket.Message{Type: 1, Body: ("Creating a instance on Google Cloud Platform\n")}
	ctx := context.Background()
	computeService, err := compute.NewService(ctx, option.WithCredentialsFile(GCPJSON))

	// Create the GCP instance configuration
	newInstance := &compute.Instance{
		MachineType: "zones/australia-southeast1-b/machineTypes/" + machine,
		Name:        image,
		Disks: []*compute.AttachedDisk{
			{
				Boot: true,
				InitializeParams: &compute.AttachedDiskInitializeParams{
					SourceImage: "global/images/" + image,
				},
			},
		},
		NetworkInterfaces: []*compute.NetworkInterface{
			{
				AccessConfigs: []*compute.AccessConfig{
					{
						Type: "ONE_TO_ONE_NAT",
					},
				},
				Network: "global/networks/default",
			},
		},
	}

	// Get the GCP project details from the JSON configuration file
	jsonFile, err := os.Open(GCPJSON)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var gcpconfig map[string]interface{}
	json.Unmarshal([]byte(byteValue), &gcpconfig)

	fmt.Println(gcpconfig["project_id"])

	resp, err := computeService.Instances.Insert((fmt.Sprint(gcpconfig["project_id"])), GCPZONE, newInstance).Context(ctx).Do()
	if err != nil {
		log.Fatal(err)
	}

OpLoop:
	for {
		time.Sleep(2 * time.Second)
		op, err := computeService.ZoneOperations.Get((fmt.Sprint(gcpconfig["project_id"])), GCPZONE, resp.Name).Do()
		if err != nil {
			fmt.Println("Failed to get op %s: %v", resp.Name, err)
			return
		}
		switch op.Status {
		case "PENDING", "RUNNING":
			continue
		case "DONE":
			if op.Error != nil {
				for _, operr := range op.Error.Errors {
					log.Printf("failed to create instance %s in zone %s: %v", newInstance.Name, GCPZONE, operr.Code)
					return
				}
				return
			}
			break OpLoop
		default:
			fmt.Println("Unknown create status %q: %+v", op.Status, op)
			return
		}
	}

	inst, err := computeService.Instances.Get((fmt.Sprint(gcpconfig["project_id"])), GCPZONE, image).Do()

	if err != nil {
		fmt.Errorf("Error getting instance %s details after creation: %v", image, err)
	}

	// Finds its internal and/or external IP addresses.
	_, extIP := instanceIPs(inst)

	var start int64 = 0

SerialLoop:
	for {
		time.Sleep(1 * time.Second)
		respSerial, _ := computeService.Instances.GetSerialPortOutput((fmt.Sprint(gcpconfig["project_id"])), GCPZONE, image).Start(start).Context(ctx).Do()

		fmt.Println("Start", start, "Next", respSerial.Next)

		// We will check whether the NGINX Unit instance has completed startup by getting the /config from
		// the instance on port 8080
		timeout := time.Duration(1 * time.Second)

		client := http.Client{
			Timeout: timeout,
		}

		resp, err := client.Get("http://" + extIP + ":8080/config")

		if err != nil {
			c.Pool.Broadcast <- websocket.Message{Type: 1, Body: respSerial.Contents}
		} else {
			c.Pool.Broadcast <- websocket.Message{Type: 1, Body: respSerial.Contents}
			defer resp.Body.Close()
			bodyBytes, err := ioutil.ReadAll(resp.Body)

			if err != nil {
				log.Fatal(err)
			}

			mesg :=
				"Deployment of the nginx-unit-go instance on Google Cloud Platform has been completed.\n" +
					"The external IP adddress for the machine is:\n" +
					"\n" +
					"IP address: " + extIP + "\n" +
					"\n" +
					"You can verify the configuration of the instance @\n" +
					"\n" +
					"\thttp://" + extIP + ":" + strconv.Itoa(port) + "/\n" +
					"\n" +
					"if your application is available at this location. Alternatively the configuration for the Unit instance\n" +
					"can be changed or verified using the following URL:\n" +
					"\n" +
					"\thttp://" + extIP + ":8080/config\n" +
					"\n" +
					string(bodyBytes)

			c.Pool.Broadcast <- websocket.Message{Type: 1, Body: mesg}

			break SerialLoop
		}
		start = respSerial.Next
	}
}

func instanceIPs(inst *compute.Instance) (intIP, extIP string) {
	for _, iface := range inst.NetworkInterfaces {
		if strings.HasPrefix(iface.NetworkIP, "10.") {
			intIP = iface.NetworkIP
		}
		for _, accessConfig := range iface.AccessConfigs {
			if accessConfig.Type == "ONE_TO_ONE_NAT" {
				extIP = accessConfig.NatIP
			}
		}
	}
	return
}
