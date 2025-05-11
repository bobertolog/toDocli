package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"todocli/internal/model"
	"todocli/internal/repository"
)

var idCounter = 0
var mu sync.Mutex

func GenerateTaskID() int {
	mu.Lock()
	defer mu.Unlock()
	idCounter++
	return idCounter
}

func AddTask(task *model.Task) {
	repository.Save(task)
}

func GetAllTasks() []*model.Task {
	return repository.GetAllTasks()
}

func FindTaskByID(id int) *model.Task {
	tasks := repository.GetAllTasks()
	for _, task := range tasks {
		if task.GetEntityID() == id {
			return task
		}
	}
	return nil
}

func DeleteTask(id int) error {
	return repository.DeleteTask(id)
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
			task := model.NewTask(GenerateTaskID(), " Авто-задача", "Сгенерирована системой")
			out <- task
		}
	}
}

func TaskSaver(ctx context.Context, in <-chan *model.Task) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Сохранение задач завершено.")
			return
		case task := <-in:
			AddTask(task)
			fmt.Println("Задача сохранена:", task.Title)
		}
	}
}

var previousCounts = make(map[string]int)

func StartLogger(ctx context.Context) {
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	prevMap := make(map[int]bool)

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Логгер завершён.")
			return
		case <-ticker.C:
			tasks := GetAllTasks()
			currentMap := make(map[int]bool)
			for _, t := range tasks {
				currentMap[t.GetEntityID()] = true
			}
			for _, t := range tasks {
				if !prevMap[t.GetEntityID()] {
					fmt.Printf("[LOG] Новая задача: ID=%d, Название=%s\n", t.GetEntityID(), t.Title)
				}
			}
			prevMap = currentMap
		}
	}
}
