package repository

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"todocli/internal/model"
)

var (
	taskStore     = make(map[model.StatusType][]*model.Task)
	loadedTaskIDs = make(map[int]bool)
	mu            sync.RWMutex
	dataDir       = "data"
	idFile        = filepath.Join(dataDir, "last_id.txt")
	idCounter     = 0
)

// Инициализация при запуске
func init() {
	_ = os.MkdirAll(dataDir, 0755)
	loadAllTasksFromFiles()
	loadLastID()
}

// Получить уникальный ID
func GenerateTaskID() int {
	mu.Lock()
	defer mu.Unlock()
	idCounter++
	saveLastID()
	return idCounter
}

// Сохранить задачу (новую или обновлённую)
func Save(entity model.Entity) error {
	mu.Lock()
	defer mu.Unlock()

	task, ok := entity.(*model.Task)
	if !ok {
		return errors.New("неизвестный тип сущности")
	}

	status := task.StatusType()
	// Проверка — если такая задача уже есть, обновим её
	for i, t := range taskStore[status] {
		if t.GetEntityID() == task.GetEntityID() {
			taskStore[status][i] = task
			return saveToFile(status)
		}
	}
	// Если новой ID — добавим
	taskStore[status] = append(taskStore[status], task)
	return saveToFile(status)
}

// Удалить задачу по ID
func DeleteTask(id int) error {
	mu.Lock()
	defer mu.Unlock()

	for status, tasks := range taskStore {
		for i, t := range tasks {
			if t.GetEntityID() == id {
				taskStore[status] = append(tasks[:i], tasks[i+1:]...)
				delete(loadedTaskIDs, id)
				return saveToFile(status)
			}
		}
	}
	return errors.New("задача не найдена")
}

// Получить все задачи
func GetAllTasks() []*model.Task {
	mu.RLock()
	defer mu.RUnlock()

	var all []*model.Task
	for _, tasks := range taskStore {
		all = append(all, tasks...)
	}
	return all
}

// Найти задачу по ID
func FindTaskByID(id int) *model.Task {
	mu.RLock()
	defer mu.RUnlock()

	for _, tasks := range taskStore {
		for _, task := range tasks {
			if task.GetEntityID() == id {
				return task
			}
		}
	}
	return nil
}

// Проверить, была ли задача загружена из файла
func WasTaskLoaded(id int) bool {
	mu.RLock()
	defer mu.RUnlock()
	return loadedTaskIDs[id]
}

// Сохранить все задачи (по статусам)
func SaveAll() error {
	mu.RLock()
	defer mu.RUnlock()

	for status := range taskStore {
		if err := saveToFile(status); err != nil {
			return err
		}
	}
	return nil
}

// ---------- Внутренние функции ----------

func filenameForStatus(status model.StatusType) string {
	switch status {
	case model.StatusTodo:
		return filepath.Join(dataDir, "todo_tasks.json")
	case model.StatusInProgress:
		return filepath.Join(dataDir, "inprogress_tasks.json")
	case model.StatusDone:
		return filepath.Join(dataDir, "done_tasks.json")
	default:
		return filepath.Join(dataDir, "unknown_tasks.json")
	}
}

// Сохранение задач одного статуса в файл
func saveToFile(status model.StatusType) error {
	file := filenameForStatus(status)
	data, err := json.MarshalIndent(taskStore[status], "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(file, data, 0644)
}

// Загрузка задач одного статуса из файла
func loadFromFile(status model.StatusType) ([]*model.Task, error) {
	file := filenameForStatus(status)
	data, err := os.ReadFile(file)
	if err != nil {
		if os.IsNotExist(err) {
			return []*model.Task{}, nil
		}
		return nil, err
	}
	var tasks []*model.Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

// Загрузка всех задач
func loadAllTasksFromFiles() {
	for _, status := range []model.StatusType{model.StatusTodo, model.StatusInProgress, model.StatusDone} {
		tasks, err := loadFromFile(status)
		if err == nil {
			taskStore[status] = tasks
			for _, task := range tasks {
				loadedTaskIDs[task.GetEntityID()] = true
			}
		}
	}
}

// Сохранение текущего значения ID
func saveLastID() {
	_ = os.WriteFile(idFile, []byte(strconv.Itoa(idCounter)), 0644)
}

// Загрузка ID при старте
func loadLastID() {
	data, err := os.ReadFile(idFile)
	if err == nil {
		val, err := strconv.Atoi(string(data))
		if err == nil {
			idCounter = val
		}
	}
}
