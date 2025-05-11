package repository

import (
	"errors"
	"sync"

	"todocli/internal/model"
)

var (
	mu    sync.RWMutex
	tasks []*model.Task
)

func Save(task *model.Task) {
	mu.Lock()
	defer mu.Unlock()
	tasks = append(tasks, task)
}

func GetAllTasks() []*model.Task {
	mu.RLock()
	defer mu.RUnlock()
	return append([]*model.Task(nil), tasks...) // безопасная копия
}

func DeleteTask(id int) error {
	mu.Lock()
	defer mu.Unlock()
	for i, task := range tasks {
		if task.GetEntityID() == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			return nil
		}
	}
	return errors.New("task not found")
}
