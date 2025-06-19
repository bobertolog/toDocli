// internal/service/service.go

package service

import (
	"todocli/internal/model"
)

// TaskRepository описывает зависимости от хранилища
type TaskRepository interface {
	Save(task *model.Task) error
	Update(task *model.Task) error
	GetAll() []*model.Task
	FindByID(id int) *model.Task
	Delete(id int) error
}

// TaskService — интерфейс бизнес-логики
type TaskService interface {
	Create(title, desc, status string) (*model.Task, error)
	GetAll() []*model.Task
	GetByID(id int) (*model.Task, error)
	Update(id int, title, desc, status string) error
	Delete(id int) error
}
