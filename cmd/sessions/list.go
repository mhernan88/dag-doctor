package sessions

import (
	"fmt"
	"slices"

	"github.com/jedib0t/go-pretty/v6/list"
	"github.com/mhernan88/dag-bisect/db/models"
	"github.com/mhernan88/dag-bisect/shared"
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

// Renders a part of a tree for a given status (and associated sessions).
func (sm SessionManager) RenderTreeBranch(
	l list.Writer, 
	sessions []models.Session,
	status string,
) (list.Writer, error) {
	statusEmoji := models.SESSION_LOGOS[status]
	l.AppendItem(fmt.Sprintf("%v %s", statusEmoji, status))
	l.Indent()

	sm.l.Debug("rendering sessions", "status", status, "n", len(sessions))
	slices.SortFunc(sessions, models.SessionUpdateSortFunc)
	for _, session := range(sessions) {
		updatedUnixTimestamp, err := session.PrettyUpdated()
		if err != nil {
			return nil, err
		}
		createdUnixTimestamp, err := session.PrettyCreated()
		if err != nil {
			return nil, err
		}

		l.AppendItem(fmt.Sprintf("Session %s", session.ID))
		l.Indent()
		l.AppendItem(fmt.Sprintf("Splits: %d", session.Splits))
		l.AppendItem(fmt.Sprintf("Original DAG File: %s", session.DAG))
		l.AppendItem(fmt.Sprintf("State File: %s", session.State))
		l.AppendItem(fmt.Sprintf("Updated: %s", updatedUnixTimestamp))
		l.AppendItem(fmt.Sprintf("Created: %s", createdUnixTimestamp))
		l.UnIndent()
	}

	l.UnIndent()
	return l, nil
}

// Renders a full tree for a map of different statuses / sessions.
func (sm SessionManager) RenderTree(
	allSessions map[string][]models.Session,
) error {
	l := list.NewWriter()
	lTemp := list.List{}
	lTemp.Render()

	var err error
	for _, status := range models.SESSION_STATUSES {
		sessionGroup, ok := allSessions[status]
		if !ok {
			continue
		}

		l, err = sm.RenderTreeBranch(l, sessionGroup, status)
		if err != nil {
			return err
		}
	}

	l.SetStyle(list.StyleConnectedRounded)
	shared.PrintTree("Sessions", l.Render())
	return nil
}

// Wrapper of multiple SessionManager methods to list all sessions grouped
// by status.
func listSessions() error {
	sm, f, err := NewDefaultSessionManager()
	if err != nil {
		sm.l.Error(
			"listSessins command failed to create session manager",
			"err", err)
		return fmt.Errorf("failed to create session manager | %v", err)
	}
	defer f.Close()

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
	return listSessions()
}

var ListSessionsCmd = cli.Command {
	Name: "list",
	Aliases: []string{"ls"},
	Usage: "list all sessions: usage -> ...session ls",
	Action: listSessionsFunc,
}
