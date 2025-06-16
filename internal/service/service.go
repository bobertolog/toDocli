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

// Генерация ID через репозиторий
func GenerateTaskID() int {
	mu.Lock()
	defer mu.Unlock()
	return repository.GenerateTaskID()
}

// Добавить задачу
func AddTask(task *model.Task) {
	repository.Save(task)
}

// Получить все задачи
func GetAllTasks() []*model.Task {
	return repository.GetAllTasks()
}

// Найти задачу по ID
func FindTaskByID(id int) *model.Task {
	return repository.FindTaskByID(id)
}

// Удалить задачу по ID
func DeleteTask(id int) error {
	return repository.DeleteTask(id)
}

// Сохранить все задачи (вызов SaveAll из repository)
func SaveAllTasks() error {
	return repository.SaveAll()
}

// Генератор задач: каждые interval секунд отправляет задачу в канал
func StartTaskGenerator(ctx context.Context, interval time.Duration, out chan<- *model.Task) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Генерация задач остановлена.")
			return
		case <-ticker.C:
			task, err := model.NewTask(GenerateTaskID(), "Авто-задача", "Сгенерировано системой", "TODO")
			if err != nil {
				fmt.Println("Ошибка генерации задачи:", err)
				continue
			}
			out <- task
		}
	}
}

// Получает задачи из канала и сохраняет
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

// Логгер задач: отслеживает новые задачи и выводит в консоль
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
			tasks := GetAllTasks()
			for _, t := range tasks {
				id := t.GetEntityID()
				if !seen[id] && !repository.WasTaskLoaded(id) {
					fmt.Printf("[LOG] Новая задача: ID=%d, Название=%s\n", id, t.Title)
					seen[id] = true
				}
			}
		}
	}
}
