package service

import (
	"errors"
	"fmt"
	"sync"
	"todocli/internal/model"
)

type taskService struct {
	repo TaskRepository
}

func NewTaskService(repo TaskRepository) TaskService {
	return &taskService{repo: repo}
}

func (s *taskService) Create(title, desc, status string) (*model.Task, error) {
	t, err := model.NewTask(GenerateTaskID(), title, desc, status)
	if err != nil {
		return nil, err
	}
	err = s.repo.Save(t)
	return t, err
}

func (s *taskService) GetAll() []*model.Task {
	return s.repo.GetAll()
}

func (s *taskService) GetByID(id int) (*model.Task, error) {
	t := s.repo.FindByID(id)
	if t == nil {
		return nil, fmt.Errorf("task with id=%d not found", id)
	}
	return t, nil
}

func (s *taskService) Update(id int, title, desc, status string) error {
	t := s.repo.FindByID(id)
	if t == nil {
		return errors.New("task not found")
	}
	t.Title = title
	t.Description = desc

	st, err := model.ParseStatus(status) // ❗ Используй другое имя переменной!
	if err != nil {
		return err
	}
	t.Status = st
	return s.repo.Update(t)
}

func (s *taskService) Delete(id int) error {
	return s.repo.Delete(id)
}

var (
	idCounter int
	mu        sync.Mutex
)

func GenerateTaskID() int {
	mu.Lock()
	defer mu.Unlock()
	idCounter++
	return idCounter
}
