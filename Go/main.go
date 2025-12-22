package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
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
		fmt.Println("3. edit - редактировать задачу")
		fmt.Println("4. status - изменить статус задачи")
		fmt.Println("5. delete - удалить задачу")
		fmt.Println("6. filter - фильтровать задачи по статусу")
		fmt.Println("7. search - поиск задач")
		fmt.Println("8. sort - сортировать задачи")
		fmt.Println("9. exit - выход")
		fmt.Print("\nВведите команду: ")

		command, _ := reader.ReadString('\n')
		command = strings.TrimSpace(command)

		switch command {
		case "list", "1":
			listTasks(storage)
		case "add", "2":
			addTask(storage, reader)
		case "edit", "3":
			editTask(storage, reader)
		case "status", "4":
			changeStatus(storage, reader)
		case "delete", "5":
			deleteTask(storage, reader)
		case "filter", "6":
			filterTasks(storage, reader)
		case "search", "7":
			searchTasks(storage, reader)
		case "sort", "8":
			sortTasks(storage, reader)
		case "exit", "9":
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
		statusEmoji := getStatusEmoji(task.Status)
		fmt.Printf("\nID: %d %s\n", task.ID, statusEmoji)
		fmt.Printf("Название: %s\n", task.Title)
		fmt.Printf("Описание: %s\n", task.Description)
		fmt.Printf("Статус: %s\n", task.Status)
		fmt.Printf("Создано: %s\n", task.CreatedAt.Format("02.01.2006 15:04"))
		fmt.Printf("Обновлено: %s\n", task.UpdatedAt.Format("02.01.2006 15:04"))
		fmt.Println(strings.Repeat("-", 40))
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

func editTask(storage *Storage, reader *bufio.Reader) {
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

	fmt.Printf("Текущее название: %s\n", task.Title)
	fmt.Print("Новое название (Enter - оставить без изменений): ")
	title, _ := reader.ReadString('\n')
	title = strings.TrimSpace(title)

	fmt.Printf("Текущее описание: %s\n", task.Description)
	fmt.Print("Новое описание (Enter - оставить без изменений): ")
	description, _ := reader.ReadString('\n')
	description = strings.TrimSpace(description)

	task.Update(title, description)

	if err := storage.Save(); err != nil {
		fmt.Printf("Ошибка сохранения: %v\n", err)
		return
	}

	fmt.Println("\n✓ Задача успешно обновлена!")
}

func changeStatus(storage *Storage, reader *bufio.Reader) {
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

	if err := storage.Save(); err != nil {
		fmt.Printf("Ошибка сохранения: %v\n", err)
		return
	}

	fmt.Println("\n✓ Статус задачи обновлен!")
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
		statusEmoji := getStatusEmoji(task.Status)
		fmt.Printf("\nID: %d %s\n", task.ID, statusEmoji)
		fmt.Printf("Название: %s\n", task.Title)
		fmt.Printf("Описание: %s\n", task.Description)
		fmt.Printf("Статус: %s\n", task.Status)
		fmt.Printf("Создано: %s\n", task.CreatedAt.Format("02.01.2006 15:04"))
		fmt.Printf("Обновлено: %s\n", task.UpdatedAt.Format("02.01.2006 15:04"))
		fmt.Println(strings.Repeat("-", 40))
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
		statusEmoji := getStatusEmoji(task.Status)
		fmt.Printf("\nID: %d %s\n", task.ID, statusEmoji)
		fmt.Printf("Название: %s\n", task.Title)
		fmt.Printf("Описание: %s\n", task.Description)
		fmt.Printf("Статус: %s\n", task.Status)
		fmt.Printf("Создано: %s\n", task.CreatedAt.Format("02.01.2006 15:04"))
		fmt.Printf("Обновлено: %s\n", task.UpdatedAt.Format("02.01.2006 15:04"))
		fmt.Println(strings.Repeat("-", 40))
	}
}

func sortTasks(storage *Storage, reader *bufio.Reader) {
	fmt.Println("\nДоступные варианты сортировки:")
	fmt.Println("1. id - по ID")
	fmt.Println("2. created - по дате создания (сначала старые)")
	fmt.Println("3. updated - по дате обновления (сначала новые)")
	fmt.Println("4. status - по статусу")
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
		"id":      "ID",
		"created": "дате создания",
		"updated": "дате обновления",
		"status":  "статусу",
	}

	fmt.Printf("\n=== Задачи, отсортированные по %s ===\n", sortNames[sortBy])
	for _, task := range tasks {
		statusEmoji := getStatusEmoji(task.Status)
		fmt.Printf("\nID: %d %s\n", task.ID, statusEmoji)
		fmt.Printf("Название: %s\n", task.Title)
		fmt.Printf("Описание: %s\n", task.Description)
		fmt.Printf("Статус: %s\n", task.Status)
		fmt.Printf("Создано: %s\n", task.CreatedAt.Format("02.01.2006 15:04"))
		fmt.Printf("Обновлено: %s\n", task.UpdatedAt.Format("02.01.2006 15:04"))
		fmt.Println(strings.Repeat("-", 40))
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
