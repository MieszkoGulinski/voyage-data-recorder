package main

import (
	"datalogger/backup"
	"flag"
	"log"
	"os"
	"path/filepath"
)

func main() {
	diagnostics := flag.Bool("diagnostics", false, "Write progress information?")
	inputFile := flag.String("input-file", "", "Path to the database file to be backed up")
	outputFile := flag.String("output-file", "", "Path to the backup file to be created")
	dir := flag.String("dir", "", "Directory where backups are stored")
	flag.Parse()

	inputPath := *inputFile
	if inputPath == "" {
		inputPath = "db.sqlite"
	}

	if *outputFile != "" && *dir != "" {
		log.Fatal("Cannot use --output-file and --dir together")
	}

	outputPath := ""
	if *outputFile != "" {
		outputPath = *outputFile
	} else if *dir != "" {
		// Create directory if it doesn't exist
		err := os.MkdirAll(*dir, 0755)
		if err != nil {
			log.Fatalf("Failed to create directory %s: %v", *dir, err)
		}
		outputPath = filepath.Join(*dir, "backup.sqlite")
	} else {
		outputPath = "backup.sqlite"
	}

	err := backup.SetupBackupFile(outputPath, *diagnostics)
	if err != nil {
		log.Fatal(err)
	}

	err = backup.RunBackup(inputPath, outputPath, *diagnostics)
	if err != nil {
		log.Fatal(err)
	}

	err = backup.RunIntegrityCheck(outputPath, *diagnostics)
	if err != nil {
		log.Fatal(err)
	}
}
