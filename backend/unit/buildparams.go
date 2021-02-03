package unit

//DockerPackages contains the verified versions for the Unit Docker packages used for conversion
var DockerPackages = map[string]string{
	"go":     "nginx/unit:1.21.0-go1.11-dev",
	"nodejs": "nginx/unit:1.21.0-minimal",
	"java":   "nginx/unit:1.21.0-jsc11",
	"perl":   "nginx/unit:1.21.0-perl5.28",
	"php":    "nginx/unit:1.21.0-php7.3",
	"python": "nginx/unit:1.21.0-python3.7",
	"ruby":   "unit:1.21.0-ruby2.5",
}

//BUILDDIR is the working directory for the builds
const BUILDDIR string = "/home/wilhelmwonigkeit/Projects/src/github.com/wwonigkeit/nginx-unit-webserver/builds"
