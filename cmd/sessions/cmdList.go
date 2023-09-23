package sessions

import (
	"fmt"

	"github.com/mhernan88/dag-bisect/db/models"
	"github.com/urfave/cli/v2"
)

// Queries sessions from sqlite db for a list of statuses.
func (sm SessionManager) QuerySessions(
	statuses []string,
) (map[string][]models.Session, error) {
	output := make(map[string][]models.Session)
	for _, status := range statuses {
		sessions, err := sm.QuerySessionsByStatus(status)
		if err != nil {
			fmt.Printf("failed to query sessions (where status='%s') from local database | %v\n", status, err)
			sm.l.Error("failed to query sessions from local database", "err", err)
		}

		if len(sessions) == 0 {
			continue
		}

		output[status] = sessions
	}
	return output, nil
}

// Wrapper of multiple SessionManager methods to list all sessions grouped
// by status.
func listSessions(sessionID string) error {
	sm, f, err := NewDefaultSessionManager()
	if err != nil {
		sm.l.Error(
			"listSessins command failed to create session manager",
			"err", err)
		return fmt.Errorf("failed to create session manager | %v", err)
	}
	defer f.Close()

	if sessionID != "" {
		session, err := sm.QuerySessionByID(sessionID)
		if err != nil {
			return fmt.Errorf("listSessions command failed to query session %s", sessionID)
		}
		err = sm.RenderSingleSessionTree(*session)
		if err != nil {
			return fmt.Errorf("listSessions command failed to render tree for session %s", sessionID)
		}
		return nil
	}

	allSessions, err := sm.QuerySessions(models.SESSION_STATUSES)
	if err != nil {
		sm.l.Error(
			"listSessions command failed to query sessions",
			"err", err)
		return fmt.Errorf("failed to query sessions | %v", err)
	}

	return sm.RenderTree(allSessions)
}

// Wrapper of listSessions() for urfave.
func listSessionsFunc(ctx *cli.Context) error {
	return listSessions(ctx.Args().Get(0))
}

var ListSessionsCmd = cli.Command {
	Name: "list",
	Aliases: []string{"ls"},
	Usage: "list all sessions: usage -> ...session ls",
	Action: listSessionsFunc,
}
