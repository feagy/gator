package main

import (
    "errors"
    "fmt"
)

type command struct {
    name string
    arguments []string
}

type commands struct {
    commands map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
    err := c.commands[cmd.name](s, cmd)
    if err != nil {
        return err
    }
    return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
    c.commands[name] = f
}

func handlerLogin(s *state, cmd command) error {
    if len(cmd.arguments) == 0 {
        return errors.New("the login command expects a single argument, the username.")
    }
    
    s.cfg.SetUser(cmd.arguments[0])
    fmt.Println("you have logged in")
    return nil
}
