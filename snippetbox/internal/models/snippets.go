package models

import (
	"database/sql"
	"errors"
	"time"
)

// Define struct to hold data for an individual snippet.
// Fields should
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// Define a SnippetModel type which wraps a sql.DB connection pool.
type SnippetModel struct {
	DB *sql.DB
}

// This will insert a new snippet into the database.
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {

	// SQL statement we want to run
	stmt := `insert into snippets (title, content, created, expires)
			 values(?, ?, utc_timestamp(), date_add(utc_timestamp(), interval ? day))`

	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	// Check LastInsertId() to get the ID of the newly inserted record
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// Convert int64 to int type
	return int(id), nil
}

// This will return a specific snippet based on its id.
func (m *SnippetModel) Get(id int) (*Snippet, error) {
	stmt := `select id, title, content, created, expires from snippets
			 where expires > utc_timestamp() and id = ?`

	row := m.DB.QueryRow(stmt, id)

	// Create a pointer to zeroed Snippet struct
	s := &Snippet{}

	// row.Scan() will copy the values from each field in sql.Row to the corresponding field in the Snippet struct.
	// Note that the arguments to row.Scan() are pointers to the place we want to copy the data into.
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {

		// If the query returns no rows, then row.Scan() will return a sql.ErrNoRows error.
		// errors.Is() will check for that error and return ErrNoRecord err instead
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	// If successful
	return s, nil
}

// This will return the 10 most recently created snippets.
func (m *SnippetModel) Latest() ([]*Snippet, error) {
	return nil, nil
}
