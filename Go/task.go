package main

import (
	"time"
)

// Task представляет структуру задачи
type Task struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`   // "todo", "in_progress", "done"
	Priority    string     `json:"priority"` // "low", "medium", "high"
	Deadline    *time.Time `json:"deadline,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// NewTask создает новую задачу
func NewTask(id int, title, description string) *Task {
	now := time.Now()
	return &Task{
		ID:          id,
		Title:       title,
		Description: description,
		Status:      "todo",
		Priority:    "medium",
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// UpdateStatus обновляет статус задачи
func (t *Task) UpdateStatus(status string) {
	t.Status = status
	t.UpdatedAt = time.Now()
}

// Update обновляет поля задачи
func (t *Task) Update(title, description string) {
	if title != "" {
		t.Title = title
	}
	if description != "" {
		t.Description = description
	}
	t.UpdatedAt = time.Now()
}

// UpdatePriority обновляет приоритет задачи
func (t *Task) UpdatePriority(priority string) {
	t.Priority = priority
	t.UpdatedAt = time.Now()
}

// UpdateDeadline обновляет дедлайн задачи
func (t *Task) UpdateDeadline(deadline *time.Time) {
	t.Deadline = deadline
	t.UpdatedAt = time.Now()
}
