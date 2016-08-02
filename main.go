package main

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"strings"
	"time"
)

func main() {

	fmt.Println("Starting Bible Search Server")

	dbsLocation := "translations/"
	kjvDb := "kjv.db"

	// Open DB
	db, err := bolt.Open(dbsLocation+kjvDb, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	query := "love"

	results := queryDb(query, db)
	fmt.Printf("results %d \n", len(results))
}

func queryDb(query string, db *bolt.DB) []string {

	results := []string{}

	err := db.View(func(tx *bolt.Tx) error {
		tx.ForEach(func(name []byte, bucket *bolt.Bucket) error {
			// For each bucket (book)
			result := queryBucket(bucket, string(name[:]), query)
			results = append(results, result...)
			return nil
		})
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	return results
}

func queryBucket(bucket *bolt.Bucket, book string, query string) []string {

	results := []string{}
	c := bucket.Cursor()
	for reference, verseText := c.First(); reference != nil; reference, verseText = c.Next() {
		if CaseInsensitiveContains(string(reference[:]), query) || CaseInsensitiveContains(string(verseText[:]), query) {
			// fmt.Printf("%s %s %s\n", book, reference, verseText)
			result := fmt.Sprintf("%s %s %s", book, reference, verseText)
			results = append(results, result)
		}
	}
	return results
}

func CaseInsensitiveContains(s, substr string) bool {
	s, substr = strings.ToUpper(s), strings.ToUpper(substr)
	return strings.Contains(s, substr)
}
