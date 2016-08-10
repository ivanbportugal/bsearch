package main

import (
	"bytes"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
	// "encoding/json"
	// "io/ioutil"
)

var db *bolt.DB

// type stanza struct {
// 	book      string
// 	reference string
// 	verse     string
// }

// type searchterms struct {
// 	terms string `json:"terms"`
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
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))
	http.Handle("/", r)

	port := ":8080"
	fmt.Println("Launching on port " + port)

	errHttp := http.ListenAndServe(port, r)
	if errHttp != nil {
		log.Println("Could not launch HTTP server: " + errHttp.Error())
	}
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {

	// decoder := json.NewDecoder(r.Body)
	// var input struct {
	// 	terms string
	// }
	// err := decoder.Decode(&input)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(input)

	// var input struct {
	// 	terms string
	// }
	// b, e := ioutil.ReadAll(r.Body)
	// fmt.Println("FROM BODY! " + string(b))
	// json.Unmarshal(b, &input)
	// if e != nil {
	// 	panic(e)
	// }
	// fmt.Println("decoded: " + input.terms)

	val := r.URL.Query().Get("query")
	w.Header().Set("Content-Type", "application/json")

	if val == "" {
		toReturn := "[]"
		io.WriteString(w, string(toReturn))
	} else {
		vals := strings.Split(val, ",")
		results := queryDb(vals, db)

		toReturn := "[" + strings.Join(results, ",") + "]"
		// toReturn, _ := json.Marshal(input)
		io.WriteString(w, string(toReturn))
	}
}

func queryDb(queries []string, db *bolt.DB) []string {

	results := []string{}

	err := db.View(func(tx *bolt.Tx) error {
		tx.ForEach(func(name []byte, bucket *bolt.Bucket) error {
			// For each bucket (book)
			result := queryBucket(bucket, string(name[:]), queries)
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

func queryBucket(bucket *bolt.Bucket, book string, queries []string) []string {

	results := []string{}
	c := bucket.Cursor()
	for reference, verseText := c.First(); reference != nil; reference, verseText = c.Next() {
		if isInQueries(book, string(verseText[:]), string(reference[:]), queries) {
			// fmt.Printf("%s %s %s\n", book, reference, verseText)
			// j, _ := json.Marshal(stanza{book: book, reference: string(reference), verse: string(verseText)})
			result := fmt.Sprintf("{\"book\":\"%s\",\"reference\":\"%s\",\"verse\":\"%s\"}", book, string(reference), string(verseText))
			results = append(results, result)
		}
	}
	return results
}

func isInQueries(book, verseText, reference string, queries []string) bool {
	toReturn := false
	matchesCount := 0
	var normalized bytes.Buffer
	normalized.WriteString(book)
	normalized.WriteString(string(verseText[:]))
	normalized.WriteString(string(reference[:]))

	for _, query := range queries {
		if caseInsensitiveContains(normalized.String(), query) {
			// one match
			matchesCount++
		}
	}

	if matchesCount == len(queries) {
		// Does an AND on the tags
		toReturn = true
	}

	return toReturn
}

func caseInsensitiveContains(s, substr string) bool {
	s, substr = strings.ToUpper(s), strings.ToUpper(substr)
	return strings.Contains(s, substr)
}
