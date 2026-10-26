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
        fmt.Println("there is no user display")
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
        return errors.New("the addfeed command expects 2 arguments, the name and the url of the feed .")
    }
    
    user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
    if err != nil {
        return err
    }

    feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{ID:uuid.New(), CreatedAt:time.Now(), UpdatedAt:time.Now(), Name:cmd.arguments[0], Url:cmd.arguments[1], UserID:user.ID})
    if err != nil {
        return err
    }
    
    _, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{ID:uuid.New(), CreatedAt:time.Now(), UpdatedAt:time.Now(), UserID:user.ID, FeedID:feed.ID})
    if err != nil {
        return err
    }

    return nil
}

func handlerFeeds(s *state, cmd command) error {
    feeds, err := s.db.GetFeeds(context.Background())
    if err != nil {
        return err
    }
    
    if len(feeds) == 0 {
        fmt.Println("there is no feed display")
        return nil
    }

    for _, feed := range feeds {
        user, err := s.db.GetUserById(context.Background(), feed.UserID)
        if err != nil {
            return err
        }
        fmt.Printf("* Name: %v\n  URL: %v\n  User who created the feed: %v\n", feed.Name, feed.Url, user.Name)
    }
    return nil
}

func handlerFollow(s *state, cmd command) error {
    if len(cmd.arguments) < 1 {
        return errors.New("the addfeed command expects 1 arguments, the url of the feed.")
    }

    feed, err := s.db.GetFeedByUrl(context.Background(), cmd.arguments[0])
    if err != nil {
        return err
    }

    user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
    if err != nil {
        return err
    }

    feedFollowed, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{ID:uuid.New(), CreatedAt:time.Now(), UpdatedAt:time.Now(), UserID:user.ID, FeedID:feed.ID})
    if err != nil {
        return err
    }

    fmt.Println(feedFollowed)
    return nil
}

func handlerFollowing(s *state, cmd command) error {
    user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
    if err != nil {
        return err
    }

    feedsFollowed, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
    if err != nil {
        return err
    }

    if len(feedsFollowed) == 0 {
        fmt.Println("you are not following any feeds")
        return nil
    }
    
    for _, feedFollowed := range feedsFollowed {
        fmt.Printf("* %v\n", feedFollowed.FeedName)
    }
    return nil
}
