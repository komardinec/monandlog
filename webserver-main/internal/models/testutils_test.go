package models

import (
	"database/sql"
	"os"
	"testing"
)

func NewTestDB(t *testing.T) *sql.DB {
	// Establish a new sql.DB connection pool for our test database.
	// Script will use multiple SQL statements, so we need to prepare
	// for it by adding the "multiStatements=true" parameter in the DSN.
	// This will instruct the driver to support multiple statement execution
	// in db.Exec() call
	db, err := sql.Open("mysql", "testserver:iamtest123@/test_the_model?parseTime=true&multiStatements=true")
	if err != nil {
		t.Fatal(err)
	}

	// Read the setup SQL script from the file and execute the statements,
	// closing the pool and t.Fatal() if something goes wrong
	script, err := os.ReadFile("./testdata/setup.sql")
	if err != nil {
		db.Close()
		t.Fatal(err)
	}
	_, err = db.Exec(string(script))
	if err != nil {
		db.Close()
		t.Fatal(err)
	}

	// Use t.Cleanup() to register a function, which will be called, when the current test (or sub-test)
	// which calls NewTestDB() has finished. In here we are executing teardown script and close the DB connection
	// pool
	t.Cleanup(func() {
		defer db.Close()

		script, err := os.ReadFile("./testdata/teardown.sql")
		if err != nil {
			t.Fatal(err)
		}
		_, err = db.Exec(string(script))
		if err != nil {
			t.Fatal(err)
		}
	})

	return db
}
