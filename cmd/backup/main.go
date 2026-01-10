package main

import (
	"flag"
	"log"
)

func main() {
	diagnostics := flag.Bool("diagnostics", false, "Write progress information?")
	flag.Parse()

	err := setupBackupFile("backup.sqlite", *diagnostics)
	if err != nil {
		log.Fatal(err)
	}

	err = runBackup("db.sqlite", "backup.sqlite", *diagnostics)
	if err != nil {
		log.Fatal(err)
	}

	err = runIntegrityCheck("backup.sqlite", *diagnostics)
	if err != nil {
		log.Fatal(err)
	}
}
