// Package that can parse and process commands
package command

import (
	"errors"
	"strings"
)

// Represents a redis command name and its args
type Command struct {
	name string
	arg1 string
	arg2 string
}

// Get the Command name
func (c *Command) Name() string {
	return c.name
}

// Get the first Command argument
func (c *Command) Arg1() string {
	return c.arg1
}

// Get the second Command argument
func (c *Command) Arg2() string {
	return c.arg2
}

// Parses a message string into a Command's components.
//
// This can look like:
//
// - "GET key"
//
// - "SET key value"
//
// - "DEL key"
func Parse(message string) (Command, error) {

	fields := strings.Fields(message)

	if len(fields) < 2 {
		return Command{
			name: "",
			arg1: "",
			arg2: "",
		}, errors.New("Invalid command, not enough arguments.")
	}

	if len(fields) == 2 {
		return Command{
			name: fields[0],
			arg1: fields[1],
			arg2: "",
		}, nil
	}

	return Command{
		name: fields[0],
		arg1: fields[1],
		arg2: fields[2],
	}, nil
}
