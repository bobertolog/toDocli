package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"todocli/internal/repository"
	"todocli/internal/service"
)

func main() {
	repo := repository.NewInMemoryRepository() // Можно заменить на PostgresRepository
	svc := service.NewTaskService(repo)

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\n=== TODO CLI ===")
		fmt.Println("1. Добавить задачу")
		fmt.Println("2. Посмотреть все задачи")
		fmt.Println("3. Найти задачу по ID")
		fmt.Println("4. Обновить задачу")
		fmt.Println("5. Удалить задачу")
		fmt.Println("0. Выход")
		fmt.Print("Выберите опцию: ")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			fmt.Print("Введите заголовок: ")
			title, _ := reader.ReadString('\n')
			fmt.Print("Введите описание: ")
			desc, _ := reader.ReadString('\n')
			fmt.Print("Введите статус (pending, in_progress, done): ")
			status, _ := reader.ReadString('\n')

			task, err := svc.Create(strings.TrimSpace(title), strings.TrimSpace(desc), strings.TrimSpace(status))
			if err != nil {
				fmt.Println("Ошибка:", err)
			} else {
				fmt.Println("Создана задача:", task)
			}

		case "2":
			tasks := svc.GetAll()
			for _, t := range tasks {
				fmt.Printf("[%d] %s - %s (%s)\n", t.ID, t.Title, t.Description, t.Status)
			}

		case "3":
			fmt.Print("Введите ID задачи: ")
			idStr, _ := reader.ReadString('\n')
			id, _ := strconv.Atoi(strings.TrimSpace(idStr))
			task, err := svc.GetByID(id)
			if err != nil {
				fmt.Println("Ошибка:", err)
			} else {
				fmt.Printf("[%d] %s - %s (%s)\n", task.ID, task.Title, task.Description, task.Status)
			}

		case "4":
			fmt.Print("ID задачи: ")
			idStr, _ := reader.ReadString('\n')
			id, _ := strconv.Atoi(strings.TrimSpace(idStr))

			fmt.Print("Новый заголовок: ")
			title, _ := reader.ReadString('\n')
			fmt.Print("Новое описание: ")
			desc, _ := reader.ReadString('\n')
			fmt.Print("Новый статус: ")
			status, _ := reader.ReadString('\n')

			err := svc.Update(id, strings.TrimSpace(title), strings.TrimSpace(desc), strings.TrimSpace(status))
			if err != nil {
				fmt.Println("Ошибка:", err)
			} else {
				fmt.Println("Задача обновлена.")
			}

		case "5":
			fmt.Print("ID задачи: ")
			idStr, _ := reader.ReadString('\n')
			id, _ := strconv.Atoi(strings.TrimSpace(idStr))

			err := svc.Delete(id)
			if err != nil {
				fmt.Println("Ошибка:", err)
			} else {
				fmt.Println("Задача удалена.")
			}
		case "6":
			fmt.Print("Введите статус (pending, in_progress, done): ")
			status, _ := reader.ReadString('\n')
			tasks := svc.GetByStatus(strings.TrimSpace(status))
			if len(tasks) == 0 {
				fmt.Println("Нет задач с таким статусом.")
			}
			for _, t := range tasks {
				fmt.Printf("[%d] %s - %s (%s)\n", t.ID, t.Title, t.Description, t.Status)
			}
		case "0":
			fmt.Println("Выход...")
			return

		default:
			fmt.Println("Неизвестная опция.")
		}
	}
}
