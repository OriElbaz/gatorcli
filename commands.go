package main

import (
	"fmt"
	"github.com/OriElbaz/gatorcli/internal/config"

)

func (c *commands) run(s *state, cmd command) error {
	commandName := cmd.name

	err := c.commands[commandName](s, cmd)
	if err != nil {
		return fmt.Errorf("Error running command: %v\n", err)
	}

	return nil
}

func (c *commands) register(name string, f func(*state, command) error) error {
	c.commands[name] = f

	// check if added correctly
	_, ok := c.commands[name]
	if !ok {
		return fmt.Errorf("ERROR command was not added\n")
	}

	fmt.Print("Command was successfully added\n")


	return nil
}

type state struct {
	config *config.Config
}

type command struct {
	name string
	arguments []string
}

type commands struct {
	commands map[string]func(*state, command) error
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) != 1 {
		return fmt.Errorf("Incorrect arguments for command\n")
	}

	username := cmd.arguments[0]

	err := s.config.SetUser(username)
	if err != nil {
		return fmt.Errorf("Unable to set config username: %v\n", err)
	}

	fmt.Printf("Config username has been set successfully\n")
	return nil
}
