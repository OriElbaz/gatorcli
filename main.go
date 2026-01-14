package main

import (
	"database/sql"
	"fmt"
	"os"
	"github.com/OriElbaz/gatorcli/internal/config"
	"github.com/OriElbaz/gatorcli/internal/database"
	_ "github.com/lib/pq"
	"github.com/OriElbaz/gatorcli/pkg/commands"
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

	configState := commands.State{
		Db:  dbQueries,
		Cfg: &configStruct,
	}

	commandMap := map[string]func(*commands.State, commands.Command) error{
		"login": commands.HandlerLogin,
		"register": commands.HandlerRegister,
		"reset": commands.Reset,
		"users": commands.Users,
		"agg": commands.Agg,
		"addfeed": commands.AddFeed,
		"feeds": commands.Feeds,
		"follow": commands.Follow,
		"following": commands.Following,
	}

	commandsStruct := commands.Commands{
		Commands: commandMap,
	}

	commandLineInputs := os.Args
	if len(commandLineInputs) < 2 {
		fmt.Println("incorrect number of arguments")
		os.Exit(1)
	}

	commandName := commandLineInputs[1]
	commandArgs := commandLineInputs[2:]
	commandToRun := commands.Command{
		Name:      commandName,
		Arguments: commandArgs,
	}

	if err = commandsStruct.Run(&configState, commandToRun); err != nil {
		fmt.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}
}
