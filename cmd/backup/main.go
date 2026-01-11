package main

import (
	"datalogger/backup"
	"flag"
	"log"
)

func main() {
	diagnostics := flag.Bool("diagnostics", false, "Write progress information?")
	flag.Parse()

	err := backup.SetupBackupFile("backup.sqlite", *diagnostics)
	if err != nil {
		log.Fatal(err)
	}

	err = backup.RunBackup("db.sqlite", "backup.sqlite", *diagnostics)
	if err != nil {
		log.Fatal(err)
	}

	err = backup.RunIntegrityCheck("backup.sqlite", *diagnostics)
	if err != nil {
		log.Fatal(err)
	}
}
