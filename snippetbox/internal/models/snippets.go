package models

import (
	"database/sql"
	"errors"
	"time"
)

// define struct to hold data for an individual snippet.
// fields should
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// define a SnippetModel type which wraps a sql.DB connection pool.
type SnippetModel struct {
	DB *sql.DB
}

// this will insert a new snippet into the database.
func (m *SnippetModel) Insert(title, content string, expires int) (int, error) {

	// SQL statement we want to run
	stmt := `
		INSERT INTO 
			snippets (
				title, content, created, expires
			)
		VALUES(
			?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY)
		);
	`

	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	// check LastInsertId() to get the ID of the newly inserted record
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// convert int64 to int type
	return int(id), nil
}

// this will return a specific snippet based on its id.
func (m *SnippetModel) Get(id int) (*Snippet, error) {

	// SQL statement to get specific id from database
	stmt := `
		SELECT
			id,
		 	title,
			content,
			created,
			expires
		FROM
			snippets
		WHERE
			expires > UTC_TIMESTAMP() 
			AND id = ?;
	`
	row := m.DB.QueryRow(stmt, id)

	// create a pointer to zeroed Snippet struct
	s := &Snippet{}

	// row.Scan() will copy the values from each field in sql.Row to the corresponding field in the Snippet struct.
	// note that the arguments to row.Scan() are pointers to the place we want to copy the data into.
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {

		// if the query returns no rows, then row.Scan() will return a sql.ErrNoRows error.
		// errors.Is() will check for that error and return ErrNoRecord err instead
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}

// this will return the 10 most recently created snippets.
func (m *SnippetModel) Latest() ([]*Snippet, error) {

	// write the SQL statement we want to execute
	stmt := `
		SELECT
			id,
			title,
			content,
			created,
			expires 
		FROM
			snippets
		WHERE 
			expires > UTC_TIMESTAMP() 
		ORDER BY 
			id 
		DESC LIMIT 10;
	`

	// use Query() on the conn pool to execute our statement. This will return a sql.Rows resultset containing the result of our query.
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	// defer rows.Close() after we check for an error so we do not attempt to close on a nil resultset
	defer rows.Close()

	// initialize empty slice to hold the Snippet structs
	snippets := []*Snippet{}

	// use rows.Next() to iterate through the rows of the resultset
	// this prepares the first and subsequent row to be acted upon using the rows.Scan() method
	for rows.Next() {
		s := &Snippet{}
		// use rows.Scan() to copy the values from each field in teh row to our Snippet struct
		// the arguments to row.Scan() must be pointers
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}

	// when the rows.Next() loop has finished, we call rows.Err() to retrieve any errors that were encountered during iteration
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return snippets, nil
}
