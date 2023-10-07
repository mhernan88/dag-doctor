package sessions

import (
	"fmt"

	"github.com/mhernan88/dag-bisect/resolver/interactive"
	"github.com/mhernan88/dag-bisect/models"
	"github.com/mhernan88/dag-bisect/resolver/pruners"
	"github.com/mhernan88/dag-bisect/resolver/splitters"
	"github.com/urfave/cli/v2"
)

func (sm SessionManager) activateSession(ID string) error{
	ID, err := sm.QuerySessionIDByPartialID(ID)
	if err != nil {
		return fmt.Errorf("failed to enrich session id | %v", err)
	}
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

	increments, err := interactive.CheckDAG(state, pruner, splitter, sm.l)
	if err != nil {
		return fmt.Errorf("failed to evaluate dag | %v", err)
	}

	err = sm.IncrementSessionSplits(ID, increments)
	if err != nil {
		return fmt.Errorf("failed to increment splits | %v", err)
	}

	return sm.cleanup(ID, sessionModel, state)
}

func activateSession(ID string) error {
	sm, f, err := NewDefaultSessionManager()
	if err != nil {
		sm.l.Error(
			"iterSession command failed to create session manager",
			"err", err)
		return fmt.Errorf("failed to create session manager | %v", err)
	}
	defer f.Close()

	err = sm.activateSession(ID)
	if err != nil {
		return fmt.Errorf("failed session iterations | %v", err)
	}
	return nil
}

func activateSessionFunc(ctx *cli.Context) error {
	return activateSession(ctx.Args().Get(0))
}


var ActivateSessionCmd = cli.Command {
	Name: "activate",
	Usage: "activates a session: usage -> ...session activate <session-id>",
	Action: activateSessionFunc,
}
