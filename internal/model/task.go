package model

import (
	"errors"
	"time"
)

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
	switch s {
	case "TODO":
		return StatusTodo, nil
	case "IN_PROGRESS":
		return StatusInProgress, nil
	case "DONE":
		return StatusDone, nil
	default:
		return 0, errors.New("invalid status")
	}
}

type Task struct {
	ID          int
	Title       string
	Description string
	Status      StatusType
	CreatedAt   time.Time
}

func NewTask(id int, title, description, status string) (*Task, error) {
	s, err := ParseStatus(status)
	if err != nil {
		return nil, err
	}
	return &Task{
		ID:          id,
		Title:       title,
		Description: description,
		Status:      s,
		CreatedAt:   time.Now(),
	}, nil
}
