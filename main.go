package main

import (
    "github.com/feagy/gator/internal/config"
    "fmt"
    "os"
)

func main() {
    cfg, err := config.Read()
    if err != nil{
        fmt.Println(err) 
        os.Exit(1)   
    }

    s := state{cfg: &cfg}

    c := commands{
        commands: make(map[string]func(*state, command) error),
    }
    
    c.register("login", handlerLogin)

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

