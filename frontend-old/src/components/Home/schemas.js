exports.baseSchema = {
    "title":"Golang Unit configuration options",
    "type":"object",
    "properties":{
       "app_type":{
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
       "limits":{
          "title":"Request Limits",
          "type":"object",
          "properties":{
             "timeout":{
                "description":"Request timeout in seconds",
                "type":"number"
             },
             "requests":{
                "description":"Maximum number of requests allowed to serve",
                "type":"number"
             }
          }
       },
       "processes":{
          "title":"Process Management",
          "type":"object",
          "properties":{
             "max":{
                "description":"Maximum number of application processes",
                "type":"number"
             },
             "spare":{
                "description":"Minimum number of idle processes that Unit tries to reserve for an app",
                "type":"number"
             },
             "idle_timeout":{
                "description":"Time in seconds that Unit waits before terminating an idle process which exceeds spare",
                "type":"number"
             }
          }
       },
       "working_directory":{
          "type":"string",
          "title":"Working Directory",
          "description":"The applications working directory"
       },
       "user":{
          "type":"string",
          "title":"Username",
          "description": "Username that runs the process"
       },
       "group":{
          "type":"string",
          "title":"Group",
          "description": "Group name that runs the process"
       },
       "environment":{
          "type":"array",
          "title":"Environmental variables",
          "description":"Environment variables to be passed to the application",
          "items":{
             "type":"object",
             "properties":{
                "title":{
                   "type":"string",
                   "title":"Variable name"
                },
                "details":{
                   "type":"string",
                   "title":"Variable value"
                }
             }
          }
       },
      "required": [
         "app_type"
      ]
    }
 }

exports.testSchema = {
   "title":"NGINX Unit configuration options",
   "type":"object",
   "properties":{
      "app_type":{
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
         "description":"The applications working directory"
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
      "app_type", "working_directory"
   ],
   "dependencies": {
      "app_type":{
         "oneOf": [
            {
               "properties":{
                  "app_type":{
                     "enum": [
                        "go", "nodejs"
                     ]
                  },
                  "executable":{
                     "type":"string",
                     "title":"Exectuable name",
                     "description":"Pathname of the application, absolute or relative to working directory"
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
                  "app_type":{
                     "enum": [
                        "java"
                     ]
                  },
                  "webapp":{
                     "type":"string",
                     "title": "Application name",
                     "description":"Pathname and name of the application’s packaged or unpackaged .war file: e.g. /www/helloworld/helloworld.war"
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
                        "type":"object",
                        "properties":{
                           "option":{
                              "type":"string",
                              "description":"JVM option"
                           },
                           "value":{
                              "type":"string",
                              "description":"Option value"
                           }
                        }
                     }
                  },
                  "threads":{
                     "type":"number",
                     "title":"Sets the number of worker threads per app process",
                     "default": 1
                  },
                  "thread_stack_size":{
                     "type":"number",
                     "title":"Stack size of a worker thread (in bytes)"
                  }
               },
               "required": [
                  "webapp"
               ]
            },
            {
               "properties":{
                  "app_type":{
                     "enum": [
                        "perl"
                     ]
                  },
                  "script":{
                     "type":"string",
                     "title":"PSGI script path",
                     "description": "Absolute pathname and name of PSGI script (or relative to working directory)"
                  },
                  "threads":{
                     "type":"number",
                     "title":"Sets the number of worker threads per app process",
                     "default": 1
                  },
                  "thread_stack_size":{
                     "type":"number",
                     "title":"Stack size of a worker thread (in bytes)"
                  }
               },
               "required":[
                  "script"
               ]
            },
            {
               "properties":{
                  "app_type":{
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
                           "description": "Pathname of the php.ini file with PHP configuration directives"
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
                           },
                           "choice": {
                              "type": "object",
                              "title": "Select the PHP startup method",
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
                           }
                        },
                        "required": [
                           "reference", "root"
                        ]
                     }
                  }
               }
            },
            {
               "properties":{
                  "app_type":{
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
                     "title":"Module name",
                     "description": "Name of the callable in module that Unit uses to run the app. The default value is application"
                  },
                  "home":{
                     "type":"string",
                     "title":"Home directory",
                     "description": "Path to the app’s virtual environment. Absolute or relative to working_directory"
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
                  "app_type":{
                     "enum": [
                        "ruby"
                     ]
                  },
                  "script":{
                     "type":"string",
                     "title":"Rack script name",
                     "description": "Rack script pathname, including the .ru extension: /www/rubyapp/script.ru"
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
   "app_type": {
      "ui:autofocus": true
   }
}

