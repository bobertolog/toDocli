package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"todocli/internal/model"
	"todocli/internal/service"
)

// POST /api/item — создать задачу
func CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var t model.Task
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "Ошибка декодирования JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	if t.Title == "" {
		http.Error(w, "Поле 'title' обязательно", http.StatusBadRequest)
		return
	}

	if err := t.NormalizeStatus(); err != nil {
		http.Error(w, "Неверный статус: "+err.Error(), http.StatusBadRequest)
		return
	}

	t.ID = service.GenerateTaskID()
	service.AddTask(&t)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)
}

// PUT /api/item/{id} — обновить задачу по ID
func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Некорректный ID", http.StatusBadRequest)
		return
	}

	existing := service.FindTaskByID(id)
	if existing == nil {
		http.Error(w, "Задача не найдена", http.StatusNotFound)
		return
	}

	var updated model.Task
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		http.Error(w, "Ошибка декодирования JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Обновление данных задачи
	if updated.Title != "" {
		existing.Title = updated.Title
	}
	if updated.Description != "" {
		existing.Description = updated.Description
	}
	if updated.StatusRaw != "" {
		existing.StatusRaw = updated.StatusRaw
		if err := existing.NormalizeStatus(); err != nil {
			http.Error(w, "Неверный статус: "+err.Error(), http.StatusBadRequest)
			return
		}
	}

	if err := service.SaveAllTasks(); err != nil {
		http.Error(w, "Ошибка при сохранении: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(existing)
}

// GET /api/items — получить все задачи
func GetAllTasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks := service.GetAllTasks()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

// GET /api/item/{id} — получить задачу по ID
func GetTaskByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Некорректный ID", http.StatusBadRequest)
		return
	}

	task := service.FindTaskByID(id)
	if task == nil {
		http.Error(w, "Задача не найдена", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

// DELETE /api/item/{id} — удалить задачу
func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Некорректный ID", http.StatusBadRequest)
		return
	}

	if err := service.DeleteTask(id); err != nil {
		http.Error(w, "Ошибка удаления: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
