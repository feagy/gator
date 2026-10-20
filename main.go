package main

import (
    "github.com/feagy/gator/internal/config"
    "fmt"
)

func main() {
    cfg, err := config.Read()
    if err != nil{
        fmt.Println(err)    
    }
    err = cfg.SetUser("ferhat")
    if err != nil{
        fmt.Println(err)    
    }
    cfg, err = config.Read()
    if err != nil{
        fmt.Println(err)    
    }

    fmt.Println(cfg)
}

