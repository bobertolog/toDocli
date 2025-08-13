package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"todocli/internal/model"
)

type APIClient struct {
	BaseURL string
}

func NewAPIClient(baseURL string) *APIClient {
	return &APIClient{BaseURL: baseURL}
}

func (c *APIClient) CreateTask(title, desc, status string) (*model.Task, error) {
	body := map[string]string{
		"title":       title,
		"description": desc,
		"status":      status,
	}
	data, _ := json.Marshal(body)

	resp, err := http.Post(c.BaseURL+"/api/item", "application/json", bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("ошибка: %s", resp.Status)
	}

	var t model.Task
	if err := json.NewDecoder(resp.Body).Decode(&t); err != nil {
		return nil, err
	}
	return &t, nil
}

func (c *APIClient) GetAll() ([]*model.Task, error) {
	resp, err := http.Get(c.BaseURL + "/api/items")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var tasks []*model.Task
	if err := json.NewDecoder(resp.Body).Decode(&tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (c *APIClient) GetByID(id int) (*model.Task, error) {
	resp, err := http.Get(fmt.Sprintf("%s/api/item/%d", c.BaseURL, id))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("task not found")
	}

	var t model.Task
	if err := json.NewDecoder(resp.Body).Decode(&t); err != nil {
		return nil, err
	}
	return &t, nil
}

func (c *APIClient) Update(id int, title, desc, status string) error {
	body := map[string]string{
		"title":       title,
		"description": desc,
		"status":      status,
	}
	data, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/api/item/%d", c.BaseURL, id), bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		out, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("ошибка: %s — %s", resp.Status, string(out))
	}
	return nil
}

func (c *APIClient) Delete(id int) error {
	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/api/item/%d", c.BaseURL, id), nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("ошибка удаления: %s", resp.Status)
	}
	return nil
}
