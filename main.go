package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/westleaf/gator/internal/config"
	"github.com/westleaf/gator/internal/database"

	_ "github.com/lib/pq"
)

type state struct {
	db     *database.Queries
	config *config.Config
}

func main() {

	if len(os.Args) < 2 {
		log.Fatal("not enough arguments")
	}

	conf, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("postgres", conf.DbURL)
	if err != nil {
		log.Fatal(err)
	}

	dbQueries := database.New(db)

	st := state{
		dbQueries,
		&conf,
	}

	cmd := commands{
		commandNames: make(map[string]func(*state, command) error),
	}

	cmd.register("login", handlerLogin)
	cmd.register("register", handlerRegister)
	cmd.register("reset", handlerReset)
	cmd.register("users", handlerGetUsers)
	cmd.register("agg", handlerGetFeed)
	cmd.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmd.register("feeds", handlerListFeeds)
	cmd.register("follow", middlewareLoggedIn(handlerFollowFeed))
	cmd.register("following", middlewareLoggedIn(handlerListFollowFeedForUser))
	cmd.register("unfollow", middlewareLoggedIn(handlerDeleteFeedFollowForUser))

	err = cmd.run(&st, command{
		name: os.Args[1],
		args: os.Args[2:],
	})
	if err != nil {
		log.Fatal(err)
	}
}
