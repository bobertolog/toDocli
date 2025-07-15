package main

import (
	"context"
	"flag"
	"log"
	"time"

	"todocli/pb"

	"google.golang.org/grpc"
)

func main() {
	// Флаг для удаления задачи по ID
	deleteID := flag.Int("delete", 0, "ID задачи для удаления")
	flag.Parse()

	// Подключение к gRPC серверу
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewTaskServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Если передан флаг --delete
	if *deleteID > 0 {
		_, err := client.DeleteTask(ctx, &pb.TaskID{Id: int32(*deleteID)})
		if err != nil {
			log.Fatalf("Ошибка удаления задачи с ID %d: %v", *deleteID, err)
		}
		log.Printf("Задача с ID %d успешно удалена.\n", *deleteID)
		return
	}

	// 1. Создание задачи
	task := &pb.Task{
		Title:       "Тут будет размещаться название задачи",
		Description: "А тут будет описание задачи",
		Status:      "TODO",
	}
	res, err := client.CreateTask(ctx, task)
	if err != nil {
		log.Fatal("CreateTask error:", err)
	}
	log.Println("Создана задача с ID:", res.Id)

	// 2. Получение задачи по ID
	got, err := client.GetTask(ctx, &pb.TaskID{Id: res.Id})
	if err != nil {
		log.Fatal("GetTask error:", err)
	}
	log.Printf("Задача получена: %+v\n", got)

	// 3. Список всех задач
	list, err := client.ListTasks(ctx, &pb.Empty{})
	if err != nil {
		log.Fatal("ListTasks error:", err)
	}
	log.Println("Все задачи:")
	for _, t := range list.Tasks {
		log.Printf(" - [%d] %s (%s)", t.Id, t.Title, t.Status)
	}
}
