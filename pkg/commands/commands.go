package commands

import (
	"context"
	"database/sql"
	"fmt"
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


/****** MIDDLEWARE ******/
func MiddlewareLoggedIn(handler func(s *State, cmd Command, user database.User) error) func(*State, Command) error {
	return func (s *State, cmd Command) error {
		userNullString := sql.NullString{
			String: s.Cfg.CurrentUserName,
			Valid: true,
		}

		user, err := s.Db.GetUser(context.Background(), userNullString);
		if err != nil {
			return fmt.Errorf("get user: %w", err)
		}

		return handler(s, cmd, user)
	}
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


func AddFeed(s *State, cmd Command, user database.User) error {

	feedName := cmd.Arguments[0]
	feedURL := cmd.Arguments[1]

	feedURLStruct := sql.NullString{
		String: feedURL,
		Valid: true,
	}

	feed := database.CreateFeedParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: feedName,
		Url: feedURLStruct,
		UserID: user.ID,
	}

	if _, err := s.Db.CreateFeed(context.Background(), feed); err != nil {
		return fmt.Errorf("create feed: %w", err)
	}

	fmt.Printf("Feed added successfully\n")

	if _, err := createFeedFollowHelper(s, user.ID, feed.ID); err != nil {
		return fmt.Errorf("create feed follow helper: %w", err)
	}

	return nil
}


func Feeds(s *State, cmd Command) error {
	feeds, err := s.Db.ListFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("list feeds: %w", err)
	}

	for _, feed := range feeds {

		fmt.Printf("== %s ==\n", feed.Name)
		fmt.Printf("- url: %s\n", feed.Url.String)
		fmt.Printf("- user: %s\n", feed.UserName.String)
	}

	return nil
}


func Follow(s *State, cmd Command, user database.User) error {
	urlToAdd := sql.NullString{
		String: cmd.Arguments[0],
		Valid: true,
	}

	feed, err := s.Db.GetFeed(context.Background(), urlToAdd)
	if err != nil {
		return fmt.Errorf("get feed: %w", err)
	}

	if _, err = createFeedFollowHelper(s, user.ID, feed.ID); err != nil {
		return fmt.Errorf("create feed follow helper: %w", err)
	}

	fmt.Printf("Feed follow created successfully!")
	return nil
}


func Following(s *State, cmd Command, user database.User) error {

	feedFollows, err := s.Db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("get feed follows: %w", err)
	}

	for _, feed := range feedFollows {
		fmt.Printf("- %s\n", feed.FeedName)
	}

	return nil
}


/** HELPER FUNCTIONS **/
func createFeedFollowHelper(s *State, userId uuid.UUID, feedId uuid.UUID) (database.CreateFeedFollowRow, error) {
	params := database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID: userId,
		FeedID: feedId,
	}

	feedFollow, err := s.Db.CreateFeedFollow(context.Background(), params)
	if err != nil {
		return database.CreateFeedFollowRow{}, fmt.Errorf("create feed follow: %w", err)
	}

	return feedFollow, nil
}