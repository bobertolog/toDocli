package main

// Мы импортируем  internal/model/task  для реализации логики проекта(Созданы директории и добавлен файл task.go)
import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"todocli/internal/model/task"
)

// Основная функция  работы программы
func main() {
	tasks := []*task.Task{}
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\n1. Добавить задачу")
		fmt.Println("2. Просмотреть задачи")
		fmt.Println("3. Обновить статус")
		fmt.Println("4. Удалить задачу")
		fmt.Println("5. Выйти")
		fmt.Print("Выберите действие: ")

		// Читаем выбор и переводим в int
		choiceStr, _ := reader.ReadString('\n')
		choiceStr = strings.TrimSpace(choiceStr)
		choice, err := strconv.Atoi(choiceStr)
		if err != nil {
			fmt.Println("Неверный ввод, введите число.")
			continue
		}

		//TODO  реализовать оставшуюся логику задач п3-4
		switch choice {
		case 1:
			fmt.Print("Название: ")
			title, _ := reader.ReadString('\n')
			title = strings.TrimSpace(title)

			fmt.Print("Статус: ")
			status, _ := reader.ReadString('\n')
			status = strings.TrimSpace(status)

			fmt.Print("Описание: ")
			desc, _ := reader.ReadString('\n')
			desc = strings.TrimSpace(desc)

			newTask := task.NewTask(len(tasks)+1, title, status, desc)
			tasks = append(tasks, newTask)
			fmt.Println("Задача добавлена!")

		case 2:
			for _, t := range tasks {
				fmt.Printf("ID: %d, Название: %s, Статус: %s, Создано: %s, Описание: %s\n",
					t.GetID(), t.Title, t.GetStatus(), t.CreatedAt.Format("2 Jan 2006 15:04"), t.Description)
			}

		case 3:
			fmt.Print("Введите ID задачи: ")
			break

		case 4:
			fmt.Print("Введите ID задачи: ")
			break

		case 5:
			fmt.Println("До встречи!") // Выходим из программы
			return

		default:
			fmt.Println("Неверный ввод, попробуйте снова.") //Обработчик неправильного ввода
		}
	}
}
