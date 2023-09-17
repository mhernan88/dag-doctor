package sessions

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/mhernan88/dag-bisect/data"
	"github.com/mhernan88/dag-bisect/db"
	"github.com/urfave/cli/v2"
)

func newSession(ctx *cli.Context) error {
	_, err := data.LoadDAG(ctx.Args().Get(0))
	if err != nil {
		return err
	}

	dbHandle, err := db.Connect()
	if err != nil {
		return err
	}
	defer dbHandle.Close()

	id := uuid.NewString()
	_, err = dbHandle.Exec(
		`INSERT INTO sessions (id, status) VALUES (?, ?)`, 
		id, "new",
	)
	if err != nil {
		return err
	}

	fmt.Printf("created session %s\n", id)
	return nil
}


var NewSessionCmd = cli.Command {
	Name: "new",
	Usage: "new session",
	Action: newSession,
}
