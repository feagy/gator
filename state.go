package main

import (
    "github.com/feagy/gator/internal/config"
    "github.com/feagy/gator/internal/database"
)

type state struct {
    db  *database.Queries
    cfg *config.Config
}
