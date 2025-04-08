package task

import "time"

// Task — структура для задачи
type Task struct {
    id          int       
    Title       string    
    status      string    
    CreatedAt   time.Time 
    Description string    
}


// NewTask — конструктор для создания задачи
func NewTask(id int, title, status, description string) *Task {
    return &Task{
        id:          id,
        Title:       title,
        status:      status,
        CreatedAt:   time.Now(),
        Description: description,
    }
}

// GetID — метод для получения ID
func (t *Task) GetID() int {
    return t.id
}

// GetStatus — метод для получения статуса
func (t *Task) GetStatus() string {
    return t.status
}

// SetStatus — метод для изменения статуса
func (t *Task) SetStatus(newStatus string) {
    t.status = newStatus
}