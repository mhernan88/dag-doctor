package sessions

import (
	"log/slog"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/mhernan88/dag-bisect/db"
	"github.com/mhernan88/dag-bisect/shared"
)


func NewSessionManager(cxn *sqlx.DB, l *slog.Logger) SessionManager {
	return SessionManager{cxn: cxn, l: l}
}

func NewDefaultSessionManager() (*SessionManager, *os.File, error) {
	l, f := shared.GetLogger()
	cxn, err := db.Connect()
	if err != nil {
		return nil, nil, err
	}
	sm := NewSessionManager(cxn, l)
	return &sm, f, nil
}

type SessionManager struct {
	cxn *sqlx.DB
	l *slog.Logger
	f *os.File
}
