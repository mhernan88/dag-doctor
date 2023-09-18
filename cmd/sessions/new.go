package sessions

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mhernan88/dag-bisect/data"
	"github.com/mhernan88/dag-bisect/db"
	"github.com/mhernan88/dag-bisect/shared"
	"github.com/urfave/cli/v2"
)

func newSession(ctx *cli.Context) error {
	dagFilename := ctx.Args().Get(0)
	_, err := data.LoadDAG(dagFilename)
	if err != nil {
		return err
	}

	cxn, err := db.Connect()
	if err != nil {
		return err
	}
	defer cxn.Close()

	id := uuid.NewString()
	dt := time.Now().Unix()
	savedDagFilename, err := shared.CopyDAGToRepo(dagFilename, id)

	_, err = cxn.Exec(
		`INSERT INTO sessions (
			id, 
			file,
			splits,
			status, 
			meta_created_datetime, 
			meta_updated_datetime
		) VALUES (?, ?, 0, ?, ?, ?)`, 
		id, savedDagFilename, "new", dt, dt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert new session | %v", err)
	}

	fmt.Printf("created session %s\n", id)
	return nil
}


var NewSessionCmd = cli.Command {
	Name: "new",
	Usage: "creates a new session: usage -> ...session new <dag-filename>",
	Action: newSession,
}
