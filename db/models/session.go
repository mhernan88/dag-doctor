package models

import (
	"fmt"
	"strconv"
	"time"
)

var SESSION_STATUSES = []string{
	"new", "in-progress", "ok", "err",
}

func SessionUpdateSortFunc(a Session, b Session) int {
	if a.MetaUpdatedDatetime > b.MetaUpdatedDatetime {
		return 1
	} else if a.MetaUpdatedDatetime < b.MetaUpdatedDatetime {
		return -1
	} else {
		return 0
	}
}

type Session struct {
	ID string `db:"id"`
	File string `db:"file"`
	Splits int `db:"splits"`
	Status string `db:"status"`
	MetaCreatedDatetime int64 `db:"meta_created_datetime"`
	MetaUpdatedDatetime int64 `db:"meta_updated_datetime"`
}

func (s Session) PrettyUpdated() (time.Time, error) {
		updatedUnixTimestamp := fmt.Sprintf("%d", s.MetaUpdatedDatetime)
		i, err := strconv.ParseInt(updatedUnixTimestamp, 10, 64)
		if err != nil {
			return time.Time{}, err
		}
		return time.Unix(i, 0), nil

}

func (s Session) PrettyCreated() (time.Time, error) {
		updatedUnixTimestamp := fmt.Sprintf("%d", s.MetaCreatedDatetime)
		i, err := strconv.ParseInt(updatedUnixTimestamp, 10, 64)
		if err != nil {
			return time.Time{}, err
		}
		return time.Unix(i, 0), nil

}
