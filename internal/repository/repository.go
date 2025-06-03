package repository

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"todocli/internal/model"
)

var (
	taskStore = make(map[model.StatusType][]*model.Task)
	mu        sync.RWMutex
	dataDir   = "data"
)

func init() {
	_ = os.MkdirAll(dataDir, 0755)
	loadAllTasksFromFiles()
}

func Save(entity model.Entity) error {
	mu.Lock()
	defer mu.Unlock()

	task, ok := entity.(*model.Task)
	if !ok {
		return errors.New("unknown entity type")
	}
	status := task.StatusType()
	taskStore[status] = append(taskStore[status], task)
	return saveToFile(status)
}

func GetAllTasks() []*model.Task {
	mu.RLock()
	defer mu.RUnlock()

	var all []*model.Task
	for _, tasks := range taskStore {
		all = append(all, tasks...)
	}
	return all
}

func DeleteTask(id int) error {
	mu.Lock()
	defer mu.Unlock()

	for status, tasks := range taskStore {
		for i, t := range tasks {
			if t.GetEntityID() == id {
				taskStore[status] = append(tasks[:i], tasks[i+1:]...)
				return saveToFile(status)
			}
		}
	}
	return errors.New("task not found")
}

// ---------- Persistence ----------

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

func saveToFile(status model.StatusType) error {
	file := filenameForStatus(status)
	data, err := json.MarshalIndent(taskStore[status], "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(file, data, 0644)
}

func loadFromFile(status model.StatusType) ([]*model.Task, error) {
	file := filenameForStatus(status)
	data, err := os.ReadFile(file)
	if err != nil {
		if os.IsNotExist(err) {
			return []*model.Task{}, nil // no data yet
		}
		return nil, err
	}
	var tasks []*model.Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

func loadAllTasksFromFiles() {
	for _, status := range []model.StatusType{model.StatusTodo, model.StatusInProgress, model.StatusDone} {
		tasks, err := loadFromFile(status)
		if err == nil {
			taskStore[status] = tasks
		}
	}
}
