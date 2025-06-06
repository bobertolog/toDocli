package model

import (
	"errors"
	"strings"
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

func ParseStatus(status string) (StatusType, error) {
	switch strings.ToUpper(status) {
	case "TODO":
		return StatusTodo, nil
	case "IN_PROGRESS":
		return StatusInProgress, nil
	case "DONE":
		return StatusDone, nil
	default:
		return StatusTodo, errors.New("некорректный статус")
	}
}

type Entity interface {
	GetEntityID() int
}

type Task struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	StatusRaw   string     `json:"status"` // Для JSON: приходит как строка
	status      StatusType // Внутреннее значение
	CreatedAt   time.Time  `json:"created_at"`
}

func NewTask(id int, title, description, status string) (*Task, error) {
	t := &Task{
		ID:          id,
		Title:       title,
		Description: description,
		StatusRaw:   status,
		CreatedAt:   time.Now(),
	}
	err := t.NormalizeStatus()
	return t, err
}

func (t *Task) NormalizeStatus() error {
	status, err := ParseStatus(t.StatusRaw)
	if err != nil {
		return err
	}
	t.status = status
	return nil
}

func (t *Task) Status() string {
	return t.status.String()
}

func (t *Task) StatusType() StatusType {
	return t.status
}

func (t *Task) SetStatusType(s StatusType) {
	t.status = s
	t.StatusRaw = s.String()
}

func (t *Task) GetEntityID() int {
	return t.ID
}
