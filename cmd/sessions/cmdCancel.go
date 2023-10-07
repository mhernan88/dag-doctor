package sessions

import (
	"fmt"
	"github.com/urfave/cli/v2"
)

func cancelSession(ID string) error{
	sm, f, err := NewDefaultSessionManager()
	if err != nil {
		sm.l.Error(
			"listSessins command failed to create session manager",
			"err", err)
		return fmt.Errorf("failed to create session manager | %v", err)
	}
	defer f.Close()

	ID, err = sm.QuerySessionIDByPartialID(ID)
	if err != nil {
		return fmt.Errorf("failed to enrich session id | %v", err)
	}

	err = sm.UpdateSessionStatus(ID, "cancelled")
	if err != nil {
		return fmt.Errorf("failed to update session status | %v", err)
	}
	return nil
}

func cancelSessionFunc(ctx *cli.Context) error {
	return cancelSession(ctx.Args().Get(0))
}

var CancelSessionCmd = cli.Command {
	Name: "cancel",
	Usage: "cancel a session: usage -> ...session cancel <session-id>",
	Action: cancelSessionFunc,
}
