package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"os"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

// The file which will be used for the database.
const DB_FILE = "vault.db"

// Struct to represent the quote table.
type databaseItemHandler struct {
	database  *sql.DB
	quotes    []item
	quoteIds  []int
	listIndex int
	mtx       *sync.Mutex
	shuffle   *sync.Once
}

// Initialize a database using the DB_FILE name.
func (d databaseItemHandler) resetDbItemHandler() error {
	// Set mtx and shuffle
	d.mtx = &sync.Mutex{}
	d.shuffle = &sync.Once{}

	var err error

	// Open database connection
	d.database, err = sql.Open("sqlite3", DB_FILE)
	if err != nil {
		return fmt.Errorf("error: %w", ErrOpenDb)
	}

	// Setup database data (if empty)
	err = CreateQuoteTable(d.database)
	if err != nil {
		return err
	}

	// Get quotes from the database
	// In the format of []item
	d.quotes, d.quoteIds, err = GetAllQuotes(d.database)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	d.listIndex = 0

	// Shuffle the array of quotes
	d.shuffle.Do(func() {
		SpliceShuffle(d.quotes)
		SpliceShuffle(d.quoteIds)
	})

	return nil
}

func SpliceShuffle[T any](s []T) {
	rand.Shuffle(len(s), func(i, j int) { s[i], s[j] = s[j], s[i] })
}

// Creates quote table if it does not exist.
func CreateQuoteTable(db *sql.DB) error {
	// Create query string.
	queryString := "CREATE TABLE IF NOT EXISTS quote (id INTEGER PRIMARY KEY, body TEXT NOT NULL, author TEXT)"

	// Prepare query to create table with the following properties:
	// id: 		integer primary key
	// body: 	string (cannot be null)
	// author: 	string
	// date: 	string (YYYY-MM-DD format. Can be just YYYY)
	_, err := db.Exec(queryString)

	// Check error.
	if err != nil {
		return fmt.Errorf("error: %w", ErrDbCreation)
	}

	return nil
}

func GetAllQuotes(db *sql.DB) ([]item, []int, error) {

	// Create query string.
	queryString := "SELECT id, body, author FROM quote"

	// Execute the string.
	rows, err := db.Query(queryString)
	if err != nil {
		return []item{}, []int{}, fmt.Errorf("error: %w", ErrMalformedQuery)
	}

	// Initialize slice of quotes.
	var items []item
	var ids []int

	// Iterate through rows and save to quotes.
	for rows.Next() {
		var i item
		var id int
		if err := rows.Scan(&id, &i.title, &i.desc); err != nil {
			return []item{}, []int{}, fmt.Errorf("error: %w", ErrMalformedQuery)
		}
		items = append(items, i)
	}

	// Check if any errors occured while iterating over rows.
	if err := rows.Err(); err != nil {
		return []item{}, []int{}, fmt.Errorf("error: %w", ErrRowIteration)
	}

	return items, ids, nil
}

func (d databaseItemHandler) GetNextQuote() item {
	if d.mtx == nil {
		d.resetDbItemHandler()
	}

	d.mtx.Lock()
	defer d.mtx.Unlock()

	i := item{
		title: d.quotes[d.listIndex].title,
		desc:  d.quotes[d.listIndex].desc,
	}

	d.listIndex++
	if d.listIndex >= len(d.quotes) {
		d.listIndex = 0
	}

	return i
}

// Insert a quote into the database.
func (d databaseItemHandler) AddQuote(db *sql.DB, i item) error {
	// Create query string.
	queryString := "INSERT INTO quote (body, author) VALUES (?, ?)"

	// Construct query using query string and parameters.
	_, err := db.Exec(queryString, i.title, i.desc)

	// Check error.
	if err != nil {
		return fmt.Errorf("error: %w", ErrMalformedQuery)
	}

	return nil
}

// Update a quote entry in the database.
func (d databaseItemHandler) UpdateQuote(db *sql.DB, i item, id int) error {
	// Query string to update quote in database.
	queryString := "UPDATE quote SET body = ?, author = ?, date = ? WHERE id = ?"

	_, err := db.Exec(queryString, i, id)
	if err != nil {
		return fmt.Errorf("error: %w", ErrMalformedQuery)
	}

	return nil
}

// Delete a quote from the database given an id.
func (d databaseItemHandler) DeleteQuote(db *sql.DB, id int) error {
	// Create query string.
	queryString := "DELETE FROM quote WHERE id = ?"

	// Construct query using query string.
	_, err := db.Exec(queryString, id)

	// Check error.
	if err != nil {
		return fmt.Errorf("error: %w", ErrMalformedQuery)
	}

	return nil
}
