package cloud

import "github.com/wwonigkeit/nginx-unit-webserver/backend/unit"

//GCPJSON contains the JSON Google Cloud Platform authentication credentials
const GCPJSON string = unit.BUILDDIR + "/templates/provisioners/vorteil-demo-gcp.json"

//GCPZONE contains the JSON Google Cloud Platform zone to use
const GCPZONE string = "australia-southeast1-b"
