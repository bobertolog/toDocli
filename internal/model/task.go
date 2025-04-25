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

type Entity interface {
	GetEntityID() int
}

type Task struct {
	id          int
	Title       string
	status      StatusType
	CreatedAt   time.Time
	Description string
}

func NewTask(id int, title, description string) *Task {
	return &Task{
		id:          id,
		Title:       title,
		status:      StatusTodo,
		CreatedAt:   time.Now(),
		Description: description,
	}
}

func (t *Task) GetEntityID() int {
	return t.id
}

func (t *Task) Status() string {
	return t.status.String()
}

func (t *Task) StatusType() StatusType {
	return t.status
}

func (t *Task) SetStatus(status string) error {
	switch status {
	case "TODO":
		t.status = StatusTodo
	case "IN_PROGRESS":
		t.status = StatusInProgress
	case "DONE":
		t.status = StatusDone
	default:
		return errors.New("invalid status value")
	}
	return nil
}

func (t *Task) SetStatusType(status StatusType) {
	t.status = status
}
