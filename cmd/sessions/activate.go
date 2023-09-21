package sessions

import (
	"fmt"

	"github.com/mhernan88/dag-bisect/cmd"
	"github.com/urfave/cli/v2"
)

func (sm SessionManager) activateSession(ID string) error{
	sessionModel, err := sm.QuerySessionByID(ID)
	if err != nil {
		return fmt.Errorf("failed to query session by id | %v", err)
	}

	fmt.Printf("loading %s", sessionModel.State)
	ui, err := cmd.LoadState(sessionModel.State)
	if err != nil {
		return fmt.Errorf("failed to load state | %v", err)
	}

	increments, err := ui.CheckDAG(sm.l)
	if err != nil {
		return fmt.Errorf("failed to evaluate dag | %v", err)
	}

	err = sm.IncrementSessionSplits(ID, increments)
	if err != nil {
		return fmt.Errorf("failed to increment splits | %v", err)
	}

	return sm.cleanup(ID, sessionModel, ui)
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
