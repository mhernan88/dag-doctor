package db

import (
	"os"
	"errors"
	_ "github.com/glebarez/go-sqlite"
	"database/sql"
)

func Connect(filepath string) (*sql.DB, error){
	if _, err := os.Stat(filepath); errors.Is(err, os.ErrNotExist) {
		file, err := os.Create(filepath)
        if err != nil {
			return nil, err
        }
        file.Close()
	}
	dbHandle, err := sql.Open("sqlite", filepath)
	if err != nil {
		return nil, err
	}
	return dbHandle, nil
}
