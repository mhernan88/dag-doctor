package db

import (
	"testing"
	"database/sql"
	_ "github.com/glebarez/go-sqlite"
)

func TestCreateTables(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Error(err)
	}

	err = CreateTables(db, false)
	if err != nil {
		t.Error(err)
	}
}
