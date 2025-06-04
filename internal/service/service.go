package service

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"todocli/internal/model"
	"todocli/internal/repository"
)

const idFilePath = "data/id_counter.txt"

var (
	idCounter int
	mu        sync.Mutex
)

func GenerateTaskID() int {
	mu.Lock()
	defer mu.Unlock()
	idCounter++
	_ = saveIDToFile()
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
func InitIDCounter() error {
	data, err := os.ReadFile(idFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			idCounter = 0
			return nil
		}
		return err
	}
	id, err := strconv.Atoi(string(data))
	if err != nil {
		return err
	}
	idCounter = id
	return nil
}

func saveIDToFile() error {
	return os.WriteFile(idFilePath, []byte(strconv.Itoa(idCounter)), 0644)
}
