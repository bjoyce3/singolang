# Singolang

[![Go Report Card](https://goreportcard.com/badge/github.com/stewartad/singolang)](https://goreportcard.com/report/github.com/stewartad/singolang)

Singolang is a library to interact with Singularity containers in Go. It is modeled from Spython. Designed for use with Singularity 3+

## Currently Supported Features

* Pulling images from Dockerhub or Singularity Hub
* Starting and Stopping instances of built images
* Executing commands in running instances
* Copying files inside a running instance to the host system

## Usage

### Create a Client

To start using Singolang, create a new client

```go
client, teardown := singolang.NewClient()
defer teardown()
```

`NewClient()` returns both a Client struct and a teardown function. The teardown function stops all instances and is intended to be deferred

### Pull an image

Pulling an image requires a struct to be filled out detailing the options with which to perform the pull.

Name is the filename to save the image as.

Pullfolder is the folder to place the final image

Force, if true, will overwrite any existing files of the same name

```go
opts := &singolang.PullOptions{
    Name: "",
    Pullfolder: filepath.Join("/tmp", "pull"),
    Force: false,
}

imgPath, err := client.Pull("docker://godlovedc/lolcow", pullOpts)

if err != nil {
    fmt.Println(err)
} else {
    fmt.Println(imgPath)
}
```

### Create an Instance

You can define environment variables as well as modify the PATH variable of a container before creating an instance by filling out an `EnvOption` struct and passing it to the client's `NewInstance()` function

```go
instanceEnv := singolang.EnvOptions {
    EnvVars: map[string]string {
        "hello": "world",
    },
    PrependPath: []string{"/usr/local/go/bin"},
    AppendPath: []string{"/usr/local/games"},
    ReplacePath: "",
}

instance, err := client.NewInstance("lolcow_latest.sif", "lolcow3", &instanceEnv)
if err != nil {
    fmt.Println(err)
}
```

`NewInstance()` creates and returns a new Instance struct and adds it to the client's list of instances.

Start the instance

```go
instance.Start()
```

Stop the instance

```go
instance.Stop()
```

### Execute a Command

Executing a command has its own configuration struct. Pwd will execute the command in the given directory. Quiet, if true, will supress output to stdout and stderr. Cleanenv will ensure the host's environment variables are not copied over to the container while while executing the command. You can specify EnvOptions like above to change the containers environment before the command runs.

```go
execOpts := singolangExecOptions{
    Pwd:   "",
    Quiet: true,
    Cleanenv: true,
    Env: DefaultEnvOptions(),
}

stdout, stderr, code, err := instance.Execute([]string{"which", "fortune"}, %execOpts)
```

### Copy a file

You can copy a file or folder from inside the container into a .tar archive, which is placed in your OS temp directory

```go
path, read, err := instance.CopyTarball(targetPath)
```

`CopyTarball()` returns the path to the archive, a Tar reader for the archive, and an error, if any

## Links
- Link to https://github.com/stewartad/container-structure-test and explanation of how it's important for this
- Link to containerized examples https://github.com/stewartad/singularity-container-test

## License

   Copyright 2019 Austen Stewart

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
