package sessions

import (
	"log/slog"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/mhernan88/dag-bisect/db"
	"github.com/mhernan88/dag-bisect/shared"
)


func NewSessionManager(cxn *sqlx.DB, l *slog.Logger, f *os.File) SessionManager {
	return SessionManager{cxn: cxn, l: l, f: f}
}

func NewDefaultSessionManager() (*SessionManager, error) {
	l, f := shared.GetLogger()
	cxn, err := db.Connect()
	if err != nil {
		return nil, err
	}
	sm := NewSessionManager(cxn, l, f)
	return &sm, nil
}

type SessionManager struct {
	cxn *sqlx.DB
	l *slog.Logger
	f *os.File
}
