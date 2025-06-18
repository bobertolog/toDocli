package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"todocli/internal/handlers"
	"todocli/internal/handlers/middleware"
	"todocli/internal/model"
	"todocli/internal/repository"
	"todocli/internal/service"

	_ "todocli/docs"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}
	if err := repository.Init(); err != nil {
		fmt.Println("init repository error:", err)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go service.StartTaskGenerator(ctx, 5*time.Second, make(chan *model.Task))
	go service.TaskSaver(ctx, make(chan *model.Task))
	go service.StartLogger(ctx)

	r := gin.Default()
	r.POST("/login", handlers.Login)
	auth := r.Group("/api", middleware.JWTAuth())
	auth.POST("/item", handlers.CreateTask)
	auth.GET("/items", handlers.GetAllTasks)
	auth.GET("/item/:id", handlers.GetTaskByID)
	auth.PUT("/item/:id", handlers.UpdateTask)
	auth.DELETE("/item/:id", handlers.DeleteTask)

	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run(":8080")
}
