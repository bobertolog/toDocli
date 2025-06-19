package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"todocli/pb"

	"google.golang.org/grpc"
)

func main() {
	// Флаги
	listFlag := flag.Bool("list", false, "Вывести все задачи")
	createFlag := flag.Bool("create", false, "Создать новую задачу")
	getFlag := flag.Int("get", 0, "Получить задачу по ID")
	deleteFlag := flag.Int("delete", 0, "Удалить задачу по ID")

	title := flag.String("title", "", "Заголовок задачи (для создания)")
	desc := flag.String("desc", "", "Описание задачи (для создания)")
	status := flag.String("status", "TODO", "Статус задачи: TODO, IN_PROGRESS, DONE")

	flag.Parse()

	// Подключение к gRPC
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Ошибка подключения к серверу: %v", err)
	}
	defer conn.Close()

	client := pb.NewTaskServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Обработка команд
	switch {
	case *listFlag:
		res, err := client.ListTasks(ctx, &pb.Empty{})
		if err != nil {
			log.Fatalf("Ошибка ListTasks: %v", err)
		}
		if len(res.Tasks) == 0 {
			fmt.Println("Список задач пуст.")
			return
		}
		fmt.Println("Список задач:")
		for _, t := range res.Tasks {
			fmt.Printf("  [%d] %s (%s)\n    %s\n", t.Id, t.Title, t.Status, t.Description)
		}

	case *createFlag:
		if *title == "" || *desc == "" {
			fmt.Println("Для создания задачи укажите --title и --desc")
			os.Exit(1)
		}
		task := &pb.Task{
			Title:       *title,
			Description: *desc,
			Status:      *status,
		}
		res, err := client.CreateTask(ctx, task)
		if err != nil {
			log.Fatalf("Ошибка CreateTask: %v", err)
		}
		fmt.Printf("Задача создана с ID: %d\n", res.Id)

	case *getFlag > 0:
		res, err := client.GetTask(ctx, &pb.TaskID{Id: int32(*getFlag)})
		if err != nil {
			log.Fatalf("Ошибка GetTask: %v", err)
		}
		fmt.Printf("Задача %d:\n  Заголовок: %s\n  Статус: %s\n  Описание: %s\n",
			res.Id, res.Title, res.Status, res.Description)

	case *deleteFlag > 0:
		_, err := client.DeleteTask(ctx, &pb.TaskID{Id: int32(*deleteFlag)})
		if err != nil {
			log.Fatalf("Ошибка DeleteTask: %v", err)
		}
		fmt.Printf("Задача %d удалена.\n", *deleteFlag)

	default:
		fmt.Println("Использование:")
		fmt.Println("  --list                        Показать все задачи")
		fmt.Println("  --get <id>                   Получить задачу по ID")
		fmt.Println("  --delete <id>                Удалить задачу по ID")
		fmt.Println("  --create --title T --desc D  Создать задачу")
		fmt.Println("         [--status TODO|IN_PROGRESS|DONE]")
	}
}
