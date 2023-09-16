package db

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

func CreateSessionsTable(dbHandle *sql.Tx, drop bool) error {
	constructor := getSessionsTableConfig()
	err := constructor.RenderAndExecute(dbHandle, drop)
	if err != nil {
		return fmt.Errorf("failed to RenderAndExecute sessions table | %v", err)
	}
	return nil
}

func CreateNodesTable(dbHandle *sql.Tx, drop bool) error {
	constructor := getNodesTableConfig()
	err := constructor.RenderAndExecute(dbHandle, drop)
	if err != nil {
		return fmt.Errorf("failed to RenderAndExecute nodes table | %v", err)
	}
	return nil
}

func CreateNodesMappingTable(dbHandle *sql.Tx, drop bool) error {
	constructor := getNodesMappingTableConfig()
	err := constructor.RenderAndExecute(dbHandle, drop)
	if err != nil {
		return fmt.Errorf("failed to RenderAndExecute nodes_mapping table | %v", err)
	}
	return nil
}

func CreateIterationsTable(dbHandle *sql.Tx, drop bool) error {
	constructor := getIterationsTableConfig()
	err := constructor.RenderAndExecute(dbHandle, drop)
	if err != nil {
		return fmt.Errorf("failed to RenderAndExecute iterations table | %v", err)
	}
	return nil
}

func CreateTables(dbHandle *sqlx.DB, drop bool) error {
	var err error
	var tx *sql.Tx

	tx, err = dbHandle.Begin()
	if err != nil {
		return err
	}
	err = CreateSessionsTable(tx, drop)
	if err != nil {
		return err
	}

	err = CreateNodesTable(tx, drop)
	if err != nil {
		return err
	}
	err = CreateNodesMappingTable(tx, drop)
	if err != nil {
		return err
	}
	err = CreateIterationsTable(tx, drop)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}


	return nil
}
