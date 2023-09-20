package sessions

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mhernan88/dag-bisect/cmd"
	"github.com/mhernan88/dag-bisect/data"
	"github.com/mhernan88/dag-bisect/shared"
	"github.com/urfave/cli/v2"
)

func (sm SessionManager) insertSession(id, savedDagFilename, savedSessionFilename string) error {
	dt := time.Now().Unix()
	_, err := sm.cxn.Exec(
		`INSERT INTO sessions (
			id, 
			dag,
			state,
			splits,
			status, 
			meta_created_datetime, 
			meta_updated_datetime
		) VALUES (?, ?, ?, 0, ?, ?, ?)`, 
		id, savedDagFilename, savedSessionFilename, "new", dt, dt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert new session | %v", err)
	}
	return nil
}

func newSession(dagFilename string) error {
	sm, f, err := NewDefaultSessionManager()
	if err != nil {
		sm.l.Error(
			"listSessins command failed to create session manager",
			"err", err)
		return fmt.Errorf("failed to create session manager | %v", err)
	}
	defer f.Close()

	dag, err := data.LoadDAG(dagFilename)
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

	err = sm.insertSession(id, savedDagFilename, sessionFilename)
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
