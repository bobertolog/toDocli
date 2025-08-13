package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"todocli/internal/model"
	"todocli/internal/repository"
	"todocli/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouterForTest() *gin.Engine {
	repo := repository.NewInMemoryRepository()
	svc := service.NewTaskService(repo)
	log := repository.NewFakeLogger() // используется заглушка логгера

	SetService(svc)
	SetLogger(log)

	r := gin.Default()
	r.POST("/api/item", CreateTask)
	r.GET("/api/items", GetAllTasks)
	r.GET("/api/item/:id", GetTaskByID)
	r.PUT("/api/item/:id", UpdateTask)
	r.DELETE("/api/item/:id", DeleteTask)

	return r
}

func TestTaskHandlers_FullFlow(t *testing.T) {
	r := setupRouterForTest()

	// --- 1. Create
	body := map[string]string{
		"title":       "API Task",
		"description": "From test",
		"status":      "TODO",
	}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/item", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)
	var created model.Task
	json.Unmarshal(resp.Body.Bytes(), &created)
	assert.Equal(t, "API Task", created.Title)

	// --- 2. Get by ID
	req = httptest.NewRequest(http.MethodGet, "/api/item/"+toStr(created.ID), nil)
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	// --- 3. Update
	update := map[string]string{
		"title":       "Updated",
		"description": "Updated",
		"status":      "DONE",
	}
	updJson, _ := json.Marshal(update)
	req = httptest.NewRequest(http.MethodPut, "/api/item/"+toStr(created.ID), bytes.NewBuffer(updJson))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	// --- 4. Delete
	req = httptest.NewRequest(http.MethodDelete, "/api/item/"+toStr(created.ID), nil)
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusNoContent, resp.Code)
}

func toStr(i int) string {
	return fmt.Sprintf("%d", i)
}
