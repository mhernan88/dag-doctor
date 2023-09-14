package db

import (
	"database/sql"
	"fmt"
)

const DROP_TABLE_TEMPLATE = "DROP TABLE IF EXISTS {{.Name}}"

func CreateSessionsTable(dbHandle *sql.Tx, drop bool) error {
	sessionsTable := SQLTable{
		Name: "sessions",
	}
	constructor := SQLTableConstructor{
		Table: sessionsTable,
		CreateTemplate:  `
		CREATE TABLE IF NOT EXISTS {{.Name}} (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			status VARCHAR(50)
		)
		`,
		DropTemplate: DROP_TABLE_TEMPLATE,
		IndexTemplates: []string{
			`
			CREATE INDEX IF NOT EXISTS {{.Name}}_status_index
			ON {{.Name}}(status)
			`,
		},
	}

	err := constructor.RenderAndExecute(dbHandle, drop)
	if err != nil {
		return fmt.Errorf("failed to RenderAndExecute sessions table | %v", err)
	}
	return nil
}

func CreateNodesTable(dbHandle *sql.Tx, drop bool) error {
	nodesTable := SQLTable{
		Name: "nodes",
	}
	constructor := SQLTableConstructor{
		Table: nodesTable,
		CreateTemplate: `
		CREATE TABLE IF NOT EXISTS {{.Name}} (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			session INTEGER,
			status VARCHAR(50),
			FOREIGN KEY(session) REFERENCES sessions(id)
		)
		`,
		DropTemplate: DROP_TABLE_TEMPLATE,
		IndexTemplates: []string{
			`
			CREATE INDEX IF NOT EXISTS {{.Name}}_status_index
			ON {{.Name}}(status)
			`,

			`
			CREATE INDEX IF NOT EXISTS {{.Name}}_session_index
			ON {{.Name}}(session)
			`,

		},
	}

	err := constructor.RenderAndExecute(dbHandle, drop)
	if err != nil {
		return fmt.Errorf("failed to RenderAndExecute nodes table | %v", err)
	}
	return nil
}

func CreateNodesMappingTable(dbHandle *sql.Tx, drop bool) error {
	nodesMappingTable := SQLTable{
		Name: "nodes_mapping",
		Fk1: "nodes",
	}
	constructor := SQLTableConstructor{
		Table: nodesMappingTable,
		CreateTemplate: `
		CREATE TABLE IF NOT EXISTS {{.Name}} (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			node INTEGER,
			relative INTEGER,
			relation VARCHAR(50),
			FOREIGN KEY(node) REFERENCES {{.Fk1}}(id),
			FOREIGN KEY(relative) REFERENCES {{.Fk1}}(id)
		)
		`,
		DropTemplate: DROP_TABLE_TEMPLATE,
		IndexTemplates: []string{
			`
			CREATE INDEX IF NOT EXISTS {{.Name}}_relation_index
			ON {{.Name}}(relation)
			`,

			`
			CREATE INDEX IF NOT EXISTS {{.Name}}_nodes_index
			ON {{.Name}}(node, relative)
			`,
		},
	}
	err := constructor.RenderAndExecute(dbHandle, drop)
	if err != nil {
		return fmt.Errorf("failed to RenderAndExecute nodes_mapping table | %v", err)
	}
	return nil
}

func CreateIterationsTable(dbHandle *sql.Tx, drop bool) error {
	iterationsTable := SQLTable{
		Name: "iterations",
	}
	constructor := SQLTableConstructor{
		Table: iterationsTable,
		CreateTemplate: `
		CREATE TABLE IF NOT EXISTS {{.Name}} (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			session INTEGER,
			evaluated_node INTEGER,
			status VARCHAR(50),
			pruned_nodes VARCHAR,
			FOREIGN KEY(session) REFERENCES sessions(id),
			FOREIGN KEY(evaluated_node) REFERENCES nodes(id)
		)
		`,
		DropTemplate: DROP_TABLE_TEMPLATE,
		IndexTemplates: []string{
			`
			CREATE INDEX IF NOT EXISTS {{.Name}}_status_index
			ON {{.Name}}(status)
			`,

			`
			CREATE INDEX IF NOT EXISTS {{.Name}}_session_and_node_index
			ON {{.Name}}(session, evaluated_node)
			`,

		},
	}
	err := constructor.RenderAndExecute(dbHandle, drop)
	if err != nil {
		return fmt.Errorf("failed to RenderAndExecute iterations table | %v", err)
	}
	return nil
}

func CreateTables(dbHandle *sql.DB, drop bool) error {
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
