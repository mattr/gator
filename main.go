package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
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

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("login handler expects a single argument (username)")
	}
	user, err := s.db.GetUserByName(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}
	err = s.config.SetUser(user.Name)
	if err != nil {
		return err
	}
	fmt.Println("Logged in as " + s.config.CurrentUserName)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 || len(cmd.args[0]) == 0 {
		return errors.New("register handler expects a single argument (username)")
	}
	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{ID: uuid.New(), Name: cmd.args[0]})
	if err != nil {
		return err
	}
	err = s.config.SetUser(user.Name)
	if err != nil {
		return err
	}
	fmt.Println("Created user ", user.Name)
	fmt.Printf("%v\n", user)
	return nil
}

func handlerReset(s *state, cmd command) error {
	return s.db.DeleteAllUsers(context.Background())
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

	c := commands{available: make(map[string]func(*state, command) error, 3)}
	c.register("login", handlerLogin)
	c.register("register", handlerRegister)
	c.register("reset", handlerReset)

	userArgs := os.Args
	if len(userArgs) < 2 {
		log.Fatal("Insufficient arguments provided")
	}
	cmd := command{name: userArgs[1], args: userArgs[2:]}
	err = c.run(s, cmd)
	if err != nil {
		log.Fatal(err)
	}
}
