package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

// Pointer to DB struct
var db *sql.DB

func init_logger(log_path string) {
	logger, err := os.OpenFile(log_path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("[FATAL] Error opening file: %v", err)
	}
	defer logger.Close()
	log.SetOutput(logger)
}

func GetConnection() *sql.DB {
	// Init logger
	init_logger("connection.go")

	// Avoid renew the DB connection in every call
	if db != nil {
		return db
	}

	var err error
	// Database connection
	db, err = sql.Open("sqlite3", "data.sqlite")
	if err != nil {
		log.Panicf(err.Error())
	}
	return db
}
