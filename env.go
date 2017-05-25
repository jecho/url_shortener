package main

import (
	"database/sql"
)

// struct object
type Env struct {
	config *Config
	db *sql.DB
	//logger *log.Logger
	//templates *template.Template
}
