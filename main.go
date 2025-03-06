package main

import (
	"errors"
	"fmt"
	"github.com/mattr/gator/internal/config"
	"log"
	"os"
)

type state struct {
	config *config.Config
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
	err := s.config.SetUser(cmd.args[0])
	if err != nil {
		return err
	}
	fmt.Println("Logged in as " + s.config.CurrentUserName)
	return nil
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	s := &state{
		config: &cfg,
	}

	c := commands{available: make(map[string]func(*state, command) error, 3)}
	c.register("login", handlerLogin)

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
