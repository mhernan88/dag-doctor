package sessions

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/list"
	"github.com/jmoiron/sqlx"
	"github.com/mhernan88/dag-bisect/db"
	"github.com/mhernan88/dag-bisect/db/models"
	"github.com/urfave/cli/v2"
)

func demoPrint(title string, content string, prefix string) {
	fmt.Printf("%s:\n", title)
	fmt.Println(strings.Repeat("-", len(title)+1))
	for _, line := range strings.Split(content, "\n") {
		fmt.Printf("%s%s\n", prefix, line)
	}
	fmt.Println()
}

func requestAndRenderTreeBranch(
	l list.Writer, 
	status string, 
	dbHandle *sqlx.DB,
) (list.Writer, error) {
	var sessions []models.Session

	err := dbHandle.Select(
		&sessions,
		fmt.Sprintf("SELECT * FROM sessions WHERE status = '%s'", status),
	)
	if err != nil {
		return nil, err
	}

	if len(sessions) == 0 {
		fmt.Printf("%s : no sessions\n", status)
		return l, nil
	}

	fmt.Printf("%s : rendering for %d sessions\n", status, len(sessions))

	l.AppendItem(status)
	l.Indent()

	slices.SortFunc(sessions, func(a models.Session, b models.Session) int {
		if a.MetaUpdatedDatetime > b.MetaUpdatedDatetime {
			return 1
		} else if a.MetaUpdatedDatetime < b.MetaUpdatedDatetime {
			return -1
		} else {
			return 0
		}
	})

	for _, session := range(sessions) {
		updatedUnixTimestamp := fmt.Sprintf("%d", session.MetaUpdatedDatetime)
		i, err := strconv.ParseInt(updatedUnixTimestamp, 10, 64)
		if err != nil {
			return nil, err
		}
		tm := time.Unix(i, 0)

		createdUnixTimestamp := fmt.Sprintf("%d", session.MetaCreatedDatetime)
		j, err := strconv.ParseInt(createdUnixTimestamp, 10, 64)
		if err != nil {
			return nil, err
		}
		tm2 := time.Unix(j, 0)

		l.AppendItem(fmt.Sprintf("Session %s", session.ID))
		l.Indent()
		l.AppendItem(fmt.Sprintf("Updated: %s", tm))
		l.AppendItem(fmt.Sprintf("Created: %s", tm2))
		l.UnIndent()
	}

	l.UnIndent()

	return l, nil
}

func renderSessionsTree(dbHandle *sqlx.DB) error {
	l := list.NewWriter()
	lTemp := list.List{}
	lTemp.Render()

	// sessionIDs := make(map[string][]string)
	// for _, status := 

	var err error
	for _, status := range []string{"new", "in-progress", "ok", "err"} {
		l, err = requestAndRenderTreeBranch(l, status, dbHandle)
		if err != nil {
			return err
		}
	}


	// l.AppendItems([]interface{}{"Winter", "Is", "Coming"})
	// l.Indent()
	// l.AppendItems([]interface{}{"This", "Is", "Known"})
	// l.UnIndent()
	// l.UnIndent()
	// l.AppendItem("The Dark Tower")
	// l.Indent()
	// l.AppendItem("The Gunslinger")
	//
	l.SetStyle(list.StyleConnectedRounded)
	demoPrint("Sessions", l.Render(), "")
	return nil
}

func requestSessions(dbHandle *sqlx.DB, statusFilter string) ([]models.Session, error) {
	var err error
	var sessions []models.Session
	if statusFilter == "" {
		err = dbHandle.Select(
			&sessions,
			`SELECT * FROM sessions WHERE status != 'closed'`)
	} else if statusFilter == "all" {
		fmt.Println("filtering to all sessions")
		err = dbHandle.Select(
			&sessions,
			`SELECT * FROM sessions WHERE status = 'all'`,
		)
	} else {
		fmt.Printf("filtering to sessions with status = '%s'\n", statusFilter)
		err = dbHandle.Select(
			&sessions,
			`SELECT * FROM sessions WHERE status = '?'`,
			statusFilter)
	}
	if err != nil {
		return nil, err
	}
	return sessions, nil
}

func listSessions(ctx *cli.Context) error {
	dbHandle, err := db.Connect()
	if err != nil {
		return err
	}

	statusFilter := ctx.String("status")
	sessions, err := requestSessions(dbHandle, statusFilter)
	if err != nil {
		return err
	}

	err = renderSessionsTree(dbHandle)
	if err != nil {
		return err
	}

	fmt.Println(sessions)
	return nil
}

var listFlags = []cli.Flag{
	&cli.StringFlag{
		Name: "status",
		Value: "",
		Usage: "session status filter",
	},
}

var ListSessionsCmd = cli.Command {
	Name: "ls",
	Usage: "list sessions",
	Action: listSessions,
	Flags: listFlags, 
}
