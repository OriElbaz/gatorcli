package commands

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
	"github.com/OriElbaz/gatorcli/internal/config"
	"github.com/OriElbaz/gatorcli/internal/database"
	"github.com/google/uuid"
	"github.com/OriElbaz/gatorcli/pkg/rss"
)


func (c *Commands) Run(s *State, cmd Command) error {
	commandName := cmd.Name

	if err := c.Commands[commandName](s, cmd); err != nil {
		return fmt.Errorf("Error running command: %v\n", err)
	}

	return nil
}


func (c *Commands) register(name string, f func(*State, Command) error) error {
	c.Commands[name] = f

	if _, ok := c.Commands[name]; !ok {
		return fmt.Errorf("ERROR command was not added\n")
	}

	fmt.Print("Command was successfully added\n")

	return nil
}

/***** STRUCTS *****/
type State struct {
	Db  *database.Queries
	Cfg *config.Config
}


type Command struct {
	Name      string
	Arguments []string
}


type Commands struct {
	Commands map[string]func(*State, Command) error
}




/****** COMMANDS ******/
func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Arguments) != 1 {
		return fmt.Errorf("Incorrect arguments for command\n")
	}

	username := sql.NullString{
		String: cmd.Arguments[0],
		Valid: true,
	}

	if _, err := s.Db.GetUser(context.Background(), username); err != nil {
		return fmt.Errorf("Getting user from db: %v\n", err)
	}

	if err := s.Cfg.SetUser(username.String); err != nil {
		return fmt.Errorf("Unable to set config username: %v\n", err)
	}

	fmt.Printf("Logged in successfully\n")
	return nil
}


func HandlerRegister(s *State, cmd Command) error { 
	userName := sql.NullString{String: cmd.Arguments[0], Valid:  true}
	params := database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: userName,
	}

	if _, err := s.Db.CreateUser(context.Background(), params); err != nil {
		return fmt.Errorf("create user in db: %w", err)
	} 
	fmt.Printf("Adding user to db was sucessful!\n")
	log.Print(params)

	if err := s.Cfg.SetUser(params.Name.String); err != nil {
		return fmt.Errorf("set config username: %w", err)
	}

	return nil
}


func Reset(s *State, cmd Command) error {
	if err := s.Db.ClearTableUsers(context.Background()); err != nil {
		return fmt.Errorf("delete all users table: %w", err)
	}

	fmt.Println("users table cleared successfully")
	return nil
}


func Users(s *State, cmd Command) error {
	users, err := s.Db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("get users names from users table: %w", err)
	}

	for _, user := range users{
		name := user.String

		switch name {
		case s.Cfg.CurrentUserName:
			fmt.Printf("* %s (current)\n", name)
		default:
			fmt.Printf("* %s\n", name)
		}
	}

	return nil
}


func Agg(s *State, cmd Command) error {
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