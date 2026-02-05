package main

import (
	"context"
	"datalogger/database"
	"fmt"

	"golang.org/x/sync/errgroup"
)

func main() {
	db := database.CreateDatabaseWriterConnection()

	g, ctx := errgroup.WithContext(context.Background())

	fmt.Println("Writer active")
}
