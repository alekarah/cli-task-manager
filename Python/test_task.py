import unittest
from datetime import datetime, timedelta
from task import Task


class TestTaskCreation(unittest.TestCase):
    def test_default_values(self):
        task = Task(1, "Test", "Description")
        self.assertEqual(task.id, 1)
        self.assertEqual(task.title, "Test")
        self.assertEqual(task.description, "Description")
        self.assertEqual(task.status, "todo")
        self.assertEqual(task.priority, "medium")
        self.assertIsNone(task.deadline)
        self.assertEqual(task.tags, [])
        self.assertIsNotNone(task.created_at)
        self.assertIsNotNone(task.updated_at)


class TestUpdateStatus(unittest.TestCase):
    def setUp(self):
        self.task = Task(1, "Test", "Desc")

    def test_valid_statuses(self):
        for status in ["todo", "in_progress", "done"]:
            self.task.update_status(status)
            self.assertEqual(self.task.status, status)

    def test_invalid_status_raises(self):
        with self.assertRaises(ValueError):
            self.task.update_status("invalid")

    def test_updates_updated_at(self):
        before = self.task.updated_at
        self.task.update_status("done")
        self.assertGreaterEqual(self.task.updated_at, before)


class TestUpdateFields(unittest.TestCase):
    def setUp(self):
        self.task = Task(1, "Old Title", "Old Desc")

    def test_update_title(self):
        self.task.update(title="New Title")
        self.assertEqual(self.task.title, "New Title")
        self.assertEqual(self.task.description, "Old Desc")

    def test_update_description(self):
        self.task.update(description="New Desc")
        self.assertEqual(self.task.description, "New Desc")
        self.assertEqual(self.task.title, "Old Title")

    def test_update_both(self):
        self.task.update(title="T", description="D")
        self.assertEqual(self.task.title, "T")
        self.assertEqual(self.task.description, "D")

    def test_empty_string_does_not_overwrite(self):
        self.task.update(title="", description="")
        self.assertEqual(self.task.title, "Old Title")
        self.assertEqual(self.task.description, "Old Desc")


class TestUpdatePriority(unittest.TestCase):
    def setUp(self):
        self.task = Task(1, "Test", "Desc")

    def test_valid_priorities(self):
        for p in ["low", "medium", "high"]:
            self.task.update_priority(p)
            self.assertEqual(self.task.priority, p)

    def test_invalid_priority_raises(self):
        with self.assertRaises(ValueError):
            self.task.update_priority("urgent")


class TestDeadline(unittest.TestCase):
    def setUp(self):
        self.task = Task(1, "Test", "Desc")

    def test_set_deadline(self):
        dl = datetime(2026, 12, 31, 23, 59)
        self.task.update_deadline(dl)
        self.assertEqual(self.task.deadline, dl)

    def test_clear_deadline(self):
        self.task.update_deadline(datetime(2026, 1, 1))
        self.task.update_deadline(None)
        self.assertIsNone(self.task.deadline)


class TestTags(unittest.TestCase):
    def setUp(self):
        self.task = Task(1, "Test", "Desc")

    def test_add_tag(self):
        self.task.add_tag("backend")
        self.assertIn("backend", self.task.tags)

    def test_add_tag_lowercased(self):
        self.task.add_tag("Backend")
        self.assertIn("backend", self.task.tags)
        self.assertNotIn("Backend", self.task.tags)

    def test_add_tag_stripped(self):
        self.task.add_tag("  api  ")
        self.assertIn("api", self.task.tags)

    def test_no_duplicate_tags(self):
        self.task.add_tag("api")
        self.task.add_tag("api")
        self.assertEqual(self.task.tags.count("api"), 1)

    def test_add_empty_tag_ignored(self):
        self.task.add_tag("")
        self.assertEqual(self.task.tags, [])

    def test_remove_tag(self):
        self.task.add_tag("api")
        self.task.remove_tag("api")
        self.assertNotIn("api", self.task.tags)

    def test_remove_nonexistent_tag(self):
        self.task.remove_tag("nope")  # должно не упасть

    def test_set_tags_replaces_all(self):
        self.task.add_tag("old")
        self.task.set_tags(["new1", "new2"])
        self.assertEqual(self.task.tags, ["new1", "new2"])

    def test_set_tags_lowercased(self):
        self.task.set_tags(["TAG"])
        self.assertEqual(self.task.tags, ["tag"])

    def test_set_tags_filters_empty(self):
        self.task.set_tags(["  ", "valid"])
        self.assertNotIn("", self.task.tags)
        self.assertIn("valid", self.task.tags)


class TestSerialization(unittest.TestCase):
    def test_to_dict_and_from_dict_roundtrip(self):
        task = Task(5, "Round Trip", "Test serialization")
        task.update_status("in_progress")
        task.update_priority("high")
        task.update_deadline(datetime(2026, 6, 15, 12, 0))
        task.add_tag("api")
        task.add_tag("backend")

        data = task.to_dict()
        restored = Task.from_dict(data)

        self.assertEqual(restored.id, task.id)
        self.assertEqual(restored.title, task.title)
        self.assertEqual(restored.description, task.description)
        self.assertEqual(restored.status, task.status)
        self.assertEqual(restored.priority, task.priority)
        self.assertEqual(restored.tags, task.tags)

    def test_from_dict_default_priority(self):
        data = {
            "id": 1, "title": "T", "description": "D",
            "status": "todo",
            "created_at": "2026-01-01T00:00:00",
            "updated_at": "2026-01-01T00:00:00"
        }
        task = Task.from_dict(data)
        self.assertEqual(task.priority, "medium")

    def test_from_dict_z_suffix_dates(self):
        data = {
            "id": 1, "title": "T", "description": "D",
            "status": "todo", "priority": "low",
            "created_at": "2026-01-01T00:00:00Z",
            "updated_at": "2026-01-01T00:00:00Z"
        }
        task = Task.from_dict(data)
        self.assertIsNotNone(task.created_at)

    def test_from_dict_no_deadline(self):
        data = {
            "id": 2, "title": "T", "description": "D",
            "status": "done", "priority": "high",
            "created_at": "2026-03-01T10:00:00",
            "updated_at": "2026-03-02T10:00:00"
        }
        task = Task.from_dict(data)
        self.assertIsNone(task.deadline)


if __name__ == "__main__":
    unittest.main()
