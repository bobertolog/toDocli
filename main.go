package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"todocli/internal/model"
)

func main() {
	tasks := []*model.Task{}
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\n1. Добавить задачу")
		fmt.Println("2. Просмотреть задачи")
		fmt.Println("3. Обновить статус")
		fmt.Println("4. Удалить задачу")
		fmt.Println("5. Выйти")
		fmt.Print("Выберите действие: ")

		choiceStr, _ := reader.ReadString('\n')
		choiceStr = strings.TrimSpace(choiceStr)
		choice, err := strconv.Atoi(choiceStr)
		if err != nil {
			fmt.Println("Неверный ввод, введите число.")
			continue
		}

		switch choice {
		case 1: // Добавление задачи
			fmt.Print("Название: ")
			title, _ := reader.ReadString('\n')
			title = strings.TrimSpace(title)

			fmt.Print("Описание: ")
			desc, _ := reader.ReadString('\n')
			desc = strings.TrimSpace(desc)

			// Статус по умолчанию будет "TODO", так как мы убрали его из конструктора
			newTask := model.NewTask(len(tasks)+1, title, desc)
			tasks = append(tasks, newTask)
			fmt.Println("Задача добавлена!")

		case 2: // Просмотр задач
			if len(tasks) == 0 {
				fmt.Println("Нет задач для отображения.")
				continue
			}
			for _, t := range tasks {
				fmt.Printf("ID: %d, Название: %s, Статус: %s, Создано: %s, Описание: %s\n",
					t.ID, t.Title, t.Status(), t.CreatedAt.Format("2 Jan 2006 15:04"), t.Description)
			}

		case 3: // Обновление статуса
			if len(tasks) == 0 {
				fmt.Println("Нет задач для обновления.")
				continue
			}

			fmt.Print("Введите ID задачи: ")
			idStr, _ := reader.ReadString('\n')
			id, err := strconv.Atoi(strings.TrimSpace(idStr))
			if err != nil || id < 1 || id > len(tasks) {
				fmt.Println("Неверный ID задачи.")
				continue
			}

			fmt.Println("Доступные статусы:")
			fmt.Println("1. TODO")
			fmt.Println("2. IN_PROGRESS")
			fmt.Println("3. DONE")
			fmt.Print("Выберите новый статус (1-3): ")

			statusChoiceStr, _ := reader.ReadString('\n')
			statusChoice, err := strconv.Atoi(strings.TrimSpace(statusChoiceStr))
			if err != nil || statusChoice < 1 || statusChoice > 3 {
				fmt.Println("Неверный выбор статуса.")
				continue
			}

			// Конвертируем выбор в StatusType
			var newStatus model.StatusType
			switch statusChoice {
			case 1:
				newStatus = model.StatusTodo
			case 2:
				newStatus = model.StatusInProgress
			case 3:
				newStatus = model.StatusDone
			}

			tasks[id-1].SetStatusType(newStatus)
			fmt.Println("Статус обновлен!")

		case 4: // Удаление задачи
			if len(tasks) == 0 {
				fmt.Println("Нет задач для удаления.")
				continue
			}

			fmt.Print("Введите ID задачи: ")
			idStr, _ := reader.ReadString('\n')
			id, err := strconv.Atoi(strings.TrimSpace(idStr))
			if err != nil || id < 1 || id > len(tasks) {
				fmt.Println("Неверный ID задачи.")
				continue
			}

			// Удаляем задачу из слайса
			tasks = append(tasks[:id-1], tasks[id:]...)

			// Обновляем ID оставшихся задач
			for i := id - 1; i < len(tasks); i++ {
				tasks[i].ID = i + 1
			}

			fmt.Println("Задача удалена!")

		case 5: // Выход
			fmt.Println("До встречи!")
			return

		default:
			fmt.Println("Неверный ввод, попробуйте снова.")
		}
	}
}
