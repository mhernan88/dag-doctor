package sessions

import (
	"fmt"
	"slices"

	"github.com/jedib0t/go-pretty/v6/list"
	"github.com/mhernan88/dag-bisect/models"
	"github.com/mhernan88/dag-bisect/shared"
)

// Renders a part of a tree for a given status (and associated sessions).
func (sm SessionManager) RenderTreeBranch(
	l list.Writer, 
	sessions []models.Session,
	status string,
) (list.Writer, error) {
	statusEmoji := models.SESSION_LOGOS[status]
	l.AppendItem(fmt.Sprintf("%v %s", statusEmoji, status))
	l.Indent()

	sm.l.Info("rendering sessions", "status", status, "n", len(sessions))
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

		if session.ErrNode != nil {
			l.AppendItem(fmt.Sprintf("Err Node: %s", *session.ErrNode))
		}
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
