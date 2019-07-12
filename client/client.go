package client

import (
	"fmt"
	"log"
	"strings"
)

type instanceError struct {
	name   string
	action string
}

func (e *instanceError) Error() string {
	return fmt.Sprintf("Error %sing Instance: %s. Must be stopped manually in Singularity", e.name, e.action)
}

// Client is a struct to hold information about the current client
type Client struct {
	simage    string // this will be assigned by the load() function
	instances map[string]*instance
	Sudo      bool // either everything or nothing you do is sudo
	Cleanenv  bool
}

// NewClient creates and returns a new client as well as a teardown function.
// Assign this teardown function and defer it to exit cleanly
func NewClient() (*Client, func(c *Client)) {
	return &Client{
			simage:    "",
			instances: make(map[string]*instance),
			Sudo:      false,
			Cleanenv:  true,
		},
		func(c *Client) { c.teardown() }
}

// Version returns the version of the system's Singularity installation
func (c *Client) Version() string {
	return GetSingularityVersion()
}

func (c *Client) String() string {
	baseClient := "[singularity-golang]"
	if c.simage != "" {
		baseClient = fmt.Sprintf("%s[%s]", baseClient, c.simage)
	}
	return baseClient
}

// NewInstance creates a new instance and adds it to the client, if it is able to be started
func (c *Client) NewInstance(image string, name string) error {
	i := getInstance(image, name)
	err := i.start(c.Sudo)
	if err != nil {
		return err
	}
	c.instances[name] = i
	return nil
}

// StopInstance stops an instance previously created in the client
// TODO: Define custom errors
func (c *Client) StopInstance(name string) error {
	err := c.instances[name].stop(c.Sudo)
	return err
}

// StopAllInstances stops all instances created in the client
func (c *Client) StopAllInstances() error {
	var err error
	for k := range c.instances {
		err = c.StopInstance(k)
	}
	return err
}

func (c *Client) GetEnv(instance string) map[string]string {
	return c.instances[instance].getEnv()
}

func (c *Client) GetEnvVar(instance string, varname string) (string, string) {
	return c.instances[instance].getEnvVar(varname)
}

// ListInstances prints all client-created instances to screen
func (c *Client) ListInstances() {
	fmt.Println("CLIENT LOADED INSTANCES")
	fmt.Println("-----------------")
	if len(c.instances) < 1 {
		fmt.Println("No Loaded Instances\n-----------------")
		return
	}
	for k, v := range c.instances {
		fmt.Printf("%s: %s\n", k, v)
	}
	fmt.Println("-----------------")
}

// ListAllInstances lists all currently running Singularity instances.
// It is equivalent to running `singularity instance list`
func ListAllInstances() {
	cmd := initCommand("instance", "list")

	output, stderr, status, err := runCommand(cmd, defaultRunCommandOptions())
	// TODO: do something with these values
	_, _, _ = output, status, stderr
	if err != nil {
		log.Printf("Error running command: %s\n", strings.Join(cmd, " "))
	}
}

func (c *Client) teardown() {
	fmt.Println("Performing Cleanup")
	c.StopAllInstances()
	ListAllInstances()
}
