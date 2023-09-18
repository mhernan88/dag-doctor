package sessions

import (
	"fmt"
	"time"

	"github.com/mhernan88/dag-bisect/db"
	"github.com/urfave/cli/v2"
)

func updateSession(ctx *cli.Context) error {
	cxn, err := db.Connect()
	if err != nil {
		return err
	}
	defer cxn.Close()

	id := ctx.Args().Get(0)
	status := ctx.Args().Get(1)
	dt := time.Now().Unix()

	query := fmt.Sprintf(
		`UPDATE sessions
		SET 
			status = '%s', 
			meta_updated_datetime = '%d' 
		WHERE id = '%s'`,
		status, dt, id,
	)

	_, err = cxn.Exec(query, status, id)

	if err != nil {
		return fmt.Errorf("failed to update session | %v", err)
	}

	fmt.Printf("updated session '%s' to status='%s'\n", id, status)
	return nil
}

var UpdateSessionCmd = cli.Command {
	Name: "update",
	Usage: "...session dev update <session-id> <new-status>",
	Action: updateSession,
}
