package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/westleaf/gator/internal/database"
)

type command struct {
	name string
	args []string
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("arguments cannot be empty")
	}

	if _, err := s.db.GetUser(context.Background(), cmd.args[0]); err != nil {
		return fmt.Errorf("user does not exist: %w", err)
	}

	err := s.config.SetUser(cmd.args[0])
	if err != nil {
		return err
	}

	fmt.Printf("User has ben set to: %s\n", cmd.args[0])

	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("not enough arguments passed")
	}

	ctx := context.Background()

	user, err := s.db.CreateUser(ctx, database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
	})
	if err != nil {
		return err
	}
	err = s.config.SetUser(cmd.args[0])
	if err != nil {
		return err
	}

	fmt.Printf("user data: %s", user)
	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.ResetUsers(context.Background())
	if err != nil {
		return err
	}

	fmt.Printf("cleared users from db\n")
	return nil
}

func handlerGetUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}

	for _, user := range users {
		if user.Name == s.config.CurrentUserName {
			fmt.Printf("%s (current)\n", user.Name)
			continue
		}
		fmt.Printf("%s\n", user.Name)
	}
	return nil
}

func handlerGetFeed(s *state, cmd command) error {
	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}

	fmt.Printf("%+v", feed)

	return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID:    user.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return err
	}

	fmt.Printf("Feed: %+v", feed)

	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return err
	}

	return nil
}

func handlerListFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return nil
	}

	for _, feed := range feeds {
		user, err := s.db.GetUserById(context.Background(), feed.UserID)
		if err != nil {
			return err
		}
		fmt.Printf("Name: %s, Url: %s, Added by: %s\n", feed.Name, feed.Url, user.Name)
	}
	return nil
}

func handlerFollowFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 || len(cmd.args) > 1 {
		return fmt.Errorf("unacceptable number of arguments")
	}

	feed, err := s.db.GetFeedByUrl(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}

	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return err
	}

	fmt.Printf("User: %s, Feed: %s", user.Name, feed.Name)

	return nil
}

func handlerListFollowFeedForUser(s *state, cmd command, user database.User) error {
	if len(cmd.args) > 0 {
		return fmt.Errorf("command takes no arguments")
	}

	feeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}

	if len(feeds) == 0 {
		return nil //fmt.Errorf("no feeds followed for user: %s", user.Name)
	}

	fmt.Printf("List of user %s followed feeds:\n", user.Name)
	for _, feed := range feeds {
		fmt.Println(feed.FeedName)
	}
	return nil
}

func handlerDeleteFeedFollowForUser(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("not enough arguments in command")
	}

	feed, err := s.db.GetFeedByUrl(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}

	s.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})

	fmt.Printf("user: %s unlfollowed feed: %s\n", user.Name, feed.Name)

	return nil
}

type commands struct {
	commandNames map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {

	if _, ok := c.commandNames[cmd.name]; ok {
		err := c.commandNames[cmd.name](s, cmd)
		if err != nil {
			return err
		}
	} else {
		log.Fatal("command not found: " + cmd.name)
	}
	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.commandNames[name] = f
}
