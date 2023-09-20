package sessions

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/mhernan88/dag-bisect/db"
	"github.com/mhernan88/dag-bisect/db/models"
)

// Queries sessions from sqlite db filtering to a given status.
func (sm SessionManager) QuerySessionsByStatus(
	status string,
) ([]models.Session, error) {
	var sessions []models.Session
	query := fmt.Sprintf("SELECT * FROM sessions WHERE status = '%s'", status)
	sm.l.Debug("executing select query", "table", "sessions", "query", query)

	err := sm.cxn.Select(
		&sessions, query,
	)
	if err != nil {
		sm.l.Error("failed select from sessions", "err", err)
		return nil, fmt.Errorf("failed select from sessions | %v", err)
	}
	return sessions, nil
}

func (sm SessionManager) QuerySessionByID(
	ID string,
) (*models.Session, error) {
	var sessions []models.Session
	query := fmt.Sprintf("SELECT * FROM sessions WHERE id = '%s'", ID)
	sm.l.Debug("executing select query", "table", "sessions", "query", query)

	err := sm.cxn.Select(
		&sessions, query,
	)
	if err != nil {
		sm.l.Error("failed select from sessions", "err", err)
		return nil, fmt.Errorf("failed select from sessions | %v", err)
	}

	if len(sessions) > 1 {
		sm.l.Error(
			"ambiguous select from sessions (expected 1)",
			"n_results", len(sessions))
	}
	return &sessions[0], nil
}

func (sm SessionManager) InsertSession(id, savedDagFilename, savedSessionFilename string) error {
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

func (sm SessionManager) UpdateSessionStatus(ID, status string) error {
	if !slices.Contains(models.SESSION_STATUSES, status) {
		sessionStatuses := strings.Join(models.SESSION_STATUSES, ", ")
		return fmt.Errorf("status must be one of: %s", sessionStatuses)
	}

	cxn, err := db.Connect()
	if err != nil {
		return err
	}
	defer cxn.Close()

	dt := time.Now().Unix()
	query := fmt.Sprintf(
		`UPDATE sessions
		SET 
			status = '%s', 
			meta_updated_datetime = '%d' 
		WHERE id = '%s'`,
		status, dt, ID,
	)

	_, err = cxn.Exec(query, status, ID)
	if err != nil {
		return fmt.Errorf("failed to update session | %v", err)
	}

	fmt.Printf("updated session '%s' to status='%s'\n", ID, status)
	return nil
}

func (sm SessionManager) IncrementSessionSplits(ID string) error {
	cxn, err := db.Connect()
	if err != nil {
		return err
	}
	defer cxn.Close()

	dt := time.Now().Unix()
	query := fmt.Sprintf(
		`UPDATE sessions
		SET 
			splits = splits + 1, 
			meta_updated_datetime = '%d' 
		WHERE id = '%s'`,
		dt, ID,
	)

	_, err = cxn.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to update session | %v", err)
	}

	fmt.Printf("incremented session '%s' splits\n", ID)
	return nil
}
