package db

import (
	"github.com/jmoiron/sqlx"
	"github.com/mhernan88/dag-bisect/db/models"
)

func SelectAllFromSessions(
	dbHandle *sqlx.Tx) ([]models.Session, error) {
	constructor := getSessionsTableConfig()
	sessions, err := constructor.RenderAndExecuteSelect(dbHandle)
	if err != nil {
		return nil, err
	}
	return sessions, nil
}

func InsertOneIntoSessions(dbHandle *sqlx.Tx, id, status string) error {
	constructor := getSessionsTableConfig()
	return constructor.RenderAndInsertSession(dbHandle, id, status)
}
