package service

import (
	"log"
	"sync"
	"time"
	"todocli/internal/model"
	"todocli/internal/repository"
)

var (
	mu          sync.Mutex
	lastTaskID  int
	prevTaskLen int
)

// Генерация ID
func GenerateTaskID() int {
	mu.Lock()
	defer mu.Unlock()
	lastTaskID++
	return lastTaskID
}

// Добавление задачи вручную
func AddTask(task *model.Task) {
	repository.Save(task)
}

// Получение всех задач
func GetAllTasks() []*model.Task {
	return repository.GetAllTasks()
}

// Поиск задачи по ID
func FindTaskByID(id int) *model.Task {
	tasks := repository.GetAllTasks()
	for _, task := range tasks {
		if task.GetEntityID() == id {
			return task
		}
	}
	return nil
}

// Удаление задачи по ID
func DeleteTask(id int) error {
	return repository.DeleteTask(id)
}

// генератор задач
func StartTaskGenerator(interval time.Duration, taskChan chan *model.Task) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		task := model.NewTask(GenerateTaskID(), "Авто-задача", "Создана автоматически")
		taskChan <- task
	}
}

func TaskSaver(taskChan chan *model.Task) {
	for task := range taskChan {
		AddTask(task)
	}
}

func StartLogger(done <-chan struct{}) {
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	var lastCount int
	var lastSeen map[int]bool = make(map[int]bool)

	for {
		select {
		case <-done:
			log.Println("Логгер завершил работу.")
			return
		case <-ticker.C:
			currentTasks := GetAllTasks()

			if len(currentTasks) != lastCount {
				log.Println("Изменения в задачах:")

				for _, task := range currentTasks {
					if !lastSeen[task.GetEntityID()] {
						log.Printf("Новая задача: ID=%d, Title=%s, Status=%s\n",
							task.GetEntityID(), task.Title, task.Status())
						lastSeen[task.GetEntityID()] = true
					}
				}

				lastCount = len(currentTasks)
			}
		}
	}
}
