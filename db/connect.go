package db

import (
	"os"
	"errors"
	_ "github.com/glebarez/go-sqlite"
	"github.com/jmoiron/sqlx"
)

func Connect(filepath string) (*sqlx.DB, error){
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
