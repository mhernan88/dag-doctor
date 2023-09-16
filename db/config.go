package db

const DROP_TABLE_TEMPLATE = "DROP TABLE IF EXISTS {{.Name}}"

func getSessionsTableConfig() SQLTableConstructor {
	sessionsTable := SQLTable{
		Name: "sessions",
	}
	return SQLTableConstructor{
		Table: sessionsTable,
		CreateTemplate:  `
		CREATE TABLE IF NOT EXISTS {{.Name}} (
			id VARCHAR(36),
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
		SelectTemplate: `
		SELECT * FROM {{.Name}}
		`,
		InsertOneTemplate: `
		INSERT INTO {{.Name}} (id, status) VALUES(?, ?)
		`,
	}
}

func getNodesTableConfig() SQLTableConstructor {
	nodesTable := SQLTable{
		Name: "nodes",
		Fk1: "sessions",
	}
	return SQLTableConstructor{
		Table: nodesTable,
		CreateTemplate: `
		CREATE TABLE IF NOT EXISTS {{.Name}} (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			session VARCHAR(36),
			status VARCHAR(50),
			FOREIGN KEY(session) REFERENCES {{.Fk1}}(id)
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
		SelectTemplate: `
		SELECT * FROM {{.Name}}
		`,
	}
}

func getNodesMappingTableConfig() SQLTableConstructor {
	nodesMappingTable := SQLTable{
		Name: "nodes_mapping",
		Fk1: "nodes",
	}
	return SQLTableConstructor{
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
		SelectTemplate: `
		SELECT * FROM {{.Name}}
		`,
	}
}

func getIterationsTableConfig() SQLTableConstructor{
	iterationsTable := SQLTable{
		Name: "iterations",
		Fk1: "sessions",
		Fk2: "nodes",
	}
	return SQLTableConstructor{
		Table: iterationsTable,
		CreateTemplate: `
		CREATE TABLE IF NOT EXISTS {{.Name}} (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			session VARCHAR(36),
			evaluated_node INTEGER,
			status VARCHAR(50),
			pruned_nodes VARCHAR,
			FOREIGN KEY(session) REFERENCES {{.Fk1}}(id),
			FOREIGN KEY(evaluated_node) REFERENCES {{.Fk2}}(id)
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
		SelectTemplate: `
		SELECT * FROM {{.Name}}
		`,
	}
}
