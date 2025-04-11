package service

import (
	"time"
	"todocli/internal/model"
	"todocli/internal/repository"
)

// Генерация тасков
func StartTaskGenerator(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		task := model.NewTask(
			len(repository.GetAllTasks())+1,
			"Auto-generated task",
			"This task was created automatically",
		)
		_ = repository.Save(task)
	}
}
