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
