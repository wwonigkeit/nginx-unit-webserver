package unit

//Limits struct used by the base structs
type Limits struct {
	Timeout  int `json:"timeout"`
	Requests int `json:"requests"`
}

//Processes struct ued by the base struct
type Processes struct {
	Max         int `json:"max"`
	Spare       int `json:"spare"`
	IdleTimeout int `json:"idle_timeout"`
}

//Options for PHP struct used for the JSON unmarshalling
type Options struct {
	Admin map[string]interface{} `json:"admin"`
	User  map[string]interface{} `json:"user"`
	File  string                 `json:"file"`
}

//Targets for PHP struct used for the JSON unmarshalling
type Targets struct {
	Script    string `json:"script"`
	Root      string `json:"root"`
	Reference string `json:"reference"`
	Index     string `json:"index"`
}

//Cloud struct for cloud deployment details used for the JSON unmarshalling
type Cloud struct {
	Platform    string `json:"platform"`
	MachineType string `json:"machinetype"`
}

//BaseStruct used for the JSON unmarshalling
type BaseStruct struct {
	Port             int                    `json:"port"`
	Limits           *Limits                `json:"limits"`
	Processes        *Processes             `json:"processes"`
	User             string                 `json:"user"`
	Group            string                 `json:"group"`
	Environment      map[string]interface{} `json:"environment"`
	Lang             string                 `json:"lang"`
	Repo             string                 `json:"repo"`
	Appname          string                 `json:"appname"`
	Cloud            *Cloud                 `json:"cloud"`
	WorkingDirectory string                 `json:"working_directory"`
}

//External struct used for JSON unmarshalling
type External struct {
	Port             int                    `json:"port"`
	Limits           *Limits                `json:"limits"`
	Processes        *Processes             `json:"processes"`
	User             string                 `json:"user"`
	Group            string                 `json:"group"`
	Environment      map[string]interface{} `json:"environment"`
	Lang             string                 `json:"lang"`
	Repo             string                 `json:"repo"`
	Appname          string                 `json:"appname"`
	Cloud            *Cloud                 `json:"cloud"`
	WorkingDirectory string                 `json:"working_directory"`
	Executable       string                 `json:"executable"`
	Arguments        []string               `json:"arguments"`
}

//Java struct used for the JSON unmarshalling
type Java struct {
	Port             int                    `json:"port"`
	Limits           *Limits                `json:"limits"`
	Processes        *Processes             `json:"processes"`
	User             string                 `json:"user"`
	Group            string                 `json:"group"`
	Environment      map[string]interface{} `json:"environment"`
	Lang             string                 `json:"lang"`
	Repo             string                 `json:"repo"`
	Appname          string                 `json:"appname"`
	Cloud            *Cloud                 `json:"cloud"`
	WorkingDirectory string                 `json:"working_directory"`
	Classpath        []string               `json:"classpath"`
	Options          []string               `json:"options"`
	Threads          int                    `json:"threads"`
	ThreadStackSize  int                    `json:"thread_stack_size"`
	Webapp           string                 `json:"webapp"`
}

//Perl struct used for the JSON unmarshalling
type Perl struct {
	Port             int                    `json:"port"`
	Limits           *Limits                `json:"limits"`
	Processes        *Processes             `json:"processes"`
	User             string                 `json:"user"`
	Group            string                 `json:"group"`
	Environment      map[string]interface{} `json:"environment"`
	Lang             string                 `json:"lang"`
	Repo             string                 `json:"repo"`
	Appname          string                 `json:"appname"`
	Cloud            *Cloud                 `json:"cloud"`
	WorkingDirectory string                 `json:"working_directory"`
	Threads          int                    `json:"threads"`
	ThreadStackSize  int                    `json:"thread_stack_size"`
	Script           string                 `json:"script"`
}

//PHP struct used for the JSON unmarshalling
type PHP struct {
	Port             int                    `json:"port"`
	Limits           *Limits                `json:"limits"`
	Processes        *Processes             `json:"processes"`
	User             string                 `json:"user"`
	Group            string                 `json:"group"`
	Environment      map[string]interface{} `json:"environment"`
	Lang             string                 `json:"lang"`
	Repo             string                 `json:"repo"`
	Appname          string                 `json:"appname"`
	Cloud            *Cloud                 `json:"cloud"`
	WorkingDirectory string                 `json:"working_directory"`
	Options          *Options               `json:"options"`
	Targets          []*Targets             `json:"targets"`
}

//Python struct used for the JSON unmarshalling
type Python struct {
	Port             int                    `json:"port"`
	Limits           *Limits                `json:"limits"`
	Processes        *Processes             `json:"processes"`
	User             string                 `json:"user"`
	Group            string                 `json:"group"`
	Environment      map[string]interface{} `json:"environment"`
	Lang             string                 `json:"lang"`
	Repo             string                 `json:"repo"`
	Appname          string                 `json:"appname"`
	Cloud            *Cloud                 `json:"cloud"`
	WorkingDirectory string                 `json:"working_directory"`
	Protocol         string                 `json:"protocol"`
	Threads          int                    `json:"threads"`
	Module           string                 `json:"module"`
	Callable         string                 `json:"callable"`
	Home             string                 `json:"home"`
	Path             string                 `json:"path"`
	ThreadStackSize  int                    `json:"thread_stack_size"`
}

//Ruby struct used for the JSON unmarshalling
type Ruby struct {
	Port             int                    `json:"port"`
	Limits           *Limits                `json:"limits"`
	Processes        *Processes             `json:"processes"`
	User             string                 `json:"user"`
	Group            string                 `json:"group"`
	Environment      map[string]interface{} `json:"environment"`
	Lang             string                 `json:"lang"`
	Repo             string                 `json:"repo"`
	Appname          string                 `json:"appname"`
	Cloud            *Cloud                 `json:"cloud"`
	WorkingDirectory string                 `json:"working_directory"`
	Threads          int                    `json:"threads"`
	Script           string                 `json:"script"`
}
