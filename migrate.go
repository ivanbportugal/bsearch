package main

import (
	"bufio"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"os"
	"strings"
	"time"
)

func main() {

	fmt.Println("Dropping and Migrating raw to Bolt")

	// Blow away DB
	dbsLocation := "translations/"
	kjvDb := "kjv.db"
	err := os.Remove(dbsLocation + kjvDb)
	if err != nil {
		fmt.Println(err)
		// Continue anyway
	}

	// New DB
	db, err := bolt.Open(dbsLocation+kjvDb, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Read from raw file
	file, err := os.Open("raw/kjv.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	copyFileToDB(file, db)
}

func copyFileToDB(file *os.File, db *bolt.DB) {

	// Start a writable transaction.
	tx, errBeginTransaction := db.Begin(true)
	if errBeginTransaction != nil {
		log.Fatal("Error beginning transaction: ", errBeginTransaction)
	}
	defer tx.Rollback()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// Per line
		currLine := scanner.Text()
		referenceAndVerse := strings.Split(currLine, "\t")
		reference := strings.Split(referenceAndVerse[0], " ")
		book := reference[0]
		chapterVerse := reference[1]
		verseText := referenceAndVerse[1]
		if len(reference) == 3 {
			// 1 Corinth, 1 Peter, etc...
			book = reference[0] + " " + reference[1]
			chapterVerse = reference[2]
		} else if len(reference) == 4 {
			// Song of Solomon
			book = reference[0] + " " + reference[1] + " " + reference[2]
			chapterVerse = reference[3]
		}

		// One bucket per book
		bucket, errBucket := tx.CreateBucketIfNotExists([]byte(book))
		if errBucket != nil {
			log.Fatal("create bucket: %s", errBucket)
		}

		// Key Value pair - reference to verse
		errKey := bucket.Put([]byte(chapterVerse), []byte(verseText))
		if errKey != nil {
			log.Fatal("Add key to bucket: %s", errKey)
		}

		// Debug output
		// fmt.Println(book + "-" + chapterVerse + "--" + verseText)
	}
	if errScanner := scanner.Err(); errScanner != nil {
		log.Fatal(errScanner)
	}

	// Commit the transaction and check for error.
	if errCommit := tx.Commit(); errCommit != nil {
		log.Fatal("Error commiting transaction: ", errCommit)
	}
}
