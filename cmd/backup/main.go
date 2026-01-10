package main

import (
	"log"
)

func main() {
	err := setupBackupFile("backup.sqlite")
	if err != nil {
		log.Fatal(err)
	}

	err = runBackup("db.sqlite", "backup.sqlite")
	if err != nil {
		log.Fatal(err)
	}

	err = runIntegrityCheck("backup.sqlite")
	if err != nil {
		log.Fatal(err)
	}
}
