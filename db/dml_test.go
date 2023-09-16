package db

import (
	"testing"

	_ "github.com/glebarez/go-sqlite"
	"github.com/jmoiron/sqlx"
)

func TestInsertSession(t *testing.T) {
	dbHandle, err := sqlx.Open("sqlite", ":memory:")
	if err != nil {
		t.Error(err)
	}

	err = CreateTables(dbHandle, false)
	if err != nil {
		t.Error(err)
	}

	tx, err := dbHandle.Beginx()
	if err != nil {
		t.Error(err)
	}

	err = InsertOneIntoSessions(tx, "abc", "123")
	if err != nil {
		t.Error(err)
	}

	err = tx.Commit()
	if err != nil {
		t.Error(err)
	}

	tx, err = dbHandle.Beginx()
	if err != nil {
		t.Error(err)
	}

	sessions, err := SelectAllFromSessions(tx)
	if len(sessions) != 1 {
		t.Errorf("expected 1 session, got %d sessions", len(sessions))
	}

	err = tx.Commit()
}
