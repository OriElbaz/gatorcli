package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/OriElbaz/gatorcli/internal/config"
	"github.com/OriElbaz/gatorcli/internal/database"
	_ "github.com/lib/pq"
)

func main() {

	// database stuff //
	dbURL := "postgres://orielbaz:@localhost:5432/gator?sslmode=disable"
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Printf("Error connecting to db: %v", err)
	}

	dbQueries := database.New(db)

	// read in JSON file into Config struct //
	configStruct, err := config.Read()
	if err != nil {
		fmt.Printf("ERROR with reading gatorconfig.json: %v\n", err)
		os.Exit(1)
	}

	// turn into state struct //
	configState := state{
		db:  dbQueries,
		cfg: &configStruct,
	}

	// initialise commands.go structs //
	commandMap := map[string]func(*state, command) error{
		"login": handlerLogin,
		"register": handlerRegister,
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
		name:      commandName,
		arguments: commandArgs,
	}

	err = commandsStruct.run(&configState, commandToRun)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}

}
