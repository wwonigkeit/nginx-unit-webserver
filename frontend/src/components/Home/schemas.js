exports.customFormats = {
   "repo":"https://.*.git",
   "perl-script": "^[^/].*.psgi",
   "external-executable": "^[^/].*",
   "java-webapp": "^[^/].*",
   "ruby-script": "^[^/].*.ru",
   "working-directory": "^/.*"
}

exports.baseSchema = {
   "title":"NGINX Unit configuration options",
   "type":"object",
   "properties":{
      "lang":{
         "type":"string",
         "title":"Application Type",
         "enum":[
            "go",
            "nodejs",
            "java",
            "perl",
            "php",
            "python",
            "ruby"
         ],
      },
      "appname": {
         "title": "Application name",
         "description":"The application name used to reference the Unit configuration instance",
         "type":"string"
      },
      "repo": {
         "title":"Repository adddress",
         "description":"Code repository address e.g. https://github.com/username/repo.git - currently only supports GitHub",
         "type":"string",
         "format":"repo"
      },
      "cloud": {
         "title": "Cloud deployment details",
         "type": "object",
         "properties":{
            "platform": {
               "type":"string",
               "title":"Cloud selection",
               "description":"Select cloud for provisioning (currently Amazon Web Services[aws], Google Cloud Platform[gcp] or Microsoft Azure [azure])",
               "enum":[
                  "aws",
                  "gcp",
                  "azure"
               ],
               "enumNames": ["Amazon Web Services", "Google Cloud Platform", "Azure"],
               "default":"gcp"
            },
         },
         "required": [
            "platform"
         ],
         "dependencies": {
            "platform": {
               "oneOf": [
                  {
                     "properties":{
                        "platform":{
                           "enum": [ "aws" ]
                        },
                        "machinetype":{
                           "type": "string",
                           "title": "Machine type",
                           "description": "Select the machine type from the Amazon Web Services offerings below",
                           "enum": ["t2.nano","t2.micro","t2.small","t2.medium","t2.large","t2.xlarge"],
                           "enumNames": ["t2.nano (1vCPU, 0.5GB mem)", "t2.micro (1vCPU, 1GB mem)", "t2.small (1vCPU, 2GB mem)", "t2.medium (2vCPU, 4GB mem)", "t2.large (2vCPU, 8GB mem)", "t2.xlarge (4vCPU, 16GB mem)"]
                        },
                     },
                     "required": [
                        "machinetype"
                     ],
                  },
                  {
                     "properties":{
                        "platform":{
                           "enum": [ "gcp" ]
                        },
                        "machinetype":{
                           "type": "string",
                           "title": "Machine type",
                           "description": "Select the machine type from the Google Cloud Platform offerings below",
                           "enum": ["f1-micro","g1-small","n1-standard-1","n1-standard-2","n1-standard-4"],
                           "enumNames": ["f1-micro (1vCPU, 0.6GB mem)", "g1-small (1vCPU, 1.7GB mem)", "n1-standard-1 (1vCPU, 3.75GB mem)", "n1-standard-2 (2vCPU, 7.5GB mem)","n1-standard-4 (4vCPU, 15GB)"]
                        },
                     },
                     "required": [
                        "machinetype"
                     ],
                  },
                  {
                     "properties":{
                        "platform":{
                           "enum": [ "azure" ]
                        },
                        "machinetype":{
                           "type": "string",
                           "title": "Machine type",
                           "description": "Select the machine type from the Google Cloud Platform offerings below",
                           "enum": ["Standard_B1ls","Standard_B1ms","Standard_B2s","Standard_B2ms","Standard_B4ms"],
                           "enumNames": ["Standard_B1ls1	(1vCPU, 0.5GB mem)", "Standard_B1ms (1vCPU, 2GB mem)", "Standard_B2s (2vCPU, 4GB mem)", "Standard_B2ms (2vCPU, 8GB mem)","Standard_B4ms (4vCPU, 16GB)"]
                        },
                     },
                     "required": [
                        "machinetype"
                     ],
                  }
               ]
            }

         }
      },
      "port": {
         "type":"number",
         "default": 80,
         "title": "Listener port",
         "description": "TCP port on which the application will be configured to listen"
      },
      "limits":{
         "title":"Request Limits",
         "type":"object",
         "properties":{
            "timeout":{
               "description":"Request timeout in seconds",
               "title":"Timeout",
               "type":"number",
               "default": 10
            },
            "requests":{
               "description":"Maximum number of requests allowed to serve",
               "title":"Maximum requests",
               "type":"number",
               "default": 100
            }
         }
      },
      "processes":{
         "title":"Process Management",
         "type":"object",
         "properties":{
            "max":{
               "title": "Maximum processes",
               "description":"Maximum number of application processes",
               "type":"number",
               "default": 1
            },
            "spare":{
               "title":"Minimum idle processes",
               "description":"Minimum number of idle processes that Unit tries to reserve for an app",
               "type":"number",
               "default": 1
            },
            "idle_timeout":{
               "title":"Idle process timeout",
               "description":"Time in seconds that Unit waits before terminating an idle process which exceeds spare",
               "type":"number",
               "default": 20
            }
         }
      },
      "working_directory":{
         "type":"string",
         "title":"Working Directory",
         "description":"The applications working directory",
         "format": "working-directory"
      },
      "user":{
         "type":"string",
         "title":"Username",
         "description": "Username that runs the process",
         "default":"root"
      },
      "group":{
         "type":"string",
         "title":"Group",
         "description": "Group name that runs the process",
         "default":"root"
      },
      "environment":{
         "title":"Environmental variables",
         "description":"Environment variables to be passed to the application",
         "type":"object",
         "additionalProperties": {
            "type": "string"
         }
      }
   },
   "required": [
      "lang", "working_directory", "appname","repo"
   ],
   "dependencies": {
      "lang":{
         "oneOf": [
            {
               "properties":{
                  "lang":{
                     "enum": [
                        "go", "nodejs"
                     ]
                  },
                  "executable":{
                     "type":"string",
                     "title":"Exectuable name",
                     "description":"Pathname of the application, relative to working directory",
                     "format": "external-executable"
                  },
                  "arguments":{
                     "type":"array",
                     "title":"Command line arguments",
                     "description":"Command line arguments to be passed to the application equivalent to /app --tmp-files /tmp/go-cache",
                     "items":{
                        "type":"string",
                        "title":"Argument Name"
                     }
                  }
               },
               "required": [
                     "executable"
               ]
            },
            {
               "properties":{
                  "lang":{
                     "enum": [
                        "java"
                     ]
                  },
                  "webapp":{
                     "type":"string",
                     "title": "Application name",
                     "description":"Pathname and name of the application’s packaged or unpackaged .war file: e.g. helloworld.war relative to the working directory",
                     "format": "java-webapp"
                  },
                  "classpath":{
                     "type":"array",
                     "title": "Classpath details",
                     "description":"Array of paths to your app’s required libraries (may list directories or .jar files)",
                     "items":{
                        "type":"string"
                     }
                  },
                  "options":{
                     "type":"array",
                     "title":"Java JVM options",
                     "description":"Array of strings defining JVM runtime options",
                     "items":{
                        "type":"string"
                     }
                  },
                  "threads":{
                     "type":"number",
                     "title":"Sets the number of worker threads per app process",
                     "default": 1
                  },
                  "thread_stack_size":{
                     "type":"number",
                     "title":"Stack size of a worker thread (in bytes)",
                     "default": 512000000
                  }
               },
               "required": [
                  "webapp"
               ]
            },
            {
               "properties":{
                  "lang":{
                     "enum": [
                        "perl"
                     ]
                  },
                  "script":{
                     "type":"string",
                     "title":"PSGI script path",
                     "description": "Name of PSGI script relative to working directory)",
                     "format": "perl-script"
                  },
                  "threads":{
                     "type":"number",
                     "title":"Sets the number of worker threads per app process",
                     "default": 1
                  },
                  "thread_stack_size":{
                     "type":"number",
                     "title":"Stack size of a worker thread (in bytes)",
                     "default": 512000000
                  }
               },
               "required":[
                  "script"
               ]
            },
            {
               "properties":{
                  "lang":{
                     "enum":[
                        "php"
                     ]
                  },
                  "options": {
                     "type": "object",
                     "title": "PHP options settings",
                     "properties": {
                        "file": {
                           "type": "string",
                           "title": "PHP.ini file location",
                           "description": "Absolute pathname of the php.ini file with PHP configuration directives"
                        },
                        "admin": {
                           "type": "object",
                           "title": "Admin options for PHP",
                           "description": "Objects for extra directives. Values in admin are set in PHP_INI_SYSTEM mode, so the app can’t alter them",
                           "additionalProperties": {
                              "type": "string",
                           }
                        },
                        "user": {
                           "type": "object",
                           "title": "User options for PHP",
                           "description": "Objects for extra directives. User values are set in PHP_INI_USER mode and may be updated in runtime",
                           "additionalProperties": {
                              "type": "string",
                           }
                        },
                     }
                  },
                  "targets": {
                     "type": "array",
                     "minItems": 1,
                     "title": "PHP targets",
                     "description": "Define application sections with custom root, script, and index values with max 254 individual entry points for a single PHP application",
                     "items": {
                        "type": "object",
                        "properties": {
                           "reference":{
                              "type":"string",
                              "title":"PHP application reference",
                              "description":"The reference name to be used for the NGINX Unit configuration"
                           },
                           "root":{
                              "type":"string",
                              "title":"Root directory",
                              "description":"Base directory of your PHP app’s file structure. All URI paths are relative to this value"
                           }
                        },
                        "oneOf": [
                           {
                              "title":"Index file name",
                              "properties": {
                                 "index":{
                                    "description":"Filename appended to any URI paths ending with a slash; applies if script is omitted. The default value is index.php",
                                    "type":"string",
                                 }
                              },
                              "required":[
                                 "index"
                              ]
                           },
                           {
                              "title":"PHP script name",   
                              "properties": {
                                 "script":{
                                    "description":"Filename of a root-based PHP script that Unit uses to serve all requests to the application. Omit this if index is used",
                                    "type":"string",
                                 }
                              },
                              "required":[
                                 "script"
                              ]
                           }
                        ]
                     },
                     "required": [
                        "reference", "root"
                     ]
                  }
               }
            },
            {
               "properties":{
                  "lang":{
                     "enum": [
                        "python"
                     ]
                  },
                  "module":{
                     "type":"string",
                     "title":"Application module name",
                     "description": "The module itself is imported just like in Python"
                  },
                  "callable":{
                     "type":"string",
                     "title":"Callable name",
                     "description": "Name of the callable in module that Unit uses to run the app. The default value is application"
                  },
                  "home":{
                     "type":"string",
                     "title":"Home directory",
                     "description": "Path to the app’s virtual environment. Relative to working_directory"
                  },
                  "path":{
                     "type":"string",
                     "title":"Lookup path",
                     "description": "Additional lookup path for Python modules; this string is inserted into sys.path"
                  },
                  "protocol":{
                     "type":"string",
                     "description": "A hint to instruct Unit that the app uses a certain interface; can be asgi or wsgi",
                     "title":"App interface",
                     "enum":[
                        "asgi",
                        "wsgi"
                     ],
                     "default":"wsgi"
                  },
                  "threads":{
                     "type":"number",
                     "description":"Thread count",
                     "title":"Sets the number of worker threads per app process",
                     "default": 1
                  },
                  "thread_stack_size":{
                     "type":"number",
                     "title":"Stack size",
                     "description":"Stack size of a worker thread (in bytes)",
                  }
               },
               "required":[
                  "module"
               ]
            },
            {
               "properties":{
                  "lang":{
                     "enum": [
                        "ruby"
                     ]
                  },
                  "script":{
                     "type":"string",
                     "title":"Rack script name",
                     "description": "Rack script name relative to the working directory, including the .ru extension: script.ru",
                     "format": "ruby-script"
                  },
                  "threads":{
                     "type":"number",
                     "description":"Thread count",
                     "title":"Sets the number of worker threads per app process",
                     "default": 1
                  }
               },
               "required": [
                  "script"
               ]
            }
         ]
      }
   }
}

exports.uiSchema = {
   "lang": {
      "ui:autofocus": true
   }
}

