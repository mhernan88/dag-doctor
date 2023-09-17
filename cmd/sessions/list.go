package sessions

import (
	"fmt"

	"github.com/mhernan88/dag-bisect/db"
	"github.com/mhernan88/dag-bisect/db/models"
	"github.com/urfave/cli/v2"
)

func list(ctx *cli.Context) error {
	dbHandle, err := db.Connect()
	if err != nil {
		return err
	}

	statusFilter := ctx.String("status")

	var sessions []models.Session
	if statusFilter == "" {
		err = dbHandle.Select(
			&sessions,
			`SELECT * FROM sessions WHERE status != 'closed'`)
	} else if statusFilter == "all" {
		fmt.Println("filtering to all sessions")
		err = dbHandle.Select(
			&sessions,
			`SELECT * FROM sessions WHERE status = 'all'`,
		)
	} else {
		fmt.Printf("filtering to sessions with status = '%s'\n", statusFilter)
		err = dbHandle.Select(
			&sessions,
			`SELECT * FROM sessions WHERE status = '?'`,
			statusFilter)
	}
	if err != nil {
		return err
	}

	fmt.Println(sessions)
	return nil
}

var listFlags = []cli.Flag{
	&cli.StringFlag{
		Name: "status",
		Value: "",
		Usage: "session status filter",
	},
}

var ListCmd = cli.Command {
	Name: "ls",
	Usage: "list sessions",
	Action: list,
	Flags: listFlags, 
}
