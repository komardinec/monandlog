package models

import (
	"database/sql"
	"errors"
	"time"
)

// Define a File type structure for an individual file.
// That structure formed to correspond to the fields in database fields
type File struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// Define a WebApp model which wraps a sql.DB connection pool.
type WebAppModel struct {
	DB *sql.DB
}

// This will insert a new file into the database.
func (m *WebAppModel) Insert(title string, content string, expires int) (int, error) {
	// Write the SQL statement we want to execute.
	stmt := `INSERT INTO files (title, content, created, expires)
	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	// Use the Exec() method on embedded connection pool to execute the
	// statement. Exec() accepts as first argument a SQL statement string, next - in order values placed in
	// statement placeholder parameters. Returns a sql.Result, which contains some
	// basic information about what happened when the statement was executed.
	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, nil
	}

	// Use the LastInsertId() method on the result to get the ID of our
	// newly inserted record in the snippets table.
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// The ID returned has the type int64, so we convert it to an int type
	// before returning
	return int(id), nil
}

// This will return a specific file based on its id.
func (m *WebAppModel) Get(id int) (File, error) {
	// The SQL statement we want to execute.
	stmt := `SELECT id, title, content, created, expires FROM files
	WHERE expires > UTC_TIMESTAMP() AND id = ?`

	// Use the QueryRow() method on the connection pool to execute our
	// SQL statement, passing in the untrusted id variable as the value for the
	// placeholder parameter. This return a pointer to a sql.Row object which
	// holds the result from the database
	row := m.DB.QueryRow(stmt, id)

	// Initialize a new zeroed data row structure
	var f File

	// Use row.Scan() to copy the values from each field in sql.Row
	// to the corresponing field in the Snippet struct.
	// Important details: arguments must be pointers and exactly,
	// the same as the number of columns returned by statement
	err := row.Scan(&f.ID, &f.Title, &f.Content, &f.Created, &f.Expires)
	if err != nil {
		// If the query returns no rows, then row.Scan()
		// will return a sql.ErrNoRows error. We use
		// the errors.Is() func check for that error
		// specifically, and return our own ErrNoRecord error
		// instead
		if errors.Is(err, sql.ErrNoRows) {
			return File{}, ErrNoRecord
		} else {
			return File{}, err
		}
	}
	// If everyting is ok, return the filled Snippet struct
	return f, nil
}

// This will return the 10 most recently created files.
func (m *WebAppModel) Latest() ([]File, error) {
	// SQL query for 10 most recent (by id) data rows
	stmt := `SELECT id, title, content, created, expires FROM files
	WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`

	// Use Query() to execute SQL statement
	// on connection pool to retrieve a sql.Rows resultset
	// of the query
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	// Defer to ensure the sql.Rows resultset is always properly
	// closed before Latest() method returns. This declaration should be
	// after query error check, otherwise, if Query() will return with an error,
	// we will get an error trying to close nil resultset
	defer rows.Close()

	// Empty slice of results of type File struct
	var files []File

	// Use rows.Next() to iterate through the rows in the resultset. This prepares
	// every row in order to be acted on by the rows.Scan(). If iteration completes
	// over all the rows completes then the resultset automatically closes itself
	// and frees-up the underlying db connections
	for rows.Next() {
		// Create a new zeroed File struct
		var f File

		// Copy values from each field in the row to corresponding object field,
		// every argument must be pointers
		err = rows.Scan(&f.ID, &f.Title, &f.Content, &f.Created, &f.Expires)
		if err != nil {
			return nil, err
		}

		// Append it to the slice of Files
		files = append(files, f)
	}

	// After rows.Next() loop completed we call rows.Err() to retreive any
	// error that was encountered during the iteration. It's important to call this -
	// den't assume that a successfull iteration was completed over the whole resultset
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// If everything is OK - return slice
	return files, nil
}
