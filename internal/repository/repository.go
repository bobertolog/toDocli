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

func Init() error {
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return err
	}
	loadAllTasksFromFiles()
	loadLastID()
	return nil
}

func GenerateTaskID() int {
	mu.Lock()
	defer mu.Unlock()
	idCounter++
	_ = os.WriteFile(idFile, []byte(strconv.Itoa(idCounter)), 0644)
	return idCounter
}

func Save(entity *model.Task) error {
	mu.Lock()
	defer mu.Unlock()
	status := entity.Status
	for i, t := range taskStore[status] {
		if t.ID == entity.ID {
			taskStore[status][i] = entity
			return saveToFile(status)
		}
	}
	taskStore[status] = append(taskStore[status], entity)
	return saveToFile(status)
}

func DeleteTask(id int) error {
	mu.Lock()
	defer mu.Unlock()
	for status, list := range taskStore {
		for i, t := range list {
			if t.ID == id {
				taskStore[status] = append(list[:i], list[i+1:]...)
				delete(loadedTaskIDs, id)
				return saveToFile(status)
			}
		}
	}
	return errors.New("задача не найдена")
}

func GetAllTasks() []*model.Task {
	mu.RLock()
	defer mu.RUnlock()
	var all []*model.Task
	for _, list := range taskStore {
		all = append(all, list...)
	}
	return all
}

func FindTaskByID(id int) *model.Task {
	mu.RLock()
	defer mu.RUnlock()
	for _, list := range taskStore {
		for _, t := range list {
			if t.ID == id {
				return t
			}
		}
	}
	return nil
}

func WasTaskLoaded(id int) bool {
	mu.RLock()
	defer mu.RUnlock()
	return loadedTaskIDs[id]
}

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

func filenameForStatus(status model.StatusType) string {
	switch status {
	case model.StatusTodo:
		return filepath.Join(dataDir, "todo_tasks.json")
	case model.StatusInProgress:
		return filepath.Join(dataDir, "inprogress_tasks.json")
	case model.StatusDone:
		return filepath.Join(dataDir, "done_tasks.json")
	default:
		return filepath.Join(dataDir, "unknown.json")
	}
}

func saveToFile(status model.StatusType) error {
	file := filenameForStatus(status)
	data, err := json.MarshalIndent(taskStore[status], "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(file, data, 0644)
}

func loadFromFile(status model.StatusType) []*model.Task {
	file := filenameForStatus(status)
	data, err := os.ReadFile(file)
	if err != nil {
		return []*model.Task{}
	}
	var list []*model.Task
	_ = json.Unmarshal(data, &list)
	return list
}

func loadAllTasksFromFiles() {
	for _, status := range []model.StatusType{model.StatusTodo, model.StatusInProgress, model.StatusDone} {
		list := loadFromFile(status)
		taskStore[status] = list
		for _, t := range list {
			loadedTaskIDs[t.ID] = true
		}
	}
}

func loadLastID() {
	data, err := os.ReadFile(idFile)
	if err != nil {
		return
	}
	if v, err := strconv.Atoi(string(data)); err == nil {
		idCounter = v
	}
}
