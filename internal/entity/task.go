package entity

import (
	"time"

	"github.com/lmnq/todolist/internal/errs"
)

type Task struct {
	Done     bool      `bson:"done" default:"false"`
	Title    string    `json:"title" bson:"title"`
	ActiveAt time.Time `json:"activeAt" bson:"activeAt"`
}

// Валидация записи
func (t Task) Validate() error {
	// валидация поля title
	switch title := t.Title; {
	case title == "":
		return errs.ErrEmptyTitle
	case len(title) > 200:
		return errs.ErrTooLongTitle
	}

	// валидация поля activeAt
	if t.ActiveAt.IsZero() {
		return errs.ErrEmptyDate
	}

	return nil
}
