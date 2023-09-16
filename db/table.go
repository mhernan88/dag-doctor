package db

import (
	"database/sql"
	"fmt"
	"strings"
	"text/template"

	"github.com/jmoiron/sqlx"
	"github.com/mhernan88/dag-bisect/db/models"
)

type SQLTable struct {
	Name string
	Fk1 string
	Fk2 string
	Fk3 string
	Fk4 string
	Fk5 string
}

type SQLTableConstructor struct {
	Table SQLTable
	CreateTemplate string
	DropTemplate string
	IndexTemplates []string
	SelectTemplate string
	InsertOneTemplate string
}

func (stc SQLTableConstructor) RenderAndExecuteCreate(db *sql.Tx) error {
	createQuery := strings.Builder{}
	compiledCreateTmpl, err := template.New("sql").Parse(stc.CreateTemplate)
	if err != nil {
		return fmt.Errorf("failed to initialize create template: %v", err)
	}

	err = compiledCreateTmpl.Execute(&createQuery, stc.Table)
	if err != nil {
		return fmt.Errorf("failed to render create template: %v", err)
	}

	_, err = db.Exec(createQuery.String())
	if err != nil {
		return fmt.Errorf(
			"failed to execute create query: %s | %v", 
			createQuery.String(),
			err,
		)
	}
	return nil
}

func (stc SQLTableConstructor) RenderAndExecuteDrop(db *sql.Tx) error {
	dropQuery := strings.Builder{}
	compiledDropTmpl, err := template.New("sql").Parse(stc.DropTemplate)
	if err != nil {
		return fmt.Errorf("failed to initialize drop template: %v", err)
	}

	err = compiledDropTmpl.Execute(&dropQuery, stc.Table)
	if err != nil {
		return fmt.Errorf("failed to render drop template: %v", err)
	}

	_, err = db.Exec(dropQuery.String())
	return fmt.Errorf(
		"failed to execute drop query: %s | %v",
		dropQuery.String(), 
		err,
	)
}

func (stc SQLTableConstructor) RenderAndExecuteIndex(db *sql.Tx) error {
	for _, indexTemplate := range stc.IndexTemplates {
		indexQuery := strings.Builder{}
		compiledIndexTmpl, err := template.New("sql").Parse(indexTemplate)
		if err != nil {
			return fmt.Errorf("failed to initialize index template | %v", err)
		}

		err = compiledIndexTmpl.Execute(&indexQuery, stc.Table)
		if err != nil {
			return fmt.Errorf("failed to render index template | %v", err)
		}

		_, err = db.Exec(indexQuery.String())
		if err != nil {
			return fmt.Errorf(
				"failed to execute index query: %s | %v", 
				indexQuery.String(), 
				err,
			)
		}
	}
	return nil
}

func (stc SQLTableConstructor) RenderAndExecuteSelect(db *sqlx.Tx) ([]models.Session, error) {
	selectQuery := strings.Builder{}
	compiledSelectTmpl, err := template.New("sql").Parse(stc.SelectTemplate)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize select template | %v", err)
	}

	err = compiledSelectTmpl.Execute(&selectQuery, stc.Table)
	if err != nil {
		return nil, fmt.Errorf("failed to render select template | %v", err)
	}

	var sessions []models.Session
	err = db.Select(&sessions, selectQuery.String())
	if err != nil {
		return nil, err
	}
	return sessions, nil
}

func (stc SQLTableConstructor) RenderAndExecute(
	db *sql.Tx, drop bool) error {
	var err error
	if drop {
		err = stc.RenderAndExecuteDrop(db)
		if err != nil {
			return err
		}
	}

	err = stc.RenderAndExecuteCreate(db)
	if err != nil {
		return err
	}

	return stc.RenderAndExecuteIndex(db)
}

func (stc SQLTableConstructor) RenderAndInsertSession(
	db *sqlx.Tx, id string, status string,
) error {
	insertQuery := strings.Builder{}
	compiledInsertTmpl, err := template.New("sql").Parse(stc.InsertOneTemplate)
	if err != nil {
		return fmt.Errorf("failed to initialize insert one template | %v", err)
	}

	err = compiledInsertTmpl.Execute(&insertQuery, stc.Table)
	if err != nil {
		return err
	}

	_, err = db.Exec(insertQuery.String(), id, status)
	if err != nil {
		return fmt.Errorf(
			"failed to execute insert query: %s | %v",
			insertQuery.String(),
			err,
		)
	}
	return nil
}
