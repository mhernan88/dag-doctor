package db

import (
	"errors"
	"os"

	_ "github.com/glebarez/go-sqlite"
	"github.com/jmoiron/sqlx"
	"github.com/mhernan88/dag-bisect/shared"
)

func Connect() (*sqlx.DB, error){
	filepath, err := shared.GetDBFilename()
	if err != nil {
		return nil, err
	}
	if _, err := os.Stat(filepath); errors.Is(err, os.ErrNotExist) {
		file, err := os.Create(filepath)
        if err != nil {
			return nil, err
        }
        file.Close()
	}
	dbHandle, err := sqlx.Connect("sqlite", filepath)
	if err != nil {
		return nil, err
	}
	return dbHandle, nil
}
