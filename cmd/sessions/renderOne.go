package sessions

import (
	"fmt"

	"github.com/jedib0t/go-pretty/v6/list"
	"github.com/mhernan88/dag-bisect/db/models"
	"github.com/mhernan88/dag-bisect/shared"
)

func (sm SessionManager) RenderSingleSessionTree(
	session models.Session,
) error {
	l := list.NewWriter()
	lTemp := list.List{}
	lTemp.Render()

	updatedUnixTimestamp, err := session.PrettyUpdated()
	if err != nil {
		return err
	}
	createdUnixTimestamp, err := session.PrettyCreated()
	if err != nil {
		return err
	}
	
	l.AppendItem(fmt.Sprintf("Splits: %d", session.Splits))
	if session.ErrNode != nil {
		l.AppendItem(fmt.Sprintf("Err Node: %s", *session.ErrNode))
	}
	l.AppendItem(fmt.Sprintf("Original DAG File: %s", session.DAG))
	l.AppendItem(fmt.Sprintf("State File: %s", session.State))
	l.AppendItem(fmt.Sprintf("Updated: %s", updatedUnixTimestamp))
	l.AppendItem(fmt.Sprintf("Created: %s", createdUnixTimestamp))

	l.SetStyle(list.StyleConnectedRounded)

	statusEmoji := models.SESSION_LOGOS[session.Status]
	shared.PrintTree(fmt.Sprintf("%v Session %s", statusEmoji, session.ID), l.Render())
	return nil
}
