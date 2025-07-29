package repository

import (
	"todocli/internal/model"
	"todocli/internal/service"
)

type InMemoryRepository struct {
	tasks map[int]*model.Task
}

var _ service.TaskRepository = (*InMemoryRepository)(nil)

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		tasks: make(map[int]*model.Task),
	}
}

func (r *InMemoryRepository) Save(task *model.Task) error {
	r.tasks[task.ID] = task
	return nil
}

func (r *InMemoryRepository) Update(task *model.Task) error {
	r.tasks[task.ID] = task
	return nil
}

func (r *InMemoryRepository) GetAll() []*model.Task {
	var all []*model.Task
	for _, t := range r.tasks {
		all = append(all, t)
	}
	return all
}

func (r *InMemoryRepository) FindByID(id int) *model.Task {
	return r.tasks[id]
}

func (r *InMemoryRepository) Delete(id int) error {
	delete(r.tasks, id)
	return nil
}

func (r *InMemoryRepository) WithTx(fn func(service.TaskRepository) error) error {
	return fn(r)
}
