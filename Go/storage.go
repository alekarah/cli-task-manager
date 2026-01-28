package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
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

// FilterTasksByTag возвращает задачи с указанным тегом
func (s *Storage) FilterTasksByTag(tag string) []*Task {
	filtered := make([]*Task, 0)
	tagLower := strings.TrimSpace(strings.ToLower(tag))

	for _, task := range s.Tasks {
		for _, t := range task.Tags {
			if t == tagLower {
				filtered = append(filtered, task)
				break
			}
		}
	}
	return filtered
}

// GetAllTags возвращает список всех уникальных тегов
func (s *Storage) GetAllTags() []string {
	tagsMap := make(map[string]bool)
	for _, task := range s.Tasks {
		for _, tag := range task.Tags {
			tagsMap[tag] = true
		}
	}

	tags := make([]string, 0, len(tagsMap))
	for tag := range tagsMap {
		tags = append(tags, tag)
	}
	sort.Strings(tags)
	return tags
}

// SortTasks сортирует и возвращает задачи по указанному критерию
func (s *Storage) SortTasks(sortBy string) []*Task {
	tasks := make([]*Task, len(s.Tasks))
	copy(tasks, s.Tasks)

	switch sortBy {
	case "id":
		sort.Slice(tasks, func(i, j int) bool {
			return tasks[i].ID < tasks[j].ID
		})
	case "created":
		sort.Slice(tasks, func(i, j int) bool {
			return tasks[i].CreatedAt.Before(tasks[j].CreatedAt)
		})
	case "updated":
		sort.Slice(tasks, func(i, j int) bool {
			return tasks[i].UpdatedAt.After(tasks[j].UpdatedAt)
		})
	case "status":
		sort.Slice(tasks, func(i, j int) bool {
			return tasks[i].Status < tasks[j].Status
		})
	case "priority":
		priorityOrder := map[string]int{"high": 1, "medium": 2, "low": 3}
		sort.Slice(tasks, func(i, j int) bool {
			return priorityOrder[tasks[i].Priority] < priorityOrder[tasks[j].Priority]
		})
	}

	return tasks
}

// ExportToCSV экспортирует задачи в CSV файл
func (s *Storage) ExportToCSV(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("ошибка создания файла: %w", err)
	}
	defer file.Close()

	// Записываем BOM для корректного отображения UTF-8 в Excel
	file.WriteString("\xEF\xBB\xBF")

	// Заголовки
	fmt.Fprintln(file, "ID,Название,Описание,Статус,Приоритет,Дедлайн,Теги,Создано,Обновлено")

	// Данные
	for _, task := range s.Tasks {
		deadlineStr := ""
		if task.Deadline != nil {
			deadlineStr = task.Deadline.Format("02.01.2006 15:04")
		}

		tagsStr := ""
		if len(task.Tags) > 0 {
			tagStrings := make([]string, len(task.Tags))
			for i, tag := range task.Tags {
				tagStrings[i] = "#" + tag
			}
			tagsStr = strings.Join(tagStrings, ", ")
		}

		// Экранируем значения для CSV
		fmt.Fprintf(file, "%d,\"%s\",\"%s\",%s,%s,\"%s\",\"%s\",\"%s\",\"%s\"\n",
			task.ID,
			escapeCSV(task.Title),
			escapeCSV(task.Description),
			task.Status,
			task.Priority,
			deadlineStr,
			tagsStr,
			task.CreatedAt.Format("02.01.2006 15:04"),
			task.UpdatedAt.Format("02.01.2006 15:04"),
		)
	}

	return nil
}

// ExportToMarkdown экспортирует задачи в Markdown файл
func (s *Storage) ExportToMarkdown(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("ошибка создания файла: %w", err)
	}
	defer file.Close()

	fmt.Fprintln(file, "# Список задач\n")

	if len(s.Tasks) == 0 {
		fmt.Fprintln(file, "*Задач нет*")
		return nil
	}

	// Группируем по статусу
	statuses := []struct {
		key  string
		name string
	}{
		{"todo", "📋 К выполнению"},
		{"in_progress", "⚙️ В процессе"},
		{"done", "✅ Выполнено"},
	}

	for _, status := range statuses {
		statusTasks := make([]*Task, 0)
		for _, task := range s.Tasks {
			if task.Status == status.key {
				statusTasks = append(statusTasks, task)
			}
		}

		if len(statusTasks) > 0 {
			fmt.Fprintf(file, "## %s\n\n", status.name)

			for _, task := range statusTasks {
				// Приоритет
				priorityEmoji := map[string]string{
					"low":    "🟢",
					"medium": "🟡",
					"high":   "🔴",
				}
				priority := priorityEmoji[task.Priority]
				if priority == "" {
					priority = "⚪"
				}

				fmt.Fprintf(file, "### %s %s\n\n", priority, task.Title)
				fmt.Fprintf(file, "**ID:** %d  \n", task.ID)
				fmt.Fprintf(file, "**Описание:** %s  \n", task.Description)
				fmt.Fprintf(file, "**Приоритет:** %s  \n", task.Priority)

				// Дедлайн
				if task.Deadline != nil {
					fmt.Fprintf(file, "**Дедлайн:** %s  \n", task.Deadline.Format("02.01.2006 15:04"))
				}

				// Теги
				if len(task.Tags) > 0 {
					tagStrings := make([]string, len(task.Tags))
					for i, tag := range task.Tags {
						tagStrings[i] = "`#" + tag + "`"
					}
					fmt.Fprintf(file, "**Теги:** %s  \n", strings.Join(tagStrings, ", "))
				}

				fmt.Fprintf(file, "**Создано:** %s  \n", task.CreatedAt.Format("02.01.2006 15:04"))
				fmt.Fprintf(file, "**Обновлено:** %s  \n", task.UpdatedAt.Format("02.01.2006 15:04"))
				fmt.Fprintln(file, "\n---\n")
			}
		}
	}

	return nil
}

// escapeCSV экранирует кавычки в строке для CSV
func escapeCSV(s string) string {
	return strings.ReplaceAll(s, "\"", "\"\"")
}
