package sessions

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/mhernan88/dag-bisect/cmd"
	"github.com/mhernan88/dag-bisect/models"
	"github.com/mhernan88/dag-bisect/shared"
	"github.com/urfave/cli/v2"
)


func newSession(dagFilename string) error {
	sm, f, err := NewDefaultSessionManager()
	if err != nil {
		sm.l.Error(
			"listSessions command failed to create session manager",
			"err", err)
		return fmt.Errorf("failed to create session manager | %v", err)
	}
	defer f.Close()

	dag, err := models.LoadDAG(dagFilename)
	if err != nil {
		return fmt.Errorf("failed to load dag | %v", err)
	}

	ui := cmd.NewDefaultUI(*dag)
	sessionFilename, err := shared.SaveStateToRepo(ui)
	if err != nil {
		return fmt.Errorf("failed to save state | %v", err)
	}
	fmt.Printf("initialized session at %s\n", sessionFilename)

	id := uuid.NewString()
	savedDagFilename, err := shared.CopyDAGToRepo(dagFilename, id)
	if err != nil {
		return fmt.Errorf("failed to copy dag to repo | %v", err)
	}

	err = sm.InsertSession(id, savedDagFilename, sessionFilename)
	if err != nil {
		return fmt.Errorf("failed to insert into sessions table | %v", err)
	}

	fmt.Printf("created session %s\n", id)
	return nil
}

func newSessionFunc(ctx *cli.Context) error {
	return newSession(ctx.Args().Get(0))
}

var NewSessionCmd = cli.Command {
	Name: "new",
	Usage: "creates a new session: usage -> ...session new <dag-filename>",
	Action: newSessionFunc,
}
