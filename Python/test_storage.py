import unittest
import os
import csv
import tempfile
from datetime import datetime
from storage import Storage


class TestStorageAddGet(unittest.TestCase):
    def setUp(self):
        self.storage = Storage(filename=":memory_test:")  # won't hit disk
        self.storage.tasks = []
        self.storage.next_id = 1

    def test_add_task_returns_task(self):
        task = self.storage.add_task("Title", "Desc")
        self.assertEqual(task.title, "Title")
        self.assertEqual(task.description, "Desc")

    def test_add_increments_next_id(self):
        self.storage.add_task("A", "a")
        self.storage.add_task("B", "b")
        self.assertEqual(self.storage.next_id, 3)

    def test_get_task_by_id(self):
        task = self.storage.add_task("Find me", "Desc")
        found = self.storage.get_task(task.id)
        self.assertIs(found, task)

    def test_get_nonexistent_task_returns_none(self):
        self.assertIsNone(self.storage.get_task(999))

    def test_list_tasks_returns_all(self):
        self.storage.add_task("A", "a")
        self.storage.add_task("B", "b")
        self.assertEqual(len(self.storage.list_tasks()), 2)


class TestStorageDelete(unittest.TestCase):
    def setUp(self):
        self.storage = Storage(filename=":memory_test:")
        self.storage.tasks = []
        self.storage.next_id = 1

    def test_delete_existing_task(self):
        task = self.storage.add_task("Delete me", "Desc")
        result = self.storage.delete_task(task.id)
        self.assertTrue(result)
        self.assertIsNone(self.storage.get_task(task.id))

    def test_delete_nonexistent_returns_false(self):
        result = self.storage.delete_task(999)
        self.assertFalse(result)

    def test_delete_reduces_list(self):
        t1 = self.storage.add_task("A", "a")
        self.storage.add_task("B", "b")
        self.storage.delete_task(t1.id)
        self.assertEqual(len(self.storage.list_tasks()), 1)


class TestStorageFilter(unittest.TestCase):
    def setUp(self):
        self.storage = Storage(filename=":memory_test:")
        self.storage.tasks = []
        self.storage.next_id = 1
        t1 = self.storage.add_task("Task 1", "Desc 1")
        t1.update_status("done")
        t2 = self.storage.add_task("Task 2", "Desc 2")
        t2.update_status("in_progress")
        self.storage.add_task("Task 3", "Desc 3")  # stays todo

    def test_filter_by_status_done(self):
        result = self.storage.filter_tasks_by_status("done")
        self.assertEqual(len(result), 1)
        self.assertEqual(result[0].title, "Task 1")

    def test_filter_by_status_todo(self):
        result = self.storage.filter_tasks_by_status("todo")
        self.assertEqual(len(result), 1)

    def test_filter_by_status_no_match(self):
        result = self.storage.filter_tasks_by_status("nonexistent")
        self.assertEqual(result, [])

    def test_filter_by_tag(self):
        task = self.storage.add_task("Tagged", "Desc")
        task.add_tag("backend")
        result = self.storage.filter_tasks_by_tag("backend")
        self.assertEqual(len(result), 1)
        self.assertEqual(result[0].title, "Tagged")

    def test_filter_by_tag_case_insensitive(self):
        task = self.storage.add_task("Tagged", "Desc")
        task.add_tag("api")
        result = self.storage.filter_tasks_by_tag("API")
        self.assertEqual(len(result), 1)

    def test_filter_by_tag_no_match(self):
        result = self.storage.filter_tasks_by_tag("nonexistent")
        self.assertEqual(result, [])


class TestStorageSearch(unittest.TestCase):
    def setUp(self):
        self.storage = Storage(filename=":memory_test:")
        self.storage.tasks = []
        self.storage.next_id = 1
        self.storage.add_task("Fix login bug", "Users can't login")
        self.storage.add_task("Update docs", "Update the README file")
        self.storage.add_task("Deploy service", "Push to production server")

    def test_search_in_title(self):
        result = self.storage.search_tasks("login")
        self.assertEqual(len(result), 1)
        self.assertEqual(result[0].title, "Fix login bug")

    def test_search_in_description(self):
        result = self.storage.search_tasks("README")
        self.assertEqual(len(result), 1)

    def test_search_case_insensitive(self):
        result = self.storage.search_tasks("DEPLOY")
        self.assertEqual(len(result), 1)

    def test_search_no_match(self):
        result = self.storage.search_tasks("xyz_no_match")
        self.assertEqual(result, [])

    def test_search_matches_multiple(self):
        result = self.storage.search_tasks("e")  # appears in all titles
        self.assertGreater(len(result), 1)


class TestStorageGetAllTags(unittest.TestCase):
    def setUp(self):
        self.storage = Storage(filename=":memory_test:")
        self.storage.tasks = []
        self.storage.next_id = 1

    def test_no_tags_returns_empty(self):
        self.storage.add_task("A", "a")
        self.assertEqual(self.storage.get_all_tags(), [])

    def test_returns_unique_sorted_tags(self):
        t1 = self.storage.add_task("A", "a")
        t1.add_tag("zebra")
        t1.add_tag("apple")
        t2 = self.storage.add_task("B", "b")
        t2.add_tag("apple")  # duplicate
        tags = self.storage.get_all_tags()
        self.assertEqual(tags, ["apple", "zebra"])


class TestStorageSort(unittest.TestCase):
    def setUp(self):
        self.storage = Storage(filename=":memory_test:")
        self.storage.tasks = []
        self.storage.next_id = 1
        t1 = self.storage.add_task("A", "a")
        t1.update_priority("low")
        t2 = self.storage.add_task("B", "b")
        t2.update_priority("high")
        t3 = self.storage.add_task("C", "c")
        t3.update_priority("medium")
        t1.update_status("done")
        t2.update_status("in_progress")

    def test_sort_by_id(self):
        result = self.storage.sort_tasks("id")
        ids = [t.id for t in result]
        self.assertEqual(ids, sorted(ids))

    def test_sort_by_priority(self):
        result = self.storage.sort_tasks("priority")
        priorities = [t.priority for t in result]
        self.assertEqual(priorities[0], "high")
        self.assertEqual(priorities[-1], "low")

    def test_sort_by_status(self):
        result = self.storage.sort_tasks("status")
        statuses = [t.status for t in result]
        self.assertEqual(statuses, sorted(statuses))

    def test_sort_does_not_mutate_original(self):
        original_order = [t.id for t in self.storage.tasks]
        self.storage.sort_tasks("priority")
        self.assertEqual([t.id for t in self.storage.tasks], original_order)

    def test_sort_by_updated_reverse(self):
        result = self.storage.sort_tasks("updated")
        times = [t.updated_at for t in result]
        self.assertEqual(times, sorted(times, reverse=True))


class TestStoragePersistence(unittest.TestCase):
    def test_save_and_load_roundtrip(self):
        with tempfile.NamedTemporaryFile(suffix=".json", delete=False) as f:
            tmpfile = f.name

        try:
            s1 = Storage(filename=tmpfile)
            t = s1.add_task("Persist me", "Check persistence")
            t.update_priority("high")
            t.add_tag("test")
            s1.save()

            s2 = Storage(filename=tmpfile)
            s2.load()
            self.assertEqual(len(s2.list_tasks()), 1)
            loaded = s2.list_tasks()[0]
            self.assertEqual(loaded.title, "Persist me")
            self.assertEqual(loaded.priority, "high")
            self.assertIn("test", loaded.tags)
        finally:
            os.unlink(tmpfile)

    def test_load_nonexistent_file(self):
        s = Storage(filename="/nonexistent/path/tasks.json")
        s.load()  # should not raise
        self.assertEqual(s.list_tasks(), [])


class TestStorageExportCSV(unittest.TestCase):
    def test_export_csv_creates_file(self):
        with tempfile.NamedTemporaryFile(suffix=".csv", delete=False) as f:
            tmpfile = f.name

        try:
            storage = Storage(filename=":memory_test:")
            storage.tasks = []
            storage.next_id = 1
            t = storage.add_task("CSV Task", "CSV Desc")
            t.add_tag("export")

            result = storage.export_to_csv(tmpfile)
            self.assertTrue(result)
            self.assertTrue(os.path.exists(tmpfile))

            with open(tmpfile, encoding="utf-8", newline="") as f:
                reader = csv.reader(f)
                rows = list(reader)
            self.assertEqual(len(rows), 2)  # header + 1 task
            self.assertEqual(rows[0][0], "ID")
        finally:
            os.unlink(tmpfile)

    def test_export_csv_empty_storage(self):
        with tempfile.NamedTemporaryFile(suffix=".csv", delete=False) as f:
            tmpfile = f.name

        try:
            storage = Storage(filename=":memory_test:")
            storage.tasks = []
            result = storage.export_to_csv(tmpfile)
            self.assertTrue(result)
        finally:
            os.unlink(tmpfile)


class TestStorageExportMarkdown(unittest.TestCase):
    def test_export_markdown_creates_file(self):
        with tempfile.NamedTemporaryFile(suffix=".md", delete=False) as f:
            tmpfile = f.name

        try:
            storage = Storage(filename=":memory_test:")
            storage.tasks = []
            storage.next_id = 1
            t = storage.add_task("MD Task", "MD Desc")
            t.update_status("done")

            result = storage.export_to_markdown(tmpfile)
            self.assertTrue(result)

            with open(tmpfile, encoding="utf-8") as f:
                content = f.read()
            self.assertIn("# Список задач", content)
            self.assertIn("MD Task", content)
        finally:
            os.unlink(tmpfile)

    def test_export_markdown_empty_storage(self):
        with tempfile.NamedTemporaryFile(suffix=".md", delete=False) as f:
            tmpfile = f.name

        try:
            storage = Storage(filename=":memory_test:")
            storage.tasks = []
            result = storage.export_to_markdown(tmpfile)
            self.assertTrue(result)

            with open(tmpfile, encoding="utf-8") as f:
                content = f.read()
            self.assertIn("Задач нет", content)
        finally:
            os.unlink(tmpfile)


if __name__ == "__main__":
    unittest.main()
