package main

import (
    "errors"
    "fmt"
    "context"
    "time"
    "github.com/google/uuid"
    "github.com/feagy/gator/internal/database"
    "github.com/feagy/gator/internal/rssfeed"
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
    _, err := s.db.GetUser(context.Background(), cmd.arguments[0])
    if err != nil {
        return err
    }
    err = s.cfg.SetUser(cmd.arguments[0])
    if err != nil {
        return err
    }

    fmt.Println("you have logged in")
    return nil
}

func handlerRegister(s *state, cmd command) error {
    if len(cmd.arguments) < 1 {
        return errors.New("the login command expects a single argument, the username.")
    }
    _, err := s.db.CreateUser(context.Background(), database.CreateUserParams{ID:uuid.New(), CreatedAt:time.Now(), UpdatedAt:time.Now(), Name:cmd.arguments[0],})
    if err != nil {
        return err
    }

    err = s.cfg.SetUser(cmd.arguments[0])
    if err != nil {
        return err
    }

    fmt.Println("new user created and logged in")
    return nil
}

func handlerReset(s *state, cmd command) error {
    err := s.db.DeleteUsers(context.Background())
    if err != nil {
        return err
    }
    fmt.Println("all users deleted")
    return nil
}

func handlerUsers(s *state, cmd command) error {
    users, err := s.db.GetUsers(context.Background())
    if err != nil {
        return err
    }
    
    if len(users) == 0 {
        fmt.Println("There is no user display")
        return nil
    }
    loggedUser := s.cfg.CurrentUserName

    for _, user := range users {
        if user.Name == loggedUser {
            fmt.Printf("* %v (current)\n", user.Name)   
            continue     
        } 
        fmt.Printf("* %v\n", user.Name)
    }
    return nil
}

func handlerAgg(s *state, cmd command) error {
    rssFeed, err := rssfeed.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
    if err != nil {
        return err
    }
    fmt.Println(rssFeed)
    return nil
}

func handlerAddFeed(s *state, cmd command) error {
    if len(cmd.arguments) < 2 {
        return errors.New("the addfeed command expects 2 arguments, the name of the feed and the url.")
    }
    
    user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
    if err != nil {
        return err
    }

    feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{ID:uuid.New(), CreatedAt:time.Now(), UpdatedAt:time.Now(), Name:cmd.arguments[0], Url:cmd.arguments[1], UserID:user.ID})
    if err != nil {
        return err
    }
    
    fmt.Println(feed)
    return nil
}
