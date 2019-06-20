package client

import (
	"strings"
	"fmt"
	"os"
	"github.com/stewartad/singolang/utils"
)

type existError struct {
	instance string
}

func (e *existError) Error() string {
	return fmt.Sprintf("%s is not a loaded instance", e.instance)
}

// Execute runs a command inside a container
func (c *Client) Execute(instance string, command string) (string, error) {
	// TODO: check install

	cmd := utils.InitCommand("exec")

	i, e := c.instances[instance]
	if !e {
		return "", &existError{instance} 
	}

	// --nv for graphics card drivers

	// use client's loaded instance by default

	// if instance, use its uri

	// TODO: bind paths
	
	// TODO: run an app

	// TODO: sudo/writable

	splitCommand := strings.Split(command, " ")
	cmd = append(cmd, i.image)
	cmd = append(cmd, splitCommand...)

	stdout, stderr, status, err := utils.RunCommand(cmd, false, false)
	// TODO: use status
	_ = status
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running command: %s\n", strings.Join(cmd, " "))
		return string(stderr.Bytes()), err
	}
	return string(stdout.Bytes()), nil
	
}