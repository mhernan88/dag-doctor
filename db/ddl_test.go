package db

import (
	"testing"

	_ "github.com/glebarez/go-sqlite"
	"github.com/jmoiron/sqlx"
)

func TestCreateTables(t *testing.T) {
	db, err := sqlx.Open("sqlite", ":memory:")
	if err != nil {
		t.Error(err)
	}

	err = CreateTables(db, false)
	if err != nil {
		t.Error(err)
	}
}
