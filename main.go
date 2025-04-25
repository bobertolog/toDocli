// main.go
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"todocli/internal/model"
	"todocli/internal/service"
)

func main() {
	taskChan := make(chan *model.Task)
	done := make(chan struct{})

	go service.StartTaskGenerator(5*time.Second, taskChan)
	go service.TaskSaver(taskChan)
	go service.StartLogger(done)

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\n1. Добавить задачу")
		fmt.Println("2. Просмотреть задачи")
		fmt.Println("3. Обновить статус")
		fmt.Println("4. Удалить задачу")
		fmt.Println("5. Выйти")
		fmt.Print("Выберите действие: ")

		choiceStr, _ := reader.ReadString('\n')
		choice, err := strconv.Atoi(strings.TrimSpace(choiceStr))
		if err != nil {
			fmt.Println("Ошибка ввода:", err)
			continue
		}

		switch choice {
		case 1:
			fmt.Print("Название: ")
			title, _ := reader.ReadString('\n')
			title = strings.TrimSpace(title)

			fmt.Print("Описание: ")
			desc, _ := reader.ReadString('\n')
			desc = strings.TrimSpace(desc)

			task := model.NewTask(service.GenerateTaskID(), title, desc)
			service.AddTask(task)
			fmt.Println("Задача добавлена!")

		case 2:
			tasks := service.GetAllTasks()
			if len(tasks) == 0 {
				fmt.Println("Нет задач.")
				continue
			}
			for _, task := range tasks {
				fmt.Printf("ID: %d, Название: %s, Статус: %s, Описание: %s\n",
					task.GetEntityID(), task.Title, task.Status(), task.Description)
			}

		case 3:
			fmt.Print("Введите ID задачи: ")
			idStr, _ := reader.ReadString('\n')
			id, err := strconv.Atoi(strings.TrimSpace(idStr))
			if err != nil {
				fmt.Println("Ошибка: неверный ID")
				continue
			}

			task := service.FindTaskByID(id)
			if task == nil {
				fmt.Println("Задача не найдена!")
				continue
			}

			fmt.Println("Выберите статус:")
			fmt.Println("1. TODO")
			fmt.Println("2. IN_PROGRESS")
			fmt.Println("3. DONE")
			fmt.Print("Ваш выбор: ")
			statusChoice, _ := reader.ReadString('\n')
			statusChoice = strings.TrimSpace(statusChoice)

			switch statusChoice {
			case "1":
				task.SetStatusType(model.StatusTodo)
			case "2":
				task.SetStatusType(model.StatusInProgress)
			case "3":
				task.SetStatusType(model.StatusDone)
			default:
				fmt.Println("Неверный статус!")
				continue
			}
			fmt.Println("Статус обновлён!")

		case 4:
			fmt.Print("Введите ID задачи: ")
			idStr, _ := reader.ReadString('\n')
			id, err := strconv.Atoi(strings.TrimSpace(idStr))
			if err != nil {
				fmt.Println("Ошибка: неверный ID")
				continue
			}

			err = service.DeleteTask(id)
			if err != nil {
				fmt.Println("Ошибка:", err)
				continue
			}
			fmt.Println("Задача удалена!")

		case 5:
			fmt.Println("Выход...")
			close(done)
			return

		default:
			fmt.Println("Неверный выбор!")
		}
	}
}
