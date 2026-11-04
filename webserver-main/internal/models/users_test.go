package models

import (
	"testing"

	"debutwebapp.komardinec/internal/assert"
)

func TestUserModelExists(t *testing.T) {
	// Check for -short flag,
	// if we don't want this long test
	// to run
	if testing.Short() {
		t.Skip("models: skipping integration test")
	}

	// Set uo a suite of table-driven tests and expected values
	tests := []struct {
		name   string
		userID int
		want   bool
	}{
		{
			name:   "Valid ID",
			userID: 1,
			want:   true,
		},
		{
			name:   "Zero ID",
			userID: 0,
			want:   false,
		},
		{
			name:   "Non-existent ID",
			userID: 99,
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call the NewTestDB() helper function to get a connection pool to
			// our test database. For each sub test, we will set up and torn the database tables
			// while calling this here
			db := NewTestDB(t)

			// Create a new instance of the UserModel
			u := UserModel{db}

			// Call the UserModel.Exists() method and check that the return
			// value and error match the expected values for the sub-test
			exists, err := u.Exists(tt.userID)

			assert.Equal(t, exists, tt.want)
			assert.NilError(t, err)
		})
	}
}
