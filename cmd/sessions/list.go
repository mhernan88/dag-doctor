package sessions

import (
	"fmt"
	"log/slog"
	"slices"

	"github.com/jedib0t/go-pretty/v6/list"
	"github.com/jmoiron/sqlx"
	"github.com/mhernan88/dag-bisect/db"
	"github.com/mhernan88/dag-bisect/db/models"
	"github.com/mhernan88/dag-bisect/shared"
	"github.com/urfave/cli/v2"
)

func NewListSessionsManager(cxn *sqlx.DB, l *slog.Logger) ListSessionsManager {
	return ListSessionsManager{
		cxn: cxn,
		l: l,
	}
}

type ListSessionsManager struct {
	cxn *sqlx.DB
	l *slog.Logger
}


func (lsm ListSessionsManager) QuerySessionsByStatus(
	status string,
) ( []models.Session, error ) {
	var sessions []models.Session

	query := fmt.Sprintf("SELECT * FROM sessions WHERE status = '%s'", status)
	lsm.l.Debug("executing select query", "table", "sessions", "query", query)

	err := lsm.cxn.Select(
		&sessions, query,
	)
	if err != nil {
		lsm.l.Error("failed select from sessions", "err", err)
		return nil, fmt.Errorf("failed select from sessions | %v", err)
	}
	return sessions, nil
}

func (lsm ListSessionsManager) QuerySessions(
	statuses []string,
) (map[string][]models.Session, error) {
	output := make(map[string][]models.Session)
	for _, status := range statuses {
		sessions, err := lsm.QuerySessionsByStatus(status)
		if err != nil {
			fmt.Printf("failed to query sessions (where status='%s') from local database | %v\n", status, err)
			lsm.l.Error("failed to query sessions from local database", "err", err)
		}

		if len(sessions) == 0 {
			continue
		}

		output[status] = sessions
	}
	return output, nil
}


func (lsm ListSessionsManager) RenderTreeBranch(
	l list.Writer, 
	sessions []models.Session,
	status string,
) (list.Writer, error) {
	statusEmoji := models.SESSION_LOGOS[status]
	l.AppendItem(fmt.Sprintf("%v %s", statusEmoji, status))
	l.Indent()

	lsm.l.Debug("rendering sessions", "status", status, "n", len(sessions))
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
		l.AppendItem(fmt.Sprintf("DAG File: %s", session.File))
		l.AppendItem(fmt.Sprintf("Updated: %s", updatedUnixTimestamp))
		l.AppendItem(fmt.Sprintf("Created: %s", createdUnixTimestamp))
		l.UnIndent()
	}

	l.UnIndent()
	return l, nil
}

func (lsm ListSessionsManager) RenderTree(
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

		l, err = lsm.RenderTreeBranch(l, sessionGroup, status)
		if err != nil {
			return err
		}
	}

	l.SetStyle(list.StyleConnectedRounded)
	shared.PrintTree("Sessions", l.Render())
	return nil
}

func listSessions(ctx *cli.Context) error {
	l, f := shared.GetLogger()
	defer f.Close()
	cxn, err := db.Connect()
	if err != nil {
		return err
	}

	lsm := NewListSessionsManager(cxn, l)
	allSessions, err := lsm.QuerySessions(models.SESSION_STATUSES)
	if err != nil {
		l.Error(
			"listSessions command failed to query sessions",
			"err", err)
		return fmt.Errorf("failed to query sessions | %v", err)
	}

	return lsm.RenderTree(allSessions)
}

var ListSessionsCmd = cli.Command {
	Name: "ls",
	Usage: "list all sessions: usage -> ...session ls",
	Action: listSessions,
}
