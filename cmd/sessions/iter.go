package sessions

import (
	"fmt"

	"github.com/mhernan88/dag-bisect/cmd"
	"github.com/urfave/cli/v2"
)


func (sm SessionManager) iterSession(ID string) error{
	sessionModel, err := sm.QuerySessionByID(ID)
	if err != nil {
		return fmt.Errorf("failed to query session by id | %v", err)
	}

	ui, err := cmd.LoadState(sessionModel.State)
	if err != nil {
		return fmt.Errorf("failed to load state | %v", err)
	}

	err = ui.CheckDAGIter(sm.l)
	if err != nil {
		return fmt.Errorf("failed to evaluate dag | %v", err)
	}

	err = sm.IncrementSessionSplits(ID)
	if err != nil {
		return fmt.Errorf("failed to increment splits | %v", err)
	}

	if (len(ui.DAG.Nodes) == 0) || (len(ui.DAG.Roots) == 0) {
		ui.Terminate()
		if ui.LastFailedNode == "" {
			err = sm.UpdateSessionStatus(ID, "ok")
			if err != nil {
				return fmt.Errorf("failed to update session status | %v", err)
			}
		} else {
			err = sm.UpdateSessionStatus(ID, "err")
			if err != nil {
				return fmt.Errorf("failed to update session status | %v", err)
			}
		}
	} else {
		fmt.Printf("successfully iterate session %s\n", ID)
	}
	return nil
}



func iterSession(ID string) error {
	sm, f, err := NewDefaultSessionManager()
	if err != nil {
		sm.l.Error(
			"iterSession command failed to create session manager",
			"err", err)
		return fmt.Errorf("failed to create session manager | %v", err)
	}
	defer f.Close()

	err = sm.iterSession(ID)
	if err != nil {
		return fmt.Errorf("failed session iteration | %v", err)
	}
	return nil
}

func iterSessionFunc(ctx * cli.Context) error {
	return iterSession(ctx.Args().Get(0))
}

var IterSessionCmd = cli.Command {
	Name: "iter",
	Usage: "iters over a session: usage -> ...session iter <session-id>",
	Action: iterSessionFunc,
}
