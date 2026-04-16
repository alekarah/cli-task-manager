package main

import (
	"testing"
	"time"
)

func TestNewTask(t *testing.T) {
	task := NewTask(1, "Test Title", "Test Description")

	if task.ID != 1 {
		t.Errorf("expected ID=1, got %d", task.ID)
	}
	if task.Title != "Test Title" {
		t.Errorf("expected title 'Test Title', got %q", task.Title)
	}
	if task.Description != "Test Description" {
		t.Errorf("expected description 'Test Description', got %q", task.Description)
	}
	if task.Status != "todo" {
		t.Errorf("expected status 'todo', got %q", task.Status)
	}
	if task.Priority != "medium" {
		t.Errorf("expected priority 'medium', got %q", task.Priority)
	}
	if task.Deadline != nil {
		t.Error("expected nil deadline")
	}
	if task.Tags != nil && len(task.Tags) != 0 {
		t.Error("expected empty tags")
	}
}

func TestUpdateStatus(t *testing.T) {
	task := NewTask(1, "T", "D")

	statuses := []string{"todo", "in_progress", "done"}
	for _, s := range statuses {
		task.UpdateStatus(s)
		if task.Status != s {
			t.Errorf("expected status %q, got %q", s, task.Status)
		}
	}
}

func TestUpdateStatusUpdatesTimestamp(t *testing.T) {
	task := NewTask(1, "T", "D")
	before := task.UpdatedAt
	time.Sleep(time.Millisecond)
	task.UpdateStatus("done")
	if !task.UpdatedAt.After(before) {
		t.Error("UpdatedAt should be updated after UpdateStatus")
	}
}

func TestUpdate(t *testing.T) {
	task := NewTask(1, "Old Title", "Old Desc")

	task.Update("New Title", "")
	if task.Title != "New Title" {
		t.Errorf("expected 'New Title', got %q", task.Title)
	}
	if task.Description != "Old Desc" {
		t.Errorf("description should not change when empty string passed")
	}

	task.Update("", "New Desc")
	if task.Description != "New Desc" {
		t.Errorf("expected 'New Desc', got %q", task.Description)
	}
	if task.Title != "New Title" {
		t.Errorf("title should not change when empty string passed")
	}
}

func TestUpdatePriority(t *testing.T) {
	task := NewTask(1, "T", "D")

	for _, p := range []string{"low", "high", "medium"} {
		task.UpdatePriority(p)
		if task.Priority != p {
			t.Errorf("expected priority %q, got %q", p, task.Priority)
		}
	}
}

func TestUpdateDeadline(t *testing.T) {
	task := NewTask(1, "T", "D")

	dl := time.Date(2026, 12, 31, 23, 59, 0, 0, time.UTC)
	task.UpdateDeadline(&dl)
	if task.Deadline == nil || !task.Deadline.Equal(dl) {
		t.Error("deadline not set correctly")
	}

	task.UpdateDeadline(nil)
	if task.Deadline != nil {
		t.Error("deadline should be nil after clearing")
	}
}

func TestAddTag(t *testing.T) {
	task := NewTask(1, "T", "D")

	task.AddTag("backend")
	if len(task.Tags) != 1 || task.Tags[0] != "backend" {
		t.Errorf("expected tag 'backend', got %v", task.Tags)
	}
}

func TestAddTagLowercased(t *testing.T) {
	task := NewTask(1, "T", "D")
	task.AddTag("Backend")
	if task.Tags[0] != "backend" {
		t.Errorf("tag should be lowercased, got %q", task.Tags[0])
	}
}

func TestAddTagStripped(t *testing.T) {
	task := NewTask(1, "T", "D")
	task.AddTag("  api  ")
	if task.Tags[0] != "api" {
		t.Errorf("tag should be stripped, got %q", task.Tags[0])
	}
}

func TestAddTagNoDuplicates(t *testing.T) {
	task := NewTask(1, "T", "D")
	task.AddTag("api")
	task.AddTag("api")
	if len(task.Tags) != 1 {
		t.Errorf("expected 1 tag, got %d", len(task.Tags))
	}
}

func TestAddTagEmptyIgnored(t *testing.T) {
	task := NewTask(1, "T", "D")
	task.AddTag("")
	task.AddTag("   ")
	if len(task.Tags) != 0 {
		t.Errorf("empty tags should be ignored, got %v", task.Tags)
	}
}

func TestRemoveTag(t *testing.T) {
	task := NewTask(1, "T", "D")
	task.AddTag("api")
	task.RemoveTag("api")
	if len(task.Tags) != 0 {
		t.Errorf("tag should be removed, got %v", task.Tags)
	}
}

func TestRemoveNonexistentTagNoOp(t *testing.T) {
	task := NewTask(1, "T", "D")
	task.AddTag("api")
	task.RemoveTag("nope") // should not panic or change anything
	if len(task.Tags) != 1 {
		t.Errorf("tags should be unchanged, got %v", task.Tags)
	}
}

func TestSetTags(t *testing.T) {
	task := NewTask(1, "T", "D")
	task.AddTag("old")
	task.SetTags([]string{"new1", "new2"})
	if len(task.Tags) != 2 || task.Tags[0] != "new1" || task.Tags[1] != "new2" {
		t.Errorf("unexpected tags after SetTags: %v", task.Tags)
	}
}

func TestSetTagsLowercasedAndFiltered(t *testing.T) {
	task := NewTask(1, "T", "D")
	task.SetTags([]string{"TAG", "  ", "valid"})
	for _, tag := range task.Tags {
		if tag == "" || tag == "TAG" || tag == "  " {
			t.Errorf("unexpected tag value: %q", tag)
		}
	}
	found := false
	for _, tag := range task.Tags {
		if tag == "tag" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected lowercase 'tag', got %v", task.Tags)
	}
}

func TestAddTagUpdatesTimestamp(t *testing.T) {
	task := NewTask(1, "T", "D")
	before := task.UpdatedAt
	time.Sleep(time.Millisecond)
	task.AddTag("x")
	if !task.UpdatedAt.After(before) {
		t.Error("UpdatedAt should be updated after AddTag")
	}
}
