package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"

	"todocli/internal/handlers"
	"todocli/internal/model"
	"todocli/internal/service"
)

func main() {
	// Канал завершения по Ctrl+C
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		fmt.Println("\nЗавершение по сигналу ОС...")
		cancel()
	}()

	taskChan := make(chan *model.Task)

	// Запуск фоновых сервисов
	go service.StartTaskGenerator(ctx, 5*time.Second, taskChan)
	go service.TaskSaver(ctx, taskChan)
	go service.StartLogger(ctx)

	// Инициализация маршрутов
	r := mux.NewRouter()
	r.HandleFunc("/api/item", handlers.CreateTaskHandler).Methods("POST")
	r.HandleFunc("/api/item/{id}", handlers.UpdateTaskHandler).Methods("PUT")
	r.HandleFunc("/api/items", handlers.GetAllTasksHandler).Methods("GET")
	r.HandleFunc("/api/item/{id}", handlers.GetTaskByIDHandler).Methods("GET")
	r.HandleFunc("/api/item/{id}", handlers.DeleteTaskHandler).Methods("DELETE")

	// Запуск сервера
	fmt.Println("Сервер запущен: http://localhost:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Println("Ошибка сервера:", err)
	}
}
