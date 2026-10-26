package main

import (
    _ "github.com/lib/pq"
    "github.com/feagy/gator/internal/config"
    "github.com/feagy/gator/internal/database"
    "fmt"
    "os"
    "database/sql"
)

func main() {
    cfg, err := config.Read()
    if err != nil{
        fmt.Println(err) 
        os.Exit(1)   
    }

    db, err := sql.Open("postgres", cfg.DbUrl)
    if err != nil{
        fmt.Println(err) 
        os.Exit(1)   
    }
    dbQueries := database.New(db)

    s := state{
        cfg: &cfg,
        db: dbQueries,
    }

    c := commands{
        commands: make(map[string]func(*state, command) error),
    }
    
    c.register("login", handlerLogin)
    c.register("register", handlerRegister)
    c.register("reset", handlerReset)
    c.register("users", handlerUsers)
    c.register("agg", handlerAgg)
    c.register("addfeed", handlerAddFeed)
    c.register("feeds", handlerFeeds)
    c.register("follow", handlerFollow)
    c.register("following", handlerFollowing)

    args := os.Args
    if len(args) < 2 {
        fmt.Println("less than 2 arguments the command should be like this: gator <command name> <arguments>") 
        os.Exit(1)   
    }

    cmd := command{
        name: args[1],
        arguments: args[2:],
    }
    err = c.run(&s, cmd)
    if err != nil{
        fmt.Println(err) 
        os.Exit(1)   
    }
    
} 

