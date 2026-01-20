package main

import (
	"fmt"
	"rss_aggregator_mod/internal/config"
	"errors"
	"os"
)

type state struct {
	cfg	*config.Config
}

type command struct {
	name 		string
	args		[]string
}

type commands struct {
	handlers	map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	handler, ok := c.handlers[cmd.name]
	if !ok {
		return errors.New("no such command exists")
	}
	return handler(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.handlers[name] = f
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		fmt.Println("username is required")
		os.Exit(1)
	}

	username := cmd.args[0]

	err := s.cfg.SetUser(username)	
	if err != nil {
		return err
	}

	fmt.Println("user has been set")
	return nil
}
	


func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	programState := &state{cfg: &cfg}
	cmds := commands{
		handlers: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handlerLogin)
	args := os.Args
	if len(args) < 2 {
		fmt.Println("Not enough arguments")
		os.Exit(1)
	}
	commandName := args[1]
	commandSlice := args[2:]
	cmd := command{
		name: commandName, args: commandSlice,
	}
	err = cmds.run(programState, cmd)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
}