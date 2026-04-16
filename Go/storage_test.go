package main

import (
	"encoding/json"
	"os"
	"testing"
	"time"
)

func TestAddTask(t *testing.T) {
	s := NewStorage()
	task := s.AddTask("Title", "Desc")
	if task.Title != "Title" || task.Description != "Desc" {
		t.Errorf("unexpected task: %+v", task)
	}
	if task.ID != 1 {
		t.Errorf("expected ID=1, got %d", task.ID)
	}
}

func TestAddTaskIncrementsID(t *testing.T) {
	s := NewStorage()
	t1 := s.AddTask("A", "a")
	t2 := s.AddTask("B", "b")
	if t2.ID != t1.ID+1 {
		t.Errorf("IDs should increment: %d, %d", t1.ID, t2.ID)
	}
}

func TestGetTask(t *testing.T) {
	s := NewStorage()
	task := s.AddTask("Find me", "Desc")
	found := s.GetTask(task.ID)
	if found == nil || found.Title != "Find me" {
		t.Errorf("GetTask failed, got %v", found)
	}
}

func TestGetTaskNotFound(t *testing.T) {
	s := NewStorage()
	if s.GetTask(999) != nil {
		t.Error("expected nil for missing task")
	}
}

func TestDeleteTask(t *testing.T) {
	s := NewStorage()
	task := s.AddTask("Delete me", "Desc")
	ok := s.DeleteTask(task.ID)
	if !ok {
		t.Error("expected true from DeleteTask")
	}
	if s.GetTask(task.ID) != nil {
		t.Error("task should be gone after delete")
	}
}

func TestDeleteTaskNotFound(t *testing.T) {
	s := NewStorage()
	if s.DeleteTask(999) {
		t.Error("expected false for missing task")
	}
}

func TestListTasks(t *testing.T) {
	s := NewStorage()
	s.AddTask("A", "a")
	s.AddTask("B", "b")
	if len(s.ListTasks()) != 2 {
		t.Errorf("expected 2 tasks, got %d", len(s.ListTasks()))
	}
}

func TestFilterTasksByStatus(t *testing.T) {
	s := NewStorage()
	t1 := s.AddTask("A", "a")
	t1.UpdateStatus("done")
	s.AddTask("B", "b") // stays todo

	done := s.FilterTasksByStatus("done")
	if len(done) != 1 || done[0].Title != "A" {
		t.Errorf("unexpected filter result: %v", done)
	}

	todo := s.FilterTasksByStatus("todo")
	if len(todo) != 1 {
		t.Errorf("expected 1 todo task, got %d", len(todo))
	}
}

func TestFilterTasksByStatusNoMatch(t *testing.T) {
	s := NewStorage()
	s.AddTask("A", "a")
	result := s.FilterTasksByStatus("nonexistent")
	if len(result) != 0 {
		t.Errorf("expected empty result, got %v", result)
	}
}

func TestSearchTasks(t *testing.T) {
	s := NewStorage()
	s.AddTask("Fix login bug", "Users cannot login")
	s.AddTask("Update README", "Update docs")

	result := s.SearchTasks("login")
	if len(result) != 1 || result[0].Title != "Fix login bug" {
		t.Errorf("unexpected search result: %v", result)
	}
}

func TestSearchTasksCaseInsensitive(t *testing.T) {
	s := NewStorage()
	s.AddTask("Fix Login Bug", "Desc")
	result := s.SearchTasks("login")
	if len(result) != 1 {
		t.Errorf("search should be case-insensitive, got %d results", len(result))
	}
}

func TestSearchTasksInDescription(t *testing.T) {
	s := NewStorage()
	s.AddTask("Task A", "deploy to production")
	result := s.SearchTasks("production")
	if len(result) != 1 {
		t.Errorf("should find task by description, got %d results", len(result))
	}
}

func TestSearchTasksNoMatch(t *testing.T) {
	s := NewStorage()
	s.AddTask("A", "a")
	result := s.SearchTasks("xyz_no_match")
	if len(result) != 0 {
		t.Errorf("expected no results, got %v", result)
	}
}

func TestFilterTasksByTag(t *testing.T) {
	s := NewStorage()
	t1 := s.AddTask("Tagged", "Desc")
	t1.AddTag("backend")
	s.AddTask("No tag", "Desc")

	result := s.FilterTasksByTag("backend")
	if len(result) != 1 || result[0].Title != "Tagged" {
		t.Errorf("unexpected tag filter result: %v", result)
	}
}

func TestFilterTasksByTagCaseInsensitive(t *testing.T) {
	s := NewStorage()
	t1 := s.AddTask("Tagged", "Desc")
	t1.AddTag("api")
	result := s.FilterTasksByTag("API")
	if len(result) != 1 {
		t.Errorf("tag filter should be case-insensitive")
	}
}

func TestGetAllTags(t *testing.T) {
	s := NewStorage()
	t1 := s.AddTask("A", "a")
	t1.AddTag("zebra")
	t1.AddTag("apple")
	t2 := s.AddTask("B", "b")
	t2.AddTag("apple") // duplicate

	tags := s.GetAllTags()
	if len(tags) != 2 || tags[0] != "apple" || tags[1] != "zebra" {
		t.Errorf("expected [apple zebra], got %v", tags)
	}
}

func TestGetAllTagsEmpty(t *testing.T) {
	s := NewStorage()
	s.AddTask("A", "a")
	if len(s.GetAllTags()) != 0 {
		t.Error("expected no tags")
	}
}

func TestSortTasksByID(t *testing.T) {
	s := NewStorage()
	s.AddTask("C", "c")
	s.AddTask("A", "a")
	s.AddTask("B", "b")

	result := s.SortTasks("id")
	for i := 1; i < len(result); i++ {
		if result[i].ID < result[i-1].ID {
			t.Errorf("tasks not sorted by id at index %d", i)
		}
	}
}

func TestSortTasksByPriority(t *testing.T) {
	s := NewStorage()
	t1 := s.AddTask("Low", "l")
	t1.UpdatePriority("low")
	t2 := s.AddTask("High", "h")
	t2.UpdatePriority("high")
	t3 := s.AddTask("Med", "m")
	t3.UpdatePriority("medium")

	result := s.SortTasks("priority")
	if result[0].Priority != "high" || result[len(result)-1].Priority != "low" {
		t.Errorf("unexpected priority sort: %v %v %v",
			result[0].Priority, result[1].Priority, result[2].Priority)
	}
}

func TestSortTasksByUpdatedReversed(t *testing.T) {
	s := NewStorage()
	t1 := s.AddTask("A", "a")
	time.Sleep(time.Millisecond)
	t2 := s.AddTask("B", "b")
	time.Sleep(time.Millisecond)
	t2.UpdateStatus("done") // make t2 most recently updated
	_ = t1

	result := s.SortTasks("updated")
	if result[0].Title != "B" {
		t.Errorf("most recently updated should be first, got %q", result[0].Title)
	}
}

func TestSortDoesNotMutateOriginal(t *testing.T) {
	s := NewStorage()
	s.AddTask("A", "a")
	s.AddTask("B", "b")
	original := []int{s.Tasks[0].ID, s.Tasks[1].ID}

	s.SortTasks("priority")
	after := []int{s.Tasks[0].ID, s.Tasks[1].ID}
	if original[0] != after[0] || original[1] != after[1] {
		t.Error("SortTasks should not mutate original Tasks slice")
	}
}

func TestSaveAndLoad(t *testing.T) {
	f, err := os.CreateTemp("", "tasks_persist_*.json")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	tmpFile := f.Name()
	defer os.Remove(tmpFile)

	// Patch dataFile via a helper — since Go uses a package-level const,
	// we write JSON directly and load it.
	s1 := NewStorage()
	task := s1.AddTask("Persist me", "Check persistence")
	task.UpdatePriority("high")
	task.AddTag("test")

	data, err := json.MarshalIndent(s1, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(tmpFile, data, 0644); err != nil {
		t.Fatal(err)
	}

	// Load into a new storage
	rawData, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatal(err)
	}
	s2 := NewStorage()
	if err := json.Unmarshal(rawData, s2); err != nil {
		t.Fatal(err)
	}

	if len(s2.Tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(s2.Tasks))
	}
	loaded := s2.Tasks[0]
	if loaded.Title != "Persist me" {
		t.Errorf("wrong title: %q", loaded.Title)
	}
	if loaded.Priority != "high" {
		t.Errorf("wrong priority: %q", loaded.Priority)
	}
	if len(loaded.Tags) != 1 || loaded.Tags[0] != "test" {
		t.Errorf("wrong tags: %v", loaded.Tags)
	}
}

func TestExportToCSV(t *testing.T) {
	f, err := os.CreateTemp("", "export_*.csv")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	tmpFile := f.Name()
	defer os.Remove(tmpFile)

	s := NewStorage()
	t1 := s.AddTask("CSV Task", "CSV Desc")
	t1.AddTag("export")

	if err := s.ExportToCSV(tmpFile); err != nil {
		t.Fatalf("ExportToCSV failed: %v", err)
	}

	data, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatal(err)
	}
	content := string(data)
	if len(content) == 0 {
		t.Error("CSV file is empty")
	}
	// BOM + header expected
	if content[3:5] != "ID" {
		t.Errorf("expected header starting with ID (after BOM), got %q", content[:20])
	}
}

func TestExportToMarkdown(t *testing.T) {
	f, err := os.CreateTemp("", "export_*.md")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	tmpFile := f.Name()
	defer os.Remove(tmpFile)

	s := NewStorage()
	t1 := s.AddTask("MD Task", "MD Desc")
	t1.UpdateStatus("done")

	if err := s.ExportToMarkdown(tmpFile); err != nil {
		t.Fatalf("ExportToMarkdown failed: %v", err)
	}

	data, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatal(err)
	}
	content := string(data)
	if len(content) == 0 {
		t.Error("Markdown file is empty")
	}
	for _, expected := range []string{"# Список задач", "MD Task"} {
		if !containsString(content, expected) {
			t.Errorf("expected %q in markdown output", expected)
		}
	}
}

func TestExportToMarkdownEmpty(t *testing.T) {
	f, err := os.CreateTemp("", "export_empty_*.md")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	tmpFile := f.Name()
	defer os.Remove(tmpFile)

	s := NewStorage()
	if err := s.ExportToMarkdown(tmpFile); err != nil {
		t.Fatal(err)
	}

	data, _ := os.ReadFile(tmpFile)
	if !containsString(string(data), "Задач нет") {
		t.Error("expected 'Задач нет' in empty export")
	}
}

func TestEscapeCSV(t *testing.T) {
	cases := []struct{ input, want string }{
		{`no quotes`, `no quotes`},
		{`say "hello"`, `say ""hello""`},
		{"\"\"", "\"\"\"\""},
	}
	for _, c := range cases {
		got := escapeCSV(c.input)
		if got != c.want {
			t.Errorf("escapeCSV(%q) = %q, want %q", c.input, got, c.want)
		}
	}
}

func containsString(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 ||
		func() bool {
			for i := 0; i <= len(s)-len(sub); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
			return false
		}())
}
