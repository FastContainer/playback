package main

import (
	"fmt"
	"os/exec"
	"strings"
)

// Cmder interface
type Cmder interface {
	Do(name string, arg ...string) ([]byte, error)
}

// Cmd struct
type Cmd struct{}

// Do execute command
func (c Cmd) Do(name string, arg ...string) ([]byte, error) {
	return exec.Command(name, arg...).Output()
}

// MockCmd struct
type MockCmd struct {
	Out string
	Err string
}

// Do execute command
func (c MockCmd) Do(name string, arg ...string) ([]byte, error) {
	var err error
	if c.Err != "" {
		err = fmt.Errorf(c.Err)
	}
	fmt.Printf("%s %s\n", name, strings.Join(arg, " "))
	return []byte(c.Out), err
}
