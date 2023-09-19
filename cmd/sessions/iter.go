package sessions

import (
	"log/slog"

	"github.com/jmoiron/sqlx"
	"github.com/urfave/cli/v2"
)

func NewIterSessionManager(cxn sqlx.DB, l *slog.Logger) IterSessionManager {
	return IterSessionManager{
		cxn: cxn,
		l: l,
	}
}

type IterSessionManager struct {
	cxn sqlx.DB
	l *slog.Logger
}

func iterSession(ctx *cli.Context) error {
	return nil
}
