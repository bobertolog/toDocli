package service

import (
	"context"
	"fmt"
	"sync"
	"time"
	"todocli/internal/model"
	"todocli/internal/repository"
)

var mu sync.Mutex

func GenerateTaskID() int {
	return repository.GenerateTaskID()
}

func AddTask(t *model.Task) {
	repository.Save(t)
}

func GetAllTasks() []*model.Task {
	return repository.GetAllTasks()
}

func FindTaskByID(id int) *model.Task {
	return repository.FindTaskByID(id)
}

func DeleteTask(id int) error {
	return repository.DeleteTask(id)
}

func SaveAllTasks() error {
	return repository.SaveAll()
}

func StartTaskGenerator(ctx context.Context, interval time.Duration, out chan<- *model.Task) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Генерация задач остановлена.")
			return
		case <-ticker.C:
			t, err := model.NewTask(GenerateTaskID(), "Авто-задача", "Сгенерировано системой", "TODO")
			if err == nil {
				out <- t
			}
		}
	}
}

func TaskSaver(ctx context.Context, in <-chan *model.Task) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Сохранение задач завершено.")
			return
		case t := <-in:
			AddTask(t)
			fmt.Println("Задача сохранена:", t.Title)
		}
	}
}

func StartLogger(ctx context.Context) {
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()
	seen := make(map[int]bool)
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Логгер завершён.")
			return
		case <-ticker.C:
			for _, t := range GetAllTasks() {
				if !seen[t.ID] && !repository.WasTaskLoaded(t.ID) {
					fmt.Printf("[LOG] Новая задача: ID=%d, Title=%s\n", t.ID, t.Title)
					seen[t.ID] = true
				}
			}
		}
	}
}
