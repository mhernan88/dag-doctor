package sessions

import (
	"os"

	"github.com/google/uuid"
	"github.com/mhernan88/dag-bisect/data"
	"github.com/mhernan88/dag-bisect/db"
	"github.com/mhernan88/dag-bisect/shared"
	"github.com/urfave/cli/v2"
)

func newSession(ctx *cli.Context) error {
	_, err := data.LoadDAG(ctx.Args().Get(0))
	if err != nil {
		return err
	}

	dbHandle, err := db.Connect(ctx.String("db-filepath"))
	if err != nil {
		return err
	}

	tx, err := dbHandle.Beginx()
	if err != nil {
		return err
	}

	dagID := uuid.New().String()
	return db.InsertOneIntoSessions(tx, dagID, "new")
}

func getNewSessionsFlags() []cli.Flag {
	DBFilename, err := shared.GetDBFilename()
	if err != nil {
		os.Exit(1)
	}
	flags := []cli.Flag{
		&cli.StringFlag{
			Name:  "db-filepath",
			Value: DBFilename,
			Usage: "language for the greeting",
		},
	}
	return flags
}

var NewSessionCmd = cli.Command {
	Name: "new",
	Usage: "new session",
	Flags: getNewSessionsFlags(),
	Action: newSession,
}
