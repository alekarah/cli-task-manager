import json
import os
from typing import List, Optional
from task import Task


class Storage:
    """Класс для управления хранением задач"""

    def __init__(self, filename: str = "tasks.json"):
        self.filename = filename
        self.tasks: List[Task] = []
        self.next_id = 1

    def load(self):
        """Загружает задачи из файла"""
        if not os.path.exists(self.filename):
            return

        try:
            with open(self.filename, 'r', encoding='utf-8') as f:
                data = json.load(f)
                self.next_id = data.get("next_id", 1)
                tasks_data = data.get("tasks", [])
                self.tasks = [Task.from_dict(task_data) for task_data in tasks_data]
        except (json.JSONDecodeError, KeyError) as e:
            print(f"Ошибка загрузки данных: {e}")
            self.tasks = []
            self.next_id = 1

    def save(self):
        """Сохраняет задачи в файл"""
        try:
            data = {
                "next_id": self.next_id,
                "tasks": [task.to_dict() for task in self.tasks]
            }
            with open(self.filename, 'w', encoding='utf-8') as f:
                json.dump(data, f, ensure_ascii=False, indent=2)
        except Exception as e:
            print(f"Ошибка сохранения данных: {e}")

    def add_task(self, title: str, description: str) -> Task:
        """Добавляет новую задачу"""
        task = Task(self.next_id, title, description)
        self.tasks.append(task)
        self.next_id += 1
        return task

    def get_task(self, task_id: int) -> Optional[Task]:
        """Получает задачу по ID"""
        for task in self.tasks:
            if task.id == task_id:
                return task
        return None

    def delete_task(self, task_id: int) -> bool:
        """Удаляет задачу по ID"""
        for i, task in enumerate(self.tasks):
            if task.id == task_id:
                self.tasks.pop(i)
                return True
        return False

    def list_tasks(self) -> List[Task]:
        """Возвращает все задачи"""
        return self.tasks
