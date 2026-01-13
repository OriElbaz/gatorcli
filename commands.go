package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/OriElbaz/gatorcli/internal/config"
	"github.com/OriElbaz/gatorcli/internal/database"
	"github.com/google/uuid"
	"github.com/OriElbaz/gatorcli/rss"
)


func (c *commands) run(s *state, cmd command) error {
	commandName := cmd.name

	if err := c.commands[commandName](s, cmd); err != nil {
		return fmt.Errorf("Error running command: %v\n", err)
	}

	return nil
}


func (c *commands) register(name string, f func(*state, command) error) error {
	c.commands[name] = f

	if _, ok := c.commands[name]; !ok {
		return fmt.Errorf("ERROR command was not added\n")
	}

	fmt.Print("Command was successfully added\n")

	return nil
}

/***** STRUCTS *****/
type state struct {
	db  *database.Queries
	cfg *config.Config
}


type command struct {
	name      string
	arguments []string
}


type commands struct {
	commands map[string]func(*state, command) error
}




/****** COMMANDS ******/
func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) != 1 {
		return fmt.Errorf("Incorrect arguments for command\n")
	}

	username := sql.NullString{
		String: cmd.arguments[0],
		Valid: true,
	}

	if _, err := s.db.GetUser(context.Background(), username); err != nil {
		return fmt.Errorf("Getting user from db: %v\n", err)
	}

	if err := s.cfg.SetUser(username.String); err != nil {
		return fmt.Errorf("Unable to set config username: %v\n", err)
	}

	fmt.Printf("Logged in successfully\n")
	return nil
}


func handlerRegister(s *state, cmd command) error { 
	userName := sql.NullString{String: cmd.arguments[0], Valid:  true}
	params := database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: userName,
	}

	if _, err := s.db.CreateUser(context.Background(), params); err != nil {
		return fmt.Errorf("create user in db: %w", err)
	} 
	fmt.Printf("Adding user to db was sucessful!\n")
	log.Print(params)

	if err := s.cfg.SetUser(params.Name.String); err != nil {
		return fmt.Errorf("set config username: %w", err)
	}

	return nil
}


func reset(s *state, cmd command) error {
	if err := s.db.ClearTableUsers(context.Background()); err != nil {
		return fmt.Errorf("delete all users table: %w", err)
	}

	fmt.Println("users table cleared successfully")
	return nil
}


func users(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("get users names from users table: %w", err)
	}

	for _, user := range users{
		name := user.String

		switch name {
		case s.cfg.CurrentUserName:
			fmt.Printf("* %s (current)\n", name)
		default:
			fmt.Printf("* %s\n", name)
		}
	}

	return nil
}


func agg(s *state, cmd command) error {
	url := "https://www.wagslane.dev/index.xml"
	feed, err := rss.FetchFeed(context.Background(), url)
	if err != nil {
		return fmt.Errorf("fetch feed: %w", err)
	}

	fmt.Println("=== FEED ===")
	fmt.Printf("%s\n", feed.Channel.Title)
	fmt.Printf("%s\n", feed.Channel.Description)
	for _, item := range feed.Channel.Item {
		fmt.Printf("- Title: %s\n", item.Title)
		fmt.Printf("- Description: %s\n", item.Description)
	}

	return nil
}