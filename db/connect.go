package db

import (
	"errors"
	"fmt"
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
			return nil, fmt.Errorf("failed to create db file | %v", err)
        }
        file.Close()
	}
	dbHandle, err := sqlx.Connect("sqlite", filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to sqlite db | %v", err)
	}
	return dbHandle, nil
}
