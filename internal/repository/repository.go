package repository

import (
	"errors"
	"todocli/internal/model"
)

var (
	tasks []*model.Task
)

func Save(entity model.Entity) error {
	switch e := entity.(type) {
	case *model.Task:
		tasks = append(tasks, e)
		return nil
	default:
		return errors.New("unknown entity type")
	}
}

func GetAllTasks() []*model.Task {
	return tasks
}

func DeleteTask(id int) error {
	for i, task := range tasks {
		if task.GetEntityID() == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			return nil
		}
	}
	return errors.New("task not found")
}
