package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

const dataFile = "tasks.json"

// Storage управляет хранением задач
type Storage struct {
	Tasks  []*Task `json:"tasks"`
	NextID int     `json:"next_id"`
}

// NewStorage создает новое хранилище
func NewStorage() *Storage {
	return &Storage{
		Tasks:  make([]*Task, 0),
		NextID: 1,
	}
}

// Load загружает задачи из файла
func (s *Storage) Load() error {
	data, err := os.ReadFile(dataFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // Файл не существует - это нормально для первого запуска
		}
		return fmt.Errorf("ошибка чтения файла: %w", err)
	}

	if len(data) == 0 {
		return nil
	}

	err = json.Unmarshal(data, s)
	if err != nil {
		return fmt.Errorf("ошибка парсинга JSON: %w", err)
	}

	return nil
}

// Save сохраняет задачи в файл
func (s *Storage) Save() error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("ошибка сериализации: %w", err)
	}

	err = os.WriteFile(dataFile, data, 0644)
	if err != nil {
		return fmt.Errorf("ошибка записи файла: %w", err)
	}

	return nil
}

// AddTask добавляет новую задачу
func (s *Storage) AddTask(title, description string) *Task {
	task := NewTask(s.NextID, title, description)
	s.Tasks = append(s.Tasks, task)
	s.NextID++
	return task
}

// GetTask получает задачу по ID
func (s *Storage) GetTask(id int) *Task {
	for _, task := range s.Tasks {
		if task.ID == id {
			return task
		}
	}
	return nil
}

// DeleteTask удаляет задачу по ID
func (s *Storage) DeleteTask(id int) bool {
	for i, task := range s.Tasks {
		if task.ID == id {
			s.Tasks = append(s.Tasks[:i], s.Tasks[i+1:]...)
			return true
		}
	}
	return false
}

// ListTasks возвращает все задачи
func (s *Storage) ListTasks() []*Task {
	return s.Tasks
}

// FilterTasksByStatus возвращает задачи с указанным статусом
func (s *Storage) FilterTasksByStatus(status string) []*Task {
	filtered := make([]*Task, 0)
	for _, task := range s.Tasks {
		if task.Status == status {
			filtered = append(filtered, task)
		}
	}
	return filtered
}

// SearchTasks ищет задачи по тексту в названии или описании
func (s *Storage) SearchTasks(query string) []*Task {
	results := make([]*Task, 0)
	queryLower := strings.ToLower(query)

	for _, task := range s.Tasks {
		titleLower := strings.ToLower(task.Title)
		descLower := strings.ToLower(task.Description)

		if strings.Contains(titleLower, queryLower) || strings.Contains(descLower, queryLower) {
			results = append(results, task)
		}
	}
	return results
}
