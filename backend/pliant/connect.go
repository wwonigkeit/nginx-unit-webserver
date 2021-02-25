package pliant

import (
	"bytes"
	"crypto/tls"
	"log"
	"net/http"
	"strings"
)

//PLIANTURL is the Pliant end-point
const PLIANTURL string = "https://vorteil.pliant.io/api/v1/trigger/admin/User/Provisioning_Endpoint?sync=true&api_key=89b313dc-e117-4ff1-be90-5c500ab716e6&worker_group=default"

//Connect send the JSON string to Pliant endpoint
func Connect(jsonString string) (response *string, err error) {

	// initialize http client
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest(http.MethodPost, PLIANTURL, strings.NewReader(jsonString))

	if err != nil {
		log.Println(err)
		return nil, err
	}

	// set the request header Content-Type for json
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	newStr := buf.String()

	return &newStr, nil
}
