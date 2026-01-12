package main

import (
	"fmt"
	"os"

	"github.com/OriElbaz/gatorcli/internal/config"
)

func main() {
	// read in JSON file into Config struct //
	configStruct, err := config.Read()
	if err != nil {
		fmt.Printf("ERROR with reading gatorconfig.json: %v\n", err)
		os.Exit(1)
	}

	// turn into state struct //
	configState := state{
		config: &configStruct,
	}

	// initialise commands.go structs //
	commandMap := map[string]func(*state, command) error{
		"login" : handlerLogin,
	}

	commandsStruct := commands{
		commands: commandMap,
	}

	// take command line arguments (inputs) and run command //
	commandLineInputs := os.Args
	if len(commandLineInputs) < 2 {
		fmt.Println("ERROR: incorrect number of arguments")
		os.Exit(1)
	} 

	commandName := commandLineInputs[1]
	commandArgs := commandLineInputs[2:]
	
	commandToRun := command{
		name: commandName,
		arguments: commandArgs,
	}

	err = commandsStruct.run(&configState, commandToRun)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}

}
