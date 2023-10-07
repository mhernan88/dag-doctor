package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

func CreateSessionsTable(cxn *sqlx.DB, drop bool) error {
	if drop {
		query := "DROP TABLE IF EXISTS sessions"
		_, err := cxn.Exec(query)
		if err != nil {
			return fmt.Errorf("failed to drop sessions table | %v", err)
		}
	}
	query := `
	CREATE TABLE IF NOT EXISTS sessions (
		id VARCHAR(36),
		dag VARCHAR,
		state VARCHAR,
		splits INT,
		status VARCHAR(50),
		err_node VARCHAR,
		meta_created_datetime BIGINT,
		meta_updated_datetime BIGINT
	)`
	_, err := cxn.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create sessions table | %v", err)
	}
	return nil
}

func CreateTables(cxn *sqlx.DB, drop bool) error {
	return CreateSessionsTable(cxn, drop)
}
