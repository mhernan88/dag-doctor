package sessions

import (
	"fmt"

	"github.com/mhernan88/dag-bisect/cmd"
	"github.com/mhernan88/dag-bisect/models"
	"github.com/mhernan88/dag-bisect/pruners"
	"github.com/mhernan88/dag-bisect/splitters"
	"github.com/urfave/cli/v2"
)


func (sm SessionManager) iterSession(ID string) error{
	sessionModel, err := sm.QuerySessionByID(ID)
	if err != nil {
		return fmt.Errorf("failed to query session by id | %v", err)
	}

	fmt.Printf("loading %s", sessionModel.State)
	state, err := models.LoadState(sessionModel.State)
	if err != nil {
		return fmt.Errorf("failed to load state | %v", err)
	}

	pruner := pruners.NewDefaultPruner()
	splitter := splitters.NewDefaultSplitter()

	_, err = cmd.CheckDAGIter(state, pruner, splitter, sm.l)
	if err != nil {
		return fmt.Errorf("failed to evaluate dag | %v", err)
	}

	err = sm.IncrementSessionSplits(ID, 1)
	if err != nil {
		return fmt.Errorf("failed to increment splits | %v", err)
	}

	return sm.cleanup(ID, sessionModel, state)
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
