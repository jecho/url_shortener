package main

import (
	"database/sql"
)

// struct object
type Env struct {
	config *Config
	db *sql.DB
}
