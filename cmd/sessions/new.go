package sessions

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/mhernan88/dag-bisect/data"
	"github.com/urfave/cli/v2"
)

func newSession(ctx * cli.Context) error {
	_, err := data.LoadDAG(ctx.Args().Get(0))
	if err != nil {
		return err
	}
	dagID := uuid.New().String()
	fmt.Println(dagID)
	return nil
}

var NewSessionCmd = cli.Command {
	Name: "new",
	Usage: "new session",
	Action: newSession,
}
