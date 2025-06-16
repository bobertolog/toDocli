package model

import (
	"errors"
	"strings"
	"time"
)

// StatusType ...
type StatusType int

const (
	StatusTodo StatusType = iota
	StatusInProgress
	StatusDone
)

func (s StatusType) String() string {
	return [...]string{"TODO", "IN_PROGRESS", "DONE"}[s]
}

func ParseStatus(s string) (StatusType, error) {
	switch strings.ToUpper(s) {
	case "TODO":
		return StatusTodo, nil
	case "IN_PROGRESS":
		return StatusInProgress, nil
	case "DONE":
		return StatusDone, nil
	default:
		return StatusTodo, errors.New("invalid status")
	}
}

// Task ...
type Task struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	StatusRaw   string     `json:"status"` // From JSON
	Status      StatusType `json:"-"`
	CreatedAt   time.Time  `json:"created_at"`
}

func (t *Task) Normalize() error {
	st, err := ParseStatus(t.StatusRaw)
	if err != nil {
		return err
	}
	t.Status = st
	t.StatusRaw = st.String()
	return nil
}

func NewTask(id int, title, description, status string) (*Task, error) {
	t := &Task{ID: id, Title: title, Description: description, StatusRaw: status, CreatedAt: time.Now()}
	if err := t.Normalize(); err != nil {
		return nil, err
	}
	return t, nil
}
