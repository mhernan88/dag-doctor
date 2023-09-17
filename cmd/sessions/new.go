package sessions

import (
	"fmt"
	"time"

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
	dt := time.Now().Unix()
	_, err = dbHandle.Exec(
		`INSERT INTO sessions (
			id, 
			status, 
			meta_created_datetime, 
			meta_updated_datetime
		) VALUES (?, ?, ?, ?)`, 
		id, "new", dt, dt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert new session | %v", err)
	}

	fmt.Printf("created session %s\n", id)
	return nil
}


var NewSessionCmd = cli.Command {
	Name: "new",
	Usage: "new session",
	Action: newSession,
}
