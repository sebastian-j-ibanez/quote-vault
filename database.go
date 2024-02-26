package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// The file which will be used for the database.
const DB_FILE = "vault.db"

// Quote information.
type Quote struct {
	ID     int
	Body   string
	Author string
	Date   string
}

// Initialize a database using the DB_FILE name.
func InitDb() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", DB_FILE)
	if err != nil {
		return nil, fmt.Errorf("error: %w", ErrOpenDb)
	}

	err = CreateQuoteTable(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// Creates quote table if it does not exist.
func CreateQuoteTable(db *sql.DB) error {
	// Create query string.
	queryString := "CREATE TABLE IF NOT EXISTS quote (id INTEGER PRIMARY KEY, body TEXT NOT NULL, author TEXT, date TEXT)"

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

// Insert a quote into the database.
func AddQuote(db *sql.DB, body string, author string, date string) error {
	// Create query string.
	queryString := "INSERT INTO quote (body, author, date) VALUES (?, ?, ?)"

	// Construct query using query string and parameters.
	_, err := db.Exec(queryString, body, author, date)

	// Check error.
	if err != nil {
		return fmt.Errorf("error: %w", ErrMalformedQuery)
	}

	return nil
}

// Update a quote entry in the database.
func UpdateQuote(db *sql.DB, id int, body string, author string, date string) error {
	// Query string to update quote in database.
	queryString := "UPDATE quote SET body = ?, author = ?, date = ? WHERE id = ?"

	_, err := db.Exec(queryString, body, author, date, id)
	if err != nil {
		return fmt.Errorf("error: %w", ErrMalformedQuery)
	}

	return nil
}

// Delete a quote from the database given an id.
func DeleteQuote(db *sql.DB, id int) error {
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

func GetAllQuotes(db *sql.DB) ([]Quote, error) {

	// Create query string.
	queryString := "SELECT * FROM quote"

	// Execute the string.
	rows, err := db.Query(queryString)
	if err != nil {
		return []Quote{}, fmt.Errorf("error: %w", ErrMalformedQuery)
	}

	// Initialize slice of quotes.
	var quotes []Quote

	// Iterate through rows and save to quotes.
	for rows.Next() {
		var q Quote
		if err := rows.Scan(&q.ID, &q.Body, &q.Author, &q.Date); err != nil {
			return []Quote{}, fmt.Errorf("error: %w", ErrMalformedQuery)
		}
		quotes = append(quotes, q)
	}

	// Check if any errors occured while iterating over rows.
	if err := rows.Err(); err != nil {
		return []Quote{}, fmt.Errorf("error: %w", ErrRowIteration)
	}

	return quotes, nil
}

// Get quote from database.
func GetQuote(db *sql.DB, id int) (Quote, error) {

	// Create query string.
	queryString := "SELECT * FROM quote WHERE id = ?"

	// Construct query with query string.
	row := db.QueryRow(queryString, id)

	var q Quote
	if err := row.Scan(&q.ID, &q.Body, &q.Author, &q.Date); err != nil {
		return Quote{}, fmt.Errorf("error: %w", ErrRowIteration)
	}

	// Return quote.
	return q, nil
}
