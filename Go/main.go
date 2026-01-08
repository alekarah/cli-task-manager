package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	storage := NewStorage()

	// Загружаем существующие задачи
	if err := storage.Load(); err != nil {
		fmt.Printf("Ошибка загрузки данных: %v\n", err)
		return
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("=== Менеджер Задач ===")

	for {
		fmt.Println("\nКоманды:")
		fmt.Println("1. list - показать все задачи")
		fmt.Println("2. add - добавить задачу")
		fmt.Println("3. update - обновить задачу (название, статус, приоритет, дедлайн)")
		fmt.Println("4. delete - удалить задачу")
		fmt.Println("5. filter - фильтровать задачи по статусу")
		fmt.Println("6. search - поиск задач")
		fmt.Println("7. sort - сортировать задачи")
		fmt.Println("8. exit - выход")
		fmt.Print("\nВведите команду: ")

		command, _ := reader.ReadString('\n')
		command = strings.TrimSpace(command)

		switch command {
		case "list", "1":
			listTasks(storage)
		case "add", "2":
			addTask(storage, reader)
		case "update", "3":
			updateTask(storage, reader)
		case "delete", "4":
			deleteTask(storage, reader)
		case "filter", "5":
			filterTasks(storage, reader)
		case "search", "6":
			searchTasks(storage, reader)
		case "sort", "7":
			sortTasks(storage, reader)
		case "exit", "8":
			fmt.Println("До свидания!")
			return
		default:
			fmt.Println("Неизвестная команда!")
		}
	}
}

func listTasks(storage *Storage) {
	tasks := storage.ListTasks()

	if len(tasks) == 0 {
		fmt.Println("\nЗадач пока нет!")
		return
	}

	fmt.Println("\n=== Список задач ===")
	for _, task := range tasks {
		printTask(task)
	}
}

func addTask(storage *Storage, reader *bufio.Reader) {
	fmt.Print("\nНазвание задачи: ")
	title, _ := reader.ReadString('\n')
	title = strings.TrimSpace(title)

	if title == "" {
		fmt.Println("Название не может быть пустым!")
		return
	}

	fmt.Print("Описание задачи: ")
	description, _ := reader.ReadString('\n')
	description = strings.TrimSpace(description)

	task := storage.AddTask(title, description)

	if err := storage.Save(); err != nil {
		fmt.Printf("Ошибка сохранения: %v\n", err)
		return
	}

	fmt.Printf("\n✓ Задача #%d успешно создана!\n", task.ID)
}

func deleteTask(storage *Storage, reader *bufio.Reader) {
	fmt.Print("\nВведите ID задачи: ")
	idStr, _ := reader.ReadString('\n')
	id, err := strconv.Atoi(strings.TrimSpace(idStr))
	if err != nil {
		fmt.Println("Некорректный ID!")
		return
	}

	if storage.DeleteTask(id) {
		if err := storage.Save(); err != nil {
			fmt.Printf("Ошибка сохранения: %v\n", err)
			return
		}
		fmt.Println("\n✓ Задача успешно удалена!")
	} else {
		fmt.Println("Задача не найдена!")
	}
}

func updateTask(storage *Storage, reader *bufio.Reader) {
	fmt.Print("\nВведите ID задачи: ")
	idStr, _ := reader.ReadString('\n')
	id, err := strconv.Atoi(strings.TrimSpace(idStr))
	if err != nil {
		fmt.Println("Некорректный ID!")
		return
	}

	task := storage.GetTask(id)
	if task == nil {
		fmt.Println("Задача не найдена!")
		return
	}

	printTask(task)

	fmt.Println("\nЧто вы хотите обновить?")
	fmt.Println("1. Название и описание")
	fmt.Println("2. Статус")
	fmt.Println("3. Приоритет")
	fmt.Println("4. Дедлайн")
	fmt.Println("5. Всё сразу")
	fmt.Print("\nВыберите опцию: ")

	optionStr, _ := reader.ReadString('\n')
	option := strings.TrimSpace(optionStr)

	switch option {
	case "1":
		fmt.Printf("\nТекущее название: %s\n", task.Title)
		fmt.Print("Новое название (Enter - оставить без изменений): ")
		title, _ := reader.ReadString('\n')
		title = strings.TrimSpace(title)

		fmt.Printf("Текущее описание: %s\n", task.Description)
		fmt.Print("Новое описание (Enter - оставить без изменений): ")
		description, _ := reader.ReadString('\n')
		description = strings.TrimSpace(description)

		task.Update(title, description)

	case "2":
		fmt.Printf("\nТекущий статус: %s\n", task.Status)
		fmt.Println("\nДоступные статусы:")
		fmt.Println("1. todo - К выполнению")
		fmt.Println("2. in_progress - В процессе")
		fmt.Println("3. done - Выполнено")
		fmt.Print("\nВыберите статус: ")

		statusInput, _ := reader.ReadString('\n')
		statusInput = strings.TrimSpace(statusInput)

		var status string
		switch statusInput {
		case "1", "todo":
			status = "todo"
		case "2", "in_progress":
			status = "in_progress"
		case "3", "done":
			status = "done"
		default:
			fmt.Println("Некорректный статус!")
			return
		}

		task.UpdateStatus(status)

	case "3":
		fmt.Printf("\nТекущий приоритет: %s\n", task.Priority)
		fmt.Println("\nДоступные приоритеты:")
		fmt.Println("1. low - Низкий 🟢")
		fmt.Println("2. medium - Средний 🟡")
		fmt.Println("3. high - Высокий 🔴")
		fmt.Print("\nВыберите приоритет: ")

		priorityInput, _ := reader.ReadString('\n')
		priorityInput = strings.TrimSpace(priorityInput)

		var priority string
		switch priorityInput {
		case "1", "low":
			priority = "low"
		case "2", "medium":
			priority = "medium"
		case "3", "high":
			priority = "high"
		default:
			fmt.Println("Некорректный приоритет!")
			return
		}

		task.UpdatePriority(priority)

	case "4":
		// Дедлайн
		if task.Deadline != nil {
			fmt.Printf("\nТекущий дедлайн: %s\n", task.Deadline.Format("02.01.2006 15:04"))
		} else {
			fmt.Println("\nДедлайн не установлен")
		}
		fmt.Println("Введите новый дедлайн в формате ДД.ММ.ГГГГ ЧЧ:ММ")
		fmt.Println("Или просто ДД.ММ.ГГГГ (время будет 23:59)")
		fmt.Println("Или оставьте пустым для удаления дедлайна")
		fmt.Print("Дедлайн: ")

		deadlineInput, _ := reader.ReadString('\n')
		deadlineInput = strings.TrimSpace(deadlineInput)

		if deadlineInput == "" {
			task.UpdateDeadline(nil)
		} else {
			var deadline time.Time
			var err error

			// Пытаемся распарсить с временем
			if strings.Contains(deadlineInput, " ") {
				deadline, err = time.Parse("02.01.2006 15:04", deadlineInput)
			} else {
				// Если только дата, устанавливаем время 23:59
				deadline, err = time.Parse("02.01.2006", deadlineInput)
				if err == nil {
					deadline = time.Date(deadline.Year(), deadline.Month(), deadline.Day(), 23, 59, 0, 0, time.Local)
				}
			}

			if err != nil {
				fmt.Println("Некорректный формат даты!")
				return
			}
			task.UpdateDeadline(&deadline)
		}

	case "5":
		// Название и описание
		fmt.Printf("\nТекущее название: %s\n", task.Title)
		fmt.Print("Новое название (Enter - оставить без изменений): ")
		title, _ := reader.ReadString('\n')
		title = strings.TrimSpace(title)

		fmt.Printf("Текущее описание: %s\n", task.Description)
		fmt.Print("Новое описание (Enter - оставить без изменений): ")
		description, _ := reader.ReadString('\n')
		description = strings.TrimSpace(description)

		task.Update(title, description)

		// Статус
		fmt.Printf("\nТекущий статус: %s\n", task.Status)
		fmt.Println("Доступные статусы:")
		fmt.Println("1. todo")
		fmt.Println("2. in_progress")
		fmt.Println("3. done")
		fmt.Print("Выберите статус (Enter - оставить без изменений): ")

		statusInput, _ := reader.ReadString('\n')
		statusInput = strings.TrimSpace(statusInput)

		if statusInput != "" {
			var status string
			switch statusInput {
			case "1", "todo":
				status = "todo"
			case "2", "in_progress":
				status = "in_progress"
			case "3", "done":
				status = "done"
			default:
				fmt.Println("Некорректный статус!")
				return
			}
			task.UpdateStatus(status)
		}

		// Приоритет
		fmt.Printf("\nТекущий приоритет: %s\n", task.Priority)
		fmt.Println("Доступные приоритеты:")
		fmt.Println("1. low 🟢")
		fmt.Println("2. medium 🟡")
		fmt.Println("3. high 🔴")
		fmt.Print("Выберите приоритет (Enter - оставить без изменений): ")

		priorityInput, _ := reader.ReadString('\n')
		priorityInput = strings.TrimSpace(priorityInput)

		if priorityInput != "" {
			var priority string
			switch priorityInput {
			case "1", "low":
				priority = "low"
			case "2", "medium":
				priority = "medium"
			case "3", "high":
				priority = "high"
			default:
				fmt.Println("Некорректный приоритет!")
				return
			}
			task.UpdatePriority(priority)
		}

		// Дедлайн
		if task.Deadline != nil {
			fmt.Printf("\nТекущий дедлайн: %s\n", task.Deadline.Format("02.01.2006 15:04"))
		} else {
			fmt.Println("\nДедлайн не установлен")
		}
		fmt.Println("Формат: ДД.ММ.ГГГГ ЧЧ:ММ или ДД.ММ.ГГГГ")
		fmt.Print("Дедлайн (Enter - оставить без изменений): ")

		deadlineInput, _ := reader.ReadString('\n')
		deadlineInput = strings.TrimSpace(deadlineInput)

		if deadlineInput != "" {
			if strings.ToLower(deadlineInput) == "удалить" {
				task.UpdateDeadline(nil)
			} else {
				var deadline time.Time
				var err error

				if strings.Contains(deadlineInput, " ") {
					deadline, err = time.Parse("02.01.2006 15:04", deadlineInput)
				} else {
					deadline, err = time.Parse("02.01.2006", deadlineInput)
					if err == nil {
						deadline = time.Date(deadline.Year(), deadline.Month(), deadline.Day(), 23, 59, 0, 0, time.Local)
					}
				}

				if err != nil {
					fmt.Println("Некорректный формат даты! Дедлайн не обновлен.")
				} else {
					task.UpdateDeadline(&deadline)
				}
			}
		}

	default:
		fmt.Println("Некорректная опция!")
		return
	}

	if err := storage.Save(); err != nil {
		fmt.Printf("Ошибка сохранения: %v\n", err)
		return
	}

	fmt.Println("\n✓ Задача успешно обновлена!")
}

func filterTasks(storage *Storage, reader *bufio.Reader) {
	fmt.Println("\nДоступные статусы:")
	fmt.Println("1. todo - К выполнению")
	fmt.Println("2. in_progress - В процессе")
	fmt.Println("3. done - Выполнено")
	fmt.Print("\nВыберите статус для фильтрации: ")

	statusInput, _ := reader.ReadString('\n')
	statusInput = strings.TrimSpace(statusInput)

	var status string
	switch statusInput {
	case "1", "todo":
		status = "todo"
	case "2", "in_progress":
		status = "in_progress"
	case "3", "done":
		status = "done"
	default:
		fmt.Println("Некорректный статус!")
		return
	}

	tasks := storage.FilterTasksByStatus(status)

	if len(tasks) == 0 {
		fmt.Printf("\nЗадач со статусом '%s' не найдено!\n", status)
		return
	}

	fmt.Printf("\n=== Задачи со статусом '%s' ===\n", status)
	for _, task := range tasks {
		printTask(task)
	}
}

func searchTasks(storage *Storage, reader *bufio.Reader) {
	fmt.Print("\nВведите текст для поиска: ")
	query, _ := reader.ReadString('\n')
	query = strings.TrimSpace(query)

	if query == "" {
		fmt.Println("Запрос не может быть пустым!")
		return
	}

	tasks := storage.SearchTasks(query)

	if len(tasks) == 0 {
		fmt.Printf("\nЗадачи, содержащие '%s', не найдены!\n", query)
		return
	}

	fmt.Printf("\n=== Результаты поиска для '%s' ===\n", query)
	for _, task := range tasks {
		printTask(task)
	}
}

func sortTasks(storage *Storage, reader *bufio.Reader) {
	fmt.Println("\nДоступные варианты сортировки:")
	fmt.Println("1. id - по ID")
	fmt.Println("2. created - по дате создания (сначала старые)")
	fmt.Println("3. updated - по дате обновления (сначала новые)")
	fmt.Println("4. status - по статусу")
	fmt.Println("5. priority - по приоритету (сначала высокий)")
	fmt.Print("\nВыберите вариант сортировки: ")

	sortInput, _ := reader.ReadString('\n')
	sortInput = strings.TrimSpace(sortInput)

	var sortBy string
	switch sortInput {
	case "1", "id":
		sortBy = "id"
	case "2", "created":
		sortBy = "created"
	case "3", "updated":
		sortBy = "updated"
	case "4", "status":
		sortBy = "status"
	case "5", "priority":
		sortBy = "priority"
	default:
		fmt.Println("Некорректный вариант сортировки!")
		return
	}

	tasks := storage.SortTasks(sortBy)

	if len(tasks) == 0 {
		fmt.Println("\nЗадач пока нет!")
		return
	}

	sortNames := map[string]string{
		"id":       "ID",
		"created":  "дате создания",
		"updated":  "дате обновления",
		"status":   "статусу",
		"priority": "приоритету",
	}

	fmt.Printf("\n=== Задачи, отсортированные по %s ===\n", sortNames[sortBy])
	for _, task := range tasks {
		printTask(task)
	}
}

func getStatusEmoji(status string) string {
	switch status {
	case "todo":
		return "📋"
	case "in_progress":
		return "⚙️"
	case "done":
		return "✅"
	default:
		return "❓"
	}
}

func getPriorityEmoji(priority string) string {
	switch priority {
	case "low":
		return "🟢"
	case "medium":
		return "🟡"
	case "high":
		return "🔴"
	default:
		return "⚪"
	}
}

func printTask(task *Task) {
	statusEmoji := getStatusEmoji(task.Status)
	priorityEmoji := getPriorityEmoji(task.Priority)
	fmt.Printf("\nID: %d %s %s\n", task.ID, statusEmoji, priorityEmoji)
	fmt.Printf("Название: %s\n", task.Title)
	fmt.Printf("Описание: %s\n", task.Description)
	fmt.Printf("Статус: %s\n", task.Status)
	fmt.Printf("Приоритет: %s\n", task.Priority)

	// Отображение дедлайна с проверкой на просрочку
	if task.Deadline != nil {
		deadlineFormatted := task.Deadline.Format("02.01.2006 15:04")
		if task.Deadline.Before(time.Now()) && task.Status != "done" {
			fmt.Printf("Дедлайн: %s ⏰ ПРОСРОЧЕН!\n", deadlineFormatted)
		} else {
			fmt.Printf("Дедлайн: %s\n", deadlineFormatted)
		}
	}

	fmt.Printf("Создано: %s\n", task.CreatedAt.Format("02.01.2006 15:04"))
	fmt.Printf("Обновлено: %s\n", task.UpdatedAt.Format("02.01.2006 15:04"))
	fmt.Println(strings.Repeat("-", 40))
}
