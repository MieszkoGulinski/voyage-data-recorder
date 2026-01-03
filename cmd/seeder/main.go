package main

import (
	"datalogger/seeder"
	"log"
)

func main() {
	options, err := seeder.GetSeederConfig()
	if err != nil {
		log.Fatal(err)
	}

	db := seeder.ResetDatabase()
	seeder.FillDatabase(db, options)
}
