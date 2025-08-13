package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"todocli/internal/model"
	"todocli/internal/service"

	"github.com/gin-gonic/gin"
)

var taskService service.TaskService

type Logger interface {
	Log(msg string) error
}

var logger Logger

func SetService(s service.TaskService) {
	taskService = s
}

func SetLogger(l Logger) {
	logger = l
}

// CreateTask godoc
// @Summary Создать задачу
// @Description Добавляет новую задачу
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body model.Task true "Задача"
// @Success 201 {object} model.Task
// @Failure 400 {object} map[string]string
// @Router /api/item [post]
// @Security ApiKeyAuth
func CreateTask(c *gin.Context) {
	var t model.Task
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}
	if t.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title required"})
		return
	}
	task, err := taskService.CreateWithLog(t.Title, t.Description, t.Status.String(), func(msg string) error {
		if logger != nil {
			return logger.Log(msg)
		}
		return nil
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, task)
}

// GetAllTasks godoc
// @Summary Получить все задачи
// @Description Возвращает список всех задач
// @Tags tasks
// @Produce json
// @Success 200 {array} model.Task
// @Router /api/items [get]
// @Security ApiKeyAuth
func GetAllTasks(c *gin.Context) {
	tasks := taskService.GetAll()
	c.JSON(http.StatusOK, tasks)
}

// GetTaskByID godoc
// @Summary Получить задачу по ID
// @Description Возвращает задачу по идентификатору
// @Tags tasks
// @Produce json
// @Param id path int true "ID задачи"
// @Success 200 {object} model.Task
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/item/{id} [get]
// @Security ApiKeyAuth
func GetTaskByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	task, err := taskService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, task)
}

// UpdateTask godoc
// @Summary Обновить задачу
// @Description Обновляет данные задачи
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "ID задачи"
// @Param task body model.Task true "Новые данные задачи"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /api/item/{id} [put]
// @Security ApiKeyAuth
func UpdateTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var upd model.Task
	if err := c.ShouldBindJSON(&upd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	err = taskService.Update(id, upd.Title, upd.Description, upd.Status.String())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if logger != nil {
		_ = logger.Log(fmt.Sprintf("update: %d", id))
	}
	c.JSON(http.StatusOK, gin.H{"result": "updated"})
}

// DeleteTask godoc
// @Summary Удалить задачу
// @Description Удаляет задачу по ID
// @Tags tasks
// @Produce json
// @Param id path int true "ID задачи"
// @Success 204 {string} string "No Content"
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/item/{id} [delete]
// @Security ApiKeyAuth
func DeleteTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	err = taskService.Delete(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if logger != nil {
		_ = logger.Log(fmt.Sprintf("delete: %d", id))
	}
	c.Status(http.StatusNoContent)
}
