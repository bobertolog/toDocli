package repository

import (
	"errors"
	"sync"
	"todocli/internal/model"
)

var (
	mu    sync.RWMutex
	tasks = make([]*model.Task, 0)
)

// Save сохраняет любую сущность типа Task безопасно для горутин
func Save(entity model.Entity) error {
	mu.Lock()
	defer mu.Unlock()

	switch e := entity.(type) {
	case *model.Task:
		tasks = append(tasks, e)
		return nil
	default:
		return errors.New("unknown entity type")
	}
}

// GetAllTasks возвращает копию всех задач безопасно для горутин
func GetAllTasks() []*model.Task {
	mu.RLock()
	defer mu.RUnlock()

	// Чтобы не отдавать оригинальный слайс (иначе другая горутина может его изменить)
	copyTasks := make([]*model.Task, len(tasks))
	copy(copyTasks, tasks)
	return copyTasks
}

// DeleteTask удаляет задачу по ID безопасно для горутин
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
