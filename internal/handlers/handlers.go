package handlers

import (
	"net/http"
	"strconv"

	"todocli/internal/model"
	"todocli/internal/repository"
	"todocli/internal/service"

	"github.com/gin-gonic/gin"
)

var taskService service.TaskService
var logger *repository.RedisLogger

func SetService(s service.TaskService) {
	taskService = s
}

func SetLogger(l *repository.RedisLogger) {
	logger = l
}

// POST /api/item
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
	task, err := taskService.Create(t.Title, t.Description, t.Status.String())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_ = logger.Log("create", task.Title)
	c.JSON(http.StatusCreated, task)
}

// GET /api/items
func GetAllTasks(c *gin.Context) {
	tasks := taskService.GetAll()
	c.JSON(http.StatusOK, tasks)
}

// GET /api/item/:id
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

// PUT /api/item/:id
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

	_ = logger.Log("update", strconv.Itoa(id))
	c.JSON(http.StatusOK, gin.H{"result": "updated"})
}

// DELETE /api/item/:id
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

	_ = logger.Log("delete", strconv.Itoa(id))
	c.Status(http.StatusNoContent)
}
