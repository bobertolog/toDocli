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
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	// Загружаем переменные окружения из .env
	_ = godotenv.Load()

	// Подключаемся к PostgreSQL
	repo, err := repository.NewPostgresRepository()
	if err != nil {
		log.Fatalf("Ошибка подключения к PostgreSQL: %v", err)
	}

	// Создаём сервис задач
	taskService := service.NewTaskService(repo)

	// Подключаемся к Redis для логгирования
	redisClient := connectRedis()
	logger := repository.NewRedisLogger(redisClient)

	// Передаём зависимости в handlers
	handlers.SetService(taskService)
	handlers.SetLogger(logger)

	// Настраиваем Gin HTTP сервер
	r := gin.Default()

	// Роуты
	r.POST("/login", handlers.Login)

	auth := r.Group("/api", middleware.JWTAuth())
	auth.POST("/item", handlers.CreateTask)
	auth.GET("/items", handlers.GetAllTasks)
	auth.GET("/item/:id", handlers.GetTaskByID)
	auth.PUT("/item/:id", handlers.UpdateTask)
	auth.DELETE("/item/:id", handlers.DeleteTask)

	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}

func connectRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}
