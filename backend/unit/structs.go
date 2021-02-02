package unit

//Limits struct used by the base structs
type Limits struct {
	Timeout  int64 `json:"timeout"`
	Requests int64 `json:"requests"`
}

//Processes struct ued by the base struct
type Processes struct {
	Max         int64 `json:"max"`
	Spare       int64 `json:"spare"`
	IdleTimeout int64 `json:"idle_timeout"`
}

//Options for PHP struct used for the JSON unmarshalling
type Options struct {
	Admin string `json:"admin"`
	User  string `json:"user"`
	File  string `json:"file"`
}

//Targets for PHP struct used for the JSON unmarshalling
type Targets struct {
	Script    string `json:"script"`
	Root      string `json:"root"`
	Reference string `json:"reference"`
	Index     string `json:"index"`
}

//BaseStruct used for the JSON unmarshalling
type BaseStruct struct {
	Port             int64      `json:"port"`
	Limits           *Limits    `json:"limits"`
	Processes        *Processes `json:"processes"`
	User             string     `json:"user"`
	Group            string     `json:"group"`
	Environment      string     `json:"environment"`
	Lang             string     `json:"lang"`
	Repo             string     `json:"repo"`
	Appname          string     `json:"appname"`
	Cloud            string     `json:"cloud"`
	WorkingDirectory string     `json:"working_directory"`
}

//External struct used for JSON unmarshalling
type External struct {
	BaseStruct
	Executable string   `json:"executable"`
	Arguments  []string `json:"arguments"`
}

//Java struct used for the JSON unmarshalling
type Java struct {
	BaseStruct
	Classpath       []string `json:"classpath"`
	Options         []string `json:"options"`
	Threads         int64    `json:"threads"`
	ThreadStackSize int64    `json:"thread_stack_size"`
	Webapp          string   `json:"webapp"`
}

//Perl struct used for the JSON unmarshalling
type Perl struct {
	BaseStruct
	Threads int64  `json:"threads"`
	Script  string `json:"script"`
}

//PHP struct used for the JSON unmarshalling
type PHP struct {
	BaseStruct
	Options *Options   `json:"options"`
	Targets []*Targets `json:"targets"`
}

//Python struct used for the JSON unmarshalling
type Python struct {
	BaseStruct
	Protocol        string `json:"protocol"`
	Threads         int64  `json:"threads"`
	Module          string `json:"module"`
	Callable        string `json:"callable"`
	Home            string `json:"home"`
	Path            string `json:"path"`
	ThreadStackSize int64  `json:"thread_stack_size"`
}

//Ruby struct used for the JSON unmarshalling
type Ruby struct {
	BaseStruct
	Threads int64  `json:"threads"`
	Script  string `json:"script"`
}
