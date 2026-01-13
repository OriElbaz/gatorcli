package main

import (
	"database/sql"
	"fmt"
	"os"
	"github.com/OriElbaz/gatorcli/internal/config"
	"github.com/OriElbaz/gatorcli/internal/database"
	_ "github.com/lib/pq"
)


const dbURL = "postgres://orielbaz:@localhost:5432/gator?sslmode=disable"


func main() {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Printf("Error connecting to db: %v", err)
	}

	dbQueries := database.New(db)

	configStruct, err := config.Read()
	if err != nil {
		fmt.Printf("ERROR with reading gatorconfig.json: %v\n", err)
		os.Exit(1)
	}

	configState := state{
		db:  dbQueries,
		cfg: &configStruct,
	}

	commandMap := map[string]func(*state, command) error{
		"login": handlerLogin,
		"register": handlerRegister,
		"reset": reset,
	}

	commandsStruct := commands{
		commands: commandMap,
	}

	commandLineInputs := os.Args
	if len(commandLineInputs) < 2 {
		fmt.Println("incorrect number of arguments")
		os.Exit(1)
	}

	commandName := commandLineInputs[1]
	commandArgs := commandLineInputs[2:]
	commandToRun := command{
		name:      commandName,
		arguments: commandArgs,
	}

	if err = commandsStruct.run(&configState, commandToRun); err != nil {
		fmt.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}
}
