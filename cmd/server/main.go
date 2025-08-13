// @title Task Manager API
// @version 1.0
// @description API для управления задачами
// @host localhost:8080
// @BasePath /
// @schemes http
package main

import (
	"log"
	"os"

	"todocli/internal/handlers"
	"todocli/internal/handlers/middleware"
	"todocli/internal/repository"
	"todocli/internal/service"

	_ "todocli/docs"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {

	repo, err := repository.NewPostgresRepository()
	if err != nil {
		log.Fatalf("Ошибка подключения к PostgreSQL: %v", err)
	}

	taskService := service.NewTaskService(repo)

	// === Redis ===
	redisClient := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"), // берётся из docker-compose
	})

	logger := repository.NewRedisLogger(redisClient)

	handlers.SetService(taskService)
	handlers.SetLogger(logger)
	handlers.SetUserRepo(repo)

	r := gin.Default()

	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)

	api := r.Group("/api", middleware.JWTAuth())
	{
		api.POST("/item", handlers.CreateTask)
		api.GET("/items", handlers.GetAllTasks)
		api.GET("/item/:id", handlers.GetTaskByID)
		api.PUT("/item/:id", handlers.UpdateTask)
		api.DELETE("/item/:id", handlers.DeleteTask)
	}

	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
