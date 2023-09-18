package sessions

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/mhernan88/dag-bisect/db"
	"github.com/mhernan88/dag-bisect/db/models"
	"github.com/urfave/cli/v2"
)

func updateSession(ctx *cli.Context) error {
	id := ctx.Args().Get(0)
	status := ctx.Args().Get(1)

	if !slices.Contains(models.SESSION_STATUSES, status) {
		sessionStatuses := strings.Join(models.SESSION_STATUSES, ", ")
		return fmt.Errorf("status must be one of: %s", sessionStatuses)
	}

	cxn, err := db.Connect()
	if err != nil {
		return err
	}
	defer cxn.Close()

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
