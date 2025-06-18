package handlers

import (
	"net/http"
	"strconv"

	"todocli/internal/model"
	"todocli/internal/service"

	"github.com/gin-gonic/gin"
)

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
	if err := t.Normalize(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status"})
		return
	}
	t.ID = service.GenerateTaskID()
	service.AddTask(&t)
	c.JSON(http.StatusCreated, t)
}

func GetAllTasks(c *gin.Context) {
	c.JSON(http.StatusOK, service.GetAllTasks())
}

func GetTaskByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if t := service.FindTaskByID(id); t != nil {
		c.JSON(http.StatusOK, t)
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	}
}

func UpdateTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	ex := service.FindTaskByID(id)
	if ex == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	var upd model.Task
	if err := c.ShouldBindJSON(&upd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}
	if upd.Title != "" {
		ex.Title = upd.Title
	}
	if upd.Description != "" {
		ex.Description = upd.Description
	}
	if upd.StatusRaw != "" {
		ex.StatusRaw = upd.StatusRaw
		if err := ex.Normalize(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status"})
			return
		}
	}
	service.SaveAllTasks()
	c.JSON(http.StatusOK, ex)
}

func DeleteTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := service.DeleteTask(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.Status(http.StatusNoContent)
}
