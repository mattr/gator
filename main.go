package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/mattr/gator/internal/config"
	"github.com/mattr/gator/internal/database"
	"log"
	"os"
)

type state struct {
	config *config.Config
	db     *database.Queries
}

type command struct {
	name string
	args []string
}

type commands struct {
	available map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.available[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	f := c.available[cmd.name]
	if f == nil {
		return fmt.Errorf("command %q not found", cmd.name)
	}
	return f(s, cmd)
}

func registerCommands(c *commands) {
	c.register("login", handlerLogin)
	c.register("register", handlerRegister)
	c.register("reset", handlerReset)
	c.register("users", handlerUsers)
	c.register("agg", handlerAggregator)
	c.register("feeds", handlerFeeds)
	c.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	c.register("follow", middlewareLoggedIn(handlerFeedFollow))
	c.register("following", middlewareLoggedIn(handlerFeedFollowing))
	c.register("unfollow", middlewareLoggedIn(handlerFeedUnfollow))
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	s := &state{
		config: &cfg,
		db:     database.New(db),
	}

	userArgs := os.Args
	if len(userArgs) < 2 {
		log.Fatal("Insufficient arguments provided")
	}

	c := &commands{available: make(map[string]func(*state, command) error)}
	registerCommands(c)
	cmd := command{name: userArgs[1], args: userArgs[2:]}
	err = c.run(s, cmd)
	if err != nil {
		log.Fatal(err)
	}
}
