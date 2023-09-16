package models

type Session struct {
	ID string `db:"id"`
	Status string `db:"status"`
}
