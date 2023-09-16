package sessions

import (
	"fmt"
	"os"

	"github.com/mhernan88/dag-bisect/db"
	"github.com/mhernan88/dag-bisect/shared"
	"github.com/urfave/cli/v2"
)

func list(ctx *cli.Context) error {
	dbHandle, err := db.Connect(ctx.String("db-filepath"))
	if err != nil {
		return err
	}

	tx, err := dbHandle.Beginx()
	if err != nil {
		return err
	}
	
	rows, err := db.SelectAllFromSessions(tx)
	if err != nil {
		tx.Rollback()
	}
	tx.Commit()

	fmt.Println(rows)
	return nil
}

func getFlags() []cli.Flag {
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

var ListCmd = cli.Command {
	Name: "ls",
	Usage: "list sessions",
	Action: list,
	Flags: getFlags(),
}
