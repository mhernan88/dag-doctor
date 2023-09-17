package models

type Session struct {
	ID string `db:"id"`
	Status string `db:"status"`
	MetaCreatedDatetime int64 `db:"meta_created_datetime"`
	MetaUpdatedDatetime int64 `db:"meta_updated_datetime"`
}
