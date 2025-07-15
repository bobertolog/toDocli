package service

import (
	"testing"
	"todocli/internal/model"

	"github.com/stretchr/testify/assert"
)

type mockRepo struct {
	tasks map[int]*model.Task
}

func newMockRepo() *mockRepo {
	return &mockRepo{tasks: make(map[int]*model.Task)}
}

func (m *mockRepo) Save(task *model.Task) error {
	m.tasks[task.ID] = task
	return nil
}

func (m *mockRepo) Update(task *model.Task) error {
	m.tasks[task.ID] = task
	return nil
}

func (m *mockRepo) GetAll() []*model.Task {
	var all []*model.Task
	for _, t := range m.tasks {
		all = append(all, t)
	}
	return all
}

func (m *mockRepo) FindByID(id int) *model.Task {
	return m.tasks[id]
}

func (m *mockRepo) Delete(id int) error {
	delete(m.tasks, id)
	return nil
}

// ✅ Заглушка для транзакций
func (m *mockRepo) WithTx(fn func(TaskRepository) error) error {
	return fn(m)
}

func TestTaskService_Create(t *testing.T) {
	repo := newMockRepo()
	svc := NewTaskService(repo)

	task, err := svc.Create("Test Title", "Test Desc", "TODO")

	assert.NoError(t, err)
	assert.Equal(t, "Test Title", task.Title)
	assert.Equal(t, "TODO", task.Status.String())
	assert.Equal(t, "Test Desc", task.Description)
}

func TestTaskService_Update(t *testing.T) {
	repo := newMockRepo()
	svc := NewTaskService(repo)

	task, _ := svc.Create("Old Title", "Old Desc", "TODO")
	err := svc.Update(task.ID, "New Title", "New Desc", "DONE")
	assert.NoError(t, err)

	updated := repo.FindByID(task.ID)
	assert.Equal(t, "New Title", updated.Title)
	assert.Equal(t, "New Desc", updated.Description)
	assert.Equal(t, "DONE", updated.Status.String())
}

func TestTaskService_Delete(t *testing.T) {
	repo := newMockRepo()
	svc := NewTaskService(repo)

	task, _ := svc.Create("To Delete", "To Delete", "TODO")
	err := svc.Delete(task.ID)
	assert.NoError(t, err)
	assert.Nil(t, repo.FindByID(task.ID))
}

func TestTaskService_Create_InvalidStatus(t *testing.T) {
	repo := newMockRepo()
	svc := NewTaskService(repo)

	_, err := svc.Create("Bad Task", "Some desc", "INVALID_STATUS")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid status")
}

func TestTaskService_Update_NotFound(t *testing.T) {
	repo := newMockRepo()
	svc := NewTaskService(repo)

	err := svc.Update(999, "Title", "Desc", "TODO")
	assert.Error(t, err)
	assert.Equal(t, "task not found", err.Error())
}

func TestTaskService_Update_InvalidStatus(t *testing.T) {
	repo := newMockRepo()
	svc := NewTaskService(repo)

	task, _ := svc.Create("ToUpdate", "desc", "TODO")

	err := svc.Update(task.ID, "Updated", "desc", "UNKNOWN")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid status")
}

func TestTaskService_GetByID_NotFound(t *testing.T) {
	repo := newMockRepo()
	svc := NewTaskService(repo)

	task, err := svc.GetByID(1234)
	assert.Nil(t, task)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}
