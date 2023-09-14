package db

import "database/sql"

func SelectAllFromSessions(dbHandle *sql.Tx) (*sql.Rows, error) {
	constructor := getSessionsTableConfig()
	rows, err := constructor.RenderAndExecuteSelect(dbHandle)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
