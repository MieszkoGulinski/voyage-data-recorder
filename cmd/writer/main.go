package main

import (
	"datalogger/database"
	"fmt"
)

func main() {
	// TODO implement writer
	db := database.CreateDatabaseWriterConnection()
	db.Exec("PRAGMA journal_mode = WAL;")
	db.Exec("PRAGMA synchronous = NORMAL;")

	fmt.Println("Writer active")
}