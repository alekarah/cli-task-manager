from datetime import datetime
from typing import Optional


class Task:
    """Класс для представления задачи"""

    def __init__(self, task_id: int, title: str, description: str):
        self.id = task_id
        self.title = title
        self.description = description
        self.status = "todo"  # todo, in_progress, done
        self.priority = "medium"  # low, medium, high
        self.deadline: Optional[datetime] = None
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

    def update_priority(self, priority: str):
        """Обновляет приоритет задачи"""
        if priority in ["low", "medium", "high"]:
            self.priority = priority
            self.updated_at = datetime.now()
        else:
            raise ValueError("Некорректный приоритет")

    def update_deadline(self, deadline: Optional[datetime]):
        """Обновляет дедлайн задачи"""
        self.deadline = deadline
        self.updated_at = datetime.now()

    def to_dict(self) -> dict:
        """Преобразует задачу в словарь"""
        return {
            "id": self.id,
            "title": self.title,
            "description": self.description,
            "status": self.status,
            "priority": self.priority,
            "deadline": self.deadline.isoformat() if self.deadline else None,
            "created_at": self.created_at.isoformat(),
            "updated_at": self.updated_at.isoformat()
        }

    @classmethod
    def from_dict(cls, data: dict) -> 'Task':
        """Создает задачу из словаря"""
        task = cls(data["id"], data["title"], data["description"])
        task.status = data["status"]
        task.priority = data.get("priority", "medium")  # По умолчанию medium для старых задач

        # Обработка формата ISO с 'Z' (UTC)
        created_str = data["created_at"].replace('Z', '+00:00')
        updated_str = data["updated_at"].replace('Z', '+00:00')

        task.created_at = datetime.fromisoformat(created_str)
        task.updated_at = datetime.fromisoformat(updated_str)

        # Обработка дедлайна
        if data.get("deadline"):
            deadline_str = data["deadline"].replace('Z', '+00:00')
            task.deadline = datetime.fromisoformat(deadline_str)

        return task

    def __str__(self) -> str:
        status_emoji = {
            "todo": "📋",
            "in_progress": "⚙️",
            "done": "✅"
        }
        priority_emoji = {
            "low": "🟢",
            "medium": "🟡",
            "high": "🔴"
        }

        # Формирование строки дедлайна с проверкой на просрочку
        deadline_str = ""
        if self.deadline:
            # Убираем timezone информацию для сравнения
            now = datetime.now()
            deadline_naive = self.deadline.replace(tzinfo=None) if self.deadline.tzinfo else self.deadline
            deadline_formatted = deadline_naive.strftime('%d.%m.%Y %H:%M')
            if deadline_naive < now and self.status != "done":
                deadline_str = f"\nДедлайн: {deadline_formatted} ⏰ ПРОСРОЧЕН!"
            else:
                deadline_str = f"\nДедлайн: {deadline_formatted}"

        return f"""
ID: {self.id} {status_emoji.get(self.status, '❓')} {priority_emoji.get(self.priority, '⚪')}
Название: {self.title}
Описание: {self.description}
Статус: {self.status}
Приоритет: {self.priority}{deadline_str}
Создано: {self.created_at.strftime('%d.%m.%Y %H:%M')}
Обновлено: {self.updated_at.strftime('%d.%m.%Y %H:%M')}
{'-' * 40}"""
