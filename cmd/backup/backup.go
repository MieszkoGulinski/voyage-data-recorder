package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/mattn/go-sqlite3"
)

func setupBackupFile(backupPath string) error {
	backupFileConn, err := sql.Open("sqlite3", "file:"+backupPath)
	if err != nil {
		return err
	}
	defer backupFileConn.Close()

	_, err = backupFileConn.Exec(`PRAGMA journal_mode = OFF;`)
	if err != nil {
		return err
	}

	_, err = backupFileConn.Exec(`PRAGMA synchronous = OFF;`)
	if err != nil {
		return err
	}

	_, err = backupFileConn.Exec(`PRAGMA temp_store = MEMORY;`)
	if err != nil {
		return err
	}

	return nil
}

func runBackup(sourcePath string, backupPath string) error {
	driver := sqlite3.SQLiteDriver{}
	sourceDriverConn, err := driver.Open("file:" + sourcePath + "?mode=ro")
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer sourceDriverConn.Close()

	destDriverConn, err := driver.Open("file:" + backupPath)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer destDriverConn.Close()

	srcConn, ok := sourceDriverConn.(*sqlite3.SQLiteConn)
	if !ok {
		return fmt.Errorf("not a sqlite3 connection")
	}

	destConn, ok := destDriverConn.(*sqlite3.SQLiteConn)
	if !ok {
		return fmt.Errorf("not a sqlite3 connection")
	}

	backup, err := destConn.Backup(
		"main", // always "main"
		srcConn,
		"main",
	)

	if err != nil {
		fmt.Println(err)
		return err
	}
	defer backup.Finish()

	for {
		done, err := backup.Step(100) // batch of 100 pages, each page has 4096 bytes, so each step processes ~400 kB
		if err != nil {
			fmt.Println(err)
			return err
		}

		if done {
			break
		}

		time.Sleep(10 * time.Millisecond) // prevents CPU and disk contention
	}

	return nil
}

func runIntegrityCheck(backupPath string) error {
	backupFileConn, err := sql.Open("sqlite3", "file:"+backupPath)
	if err != nil {
		return err
	}
	defer backupFileConn.Close()

	rows, err := backupFileConn.Query("PRAGMA integrity_check")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var result string
		err := rows.Scan(&result)

		if err != nil {
			return err
		}

		if result != "ok" {
			return fmt.Errorf("integrity check failed: %s", result)
		}
	}

	return nil
}
