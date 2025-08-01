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
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Файл .env не найден, переменные окружения должны быть заданы вручную")
	}

	// Инициализация PostgreSQL
	repo, err := repository.NewPostgresRepository()
	if err != nil {
		log.Fatalf("Ошибка подключения к PostgreSQL: %v", err)
	}

	// Создание сервиса задач
	taskService := service.NewTaskService(repo)

	// Инициализация Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	logger := repository.NewRedisLogger(redisClient)

	// Подключение зависимостей
	handlers.SetService(taskService)
	handlers.SetLogger(logger)

	// Настройка маршрутов
	r := gin.Default()

	r.POST("/login", handlers.Login)

	api := r.Group("/api", middleware.JWTAuth())
	{
		api.POST("/item", handlers.CreateTask)
		api.GET("/items", handlers.GetAllTasks)
		api.GET("/item/:id", handlers.GetTaskByID)
		api.PUT("/item/:id", handlers.UpdateTask)
		api.DELETE("/item/:id", handlers.DeleteTask)
	}

	// Swagger
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Запуск сервера
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
