package db

const (
	CREATE_SESSIONS_TABLE = `
	CREATE TABLE IF NOT EXISTS sessions (
		id PRIMARY KEY AUTO INCREMENT
		status VARCHAR(50)
	)
	`

	INDEX_SESSIONS_TABLE_ON_STATUS = `
	CREATE INDEX IF NOT EXISTS sessions_status_index
	ON sessions(status)
	`

	CREATE_NODES_TABLE = `
	CREATE TABLE IF NOT EXISTS nodes (
		id PRIMARY KEY AUTO INCREMENT
		FOREIGN_KEY(session) REFERENCES sessions(id)
		status VARCHAR(50)
	)
	`

	INDEX_NODES_TABLE_ON_STATUS = `
	CREATE INDEX IF NOT EXISTS nodes_status_index
	ON nodes(status)
	`

	INDEX_NODES_TABLE_ON_SESSION = `
	CREATE INDEX IF NOT EXISTS nodes_session_index
	ON nodes(session)
	`

	CREATE_NODES_MAPPING_TABLE = `
	CREATE TABLE IF NOT EXISTS nodes_mapping (
		id PRIMARY KEY AUTO INCREMENT
		FOREIGN_KEY(node) REFERENCES nodes(id)
		FOREIGN_KEY(relative) REFERENCES nodes(id)
		relation VARCHAR(50)
	)
	`

	INDEX_NODES_MAPPING_ON_RELATION = `
	CREATE INDEX IF NOT EXISTS nodes_mapping_relation_index
	ON nodes_mapping(relation)
	`

	INDEX_NODES_MAPPING_ON_NODES = `
	CREATE INDEX IF NOT EXISTS nodes_mapping_nodes_index
	ON nodes_mapping(node, relative)
	`

	CREATE_ITERATIONS_TABLE = `
	CREATE TABLE IF NOT EXISTS iterations (
		id PRIMARY KEY AUTO INCREMENT
		FOREIGN_KEY(session) REFERENCES sessions(id)
		FOREIGN_KEY(evaluated_node) REFERENCES nodes(id)
		status VARCHAR(50)
		pruned_nodes VARCHAR
	)
	`

	INDEX_ITERATIONS_TABLE_ON_STATUS = `
	CREATE INDEX IF NOT EXISTS iterations_status_index
	ON iterations(status)
	`

	INDEX_ITERATIONS_TABLE_ON_SESSION_AND_NODE = `
	CREATE INDEX IF NOT EXISTS iterations_session_and_node_index
	ON iterations(session, evaluated_node)
	`
)
