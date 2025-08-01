package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

const baseURL = "http://localhost:8080/api"
const loginURL = "http://localhost:8080/login"

var jwtToken string

func main() {
	login()

	for {
		fmt.Println("\n== TODO CLI ==")
		fmt.Println("1. Создать задачу")
		fmt.Println("2. Показать все задачи")
		fmt.Println("3. Найти задачу по ID")
		fmt.Println("4. Обновить задачу")
		fmt.Println("5. Удалить задачу")
		fmt.Println("0. Выход")

		var choice int
		fmt.Print("Выберите опцию: ")
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			createTask()
		case 2:
			getAllTasks()
		case 3:
			getTaskByID()
		case 4:
			updateTask()
		case 5:
			deleteTask()
		case 0:
			fmt.Println("Выход.")
			return
		default:
			fmt.Println("Неверная опция")
		}
	}
}

func login() {
	fmt.Println("== Авторизация ==")
	fmt.Print("Логин: ")
	var username string
	fmt.Scanln(&username)

	fmt.Print("Пароль: ")
	var password string
	fmt.Scanln(&password)

	creds := map[string]string{
		"username": username,
		"password": password,
	}
	body, _ := json.Marshal(creds)

	resp, err := http.Post(loginURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("Ошибка запроса:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Ошибка авторизации:", resp.Status)
		io.Copy(os.Stdout, resp.Body)
		os.Exit(1)
	}

	var result map[string]string
	json.NewDecoder(resp.Body).Decode(&result)
	jwtToken = result["token"]
	fmt.Println("Успешно авторизован!")
}

func createTask() {
	var title, description, status string
	fmt.Print("Название: ")
	fmt.Scanln(&title)
	fmt.Print("Описание: ")
	fmt.Scanln(&description)
	fmt.Print("Статус (TODO, IN_PROGRESS, DONE): ")
	fmt.Scanln(&status)

	task := map[string]string{
		"title":       title,
		"description": description,
		"status":      status,
	}
	body, _ := json.Marshal(task)

	req, _ := http.NewRequest("POST", baseURL+"/item", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+jwtToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	defer resp.Body.Close()

	io.Copy(os.Stdout, resp.Body)
}

func getAllTasks() {
	req, _ := http.NewRequest("GET", baseURL+"/items", nil)
	req.Header.Set("Authorization", "Bearer "+jwtToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	defer resp.Body.Close()

	io.Copy(os.Stdout, resp.Body)
}

func getTaskByID() {
	var id int
	fmt.Print("ID задачи: ")
	fmt.Scanln(&id)

	req, _ := http.NewRequest("GET", baseURL+"/item/"+strconv.Itoa(id), nil)
	req.Header.Set("Authorization", "Bearer "+jwtToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	defer resp.Body.Close()

	io.Copy(os.Stdout, resp.Body)
}

func updateTask() {
	var id int
	fmt.Print("ID задачи: ")
	fmt.Scanln(&id)

	var title, description, status string
	fmt.Print("Новое название: ")
	fmt.Scanln(&title)
	fmt.Print("Новое описание: ")
	fmt.Scanln(&description)
	fmt.Print("Новый статус (TODO, IN_PROGRESS, DONE): ")
	fmt.Scanln(&status)

	update := map[string]string{
		"title":       title,
		"description": description,
		"status":      status,
	}
	body, _ := json.Marshal(update)

	req, _ := http.NewRequest("PUT", baseURL+"/item/"+strconv.Itoa(id), bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+jwtToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	defer resp.Body.Close()

	io.Copy(os.Stdout, resp.Body)
}

func deleteTask() {
	var id int
	fmt.Print("ID задачи: ")
	fmt.Scanln(&id)

	req, _ := http.NewRequest("DELETE", baseURL+"/item/"+strconv.Itoa(id), nil)
	req.Header.Set("Authorization", "Bearer "+jwtToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		fmt.Println("Задача удалена.")
	} else {
		io.Copy(os.Stdout, resp.Body)
	}
}
