package main

import (
	"context"
	"log"

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

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Загрузка .env (если есть)
	err := godotenv.Load()
	if err != nil {
		log.Println("Не найден .env, продолжаем без него")
	}

	// Подключение к MongoDB
	mongoCol := connectMongo()
	repo := repository.NewMongoRepo(mongoCol)

	// Подключение к Redis
	redisClient := connectRedis()
	logger := repository.NewRedisLogger(redisClient)

	// Создаём сервис с Mongo-репозиторием
	taskService := service.NewTaskService(repo)

	// Передаём зависимости в handlers
	handlers.SetService(taskService)
	handlers.SetLogger(logger)

	// HTTP-сервер
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

func connectMongo() *mongo.Collection {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal("Mongo connect error:", err)
	}
	return client.Database("tododb").Collection("tasks")
}

func connectRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}
