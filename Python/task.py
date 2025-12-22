from datetime import datetime
from typing import Optional


class Task:
    """Класс для представления задачи"""

    def __init__(self, task_id: int, title: str, description: str):
        self.id = task_id
        self.title = title
        self.description = description
        self.status = "todo"  # todo, in_progress, done
        self.created_at = datetime.now()
        self.updated_at = datetime.now()

    def update_status(self, status: str):
        """Обновляет статус задачи"""
        if status in ["todo", "in_progress", "done"]:
            self.status = status
            self.updated_at = datetime.now()
        else:
            raise ValueError("Некорректный статус")

    def update(self, title: Optional[str] = None, description: Optional[str] = None):
        """Обновляет поля задачи"""
        if title:
            self.title = title
        if description:
            self.description = description
        self.updated_at = datetime.now()

    def to_dict(self) -> dict:
        """Преобразует задачу в словарь"""
        return {
            "id": self.id,
            "title": self.title,
            "description": self.description,
            "status": self.status,
            "created_at": self.created_at.isoformat(),
            "updated_at": self.updated_at.isoformat()
        }

    @classmethod
    def from_dict(cls, data: dict) -> 'Task':
        """Создает задачу из словаря"""
        task = cls(data["id"], data["title"], data["description"])
        task.status = data["status"]

        # Обработка формата ISO с 'Z' (UTC)
        created_str = data["created_at"].replace('Z', '+00:00')
        updated_str = data["updated_at"].replace('Z', '+00:00')

        task.created_at = datetime.fromisoformat(created_str)
        task.updated_at = datetime.fromisoformat(updated_str)
        return task

    def __str__(self) -> str:
        status_emoji = {
            "todo": "📋",
            "in_progress": "⚙️",
            "done": "✅"
        }
        return f"""
ID: {self.id} {status_emoji.get(self.status, '❓')}
Название: {self.title}
Описание: {self.description}
Статус: {self.status}
Создано: {self.created_at.strftime('%d.%m.%Y %H:%M')}
Обновлено: {self.updated_at.strftime('%d.%m.%Y %H:%M')}
{'-' * 40}"""
