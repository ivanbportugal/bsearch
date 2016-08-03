package main

import (
	// "encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

var db *bolt.DB

// type stanza struct {
// 	book      string
// 	reference string
// 	verse     string
// }

func main() {

	fmt.Println("Starting Bible Search Server")

	dbsLocation := "translations/"
	kjvDb := "kjv.db"

	// Open DB
	var err error
	db, err = bolt.Open(dbsLocation+kjvDb, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Debug
	// query := "love"
	// results := queryDb(query, db)
	// fmt.Printf("results %d \n", len(results))

	// Begin Mux
	r := mux.NewRouter()
	r.HandleFunc("/search", SearchHandler).Methods("POST")
	http.Handle("/", r)

	errHttp := http.ListenAndServe(":8080", r)
	if errHttp != nil {
		log.Println("Could not launch HTTP server: " + errHttp.Error())
	}
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {

	val := r.URL.Query().Get("query")
	results := queryDb(val, db)

	w.Header().Set("Content-Type", "application/json")
	toReturn := "[" + strings.Join(results, ",") + "]"
	io.WriteString(w, toReturn)
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
		if CaseInsensitiveContains(book, query) || CaseInsensitiveContains(string(verseText[:]), query) || CaseInsensitiveContains(string(reference[:]), query) {
			// fmt.Printf("%s %s %s\n", book, reference, verseText)
			// j, _ := json.Marshal(stanza{book: book, reference: string(reference), verse: string(verseText)})
			result := fmt.Sprintf("{\"book\":\"%s\",\"reference\":\"%s\",\"verse\":\"%s\"}", book, string(reference), string(verseText))
			results = append(results, result)
		}
	}
	return results
}

func CaseInsensitiveContains(s, substr string) bool {
	s, substr = strings.ToUpper(s), strings.ToUpper(substr)
	return strings.Contains(s, substr)
}
