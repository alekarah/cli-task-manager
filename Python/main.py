#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import sys
import io

# Настройка UTF-8 для Windows консоли
if sys.platform == 'win32':
    sys.stdout = io.TextIOWrapper(sys.stdout.buffer, encoding='utf-8')
    sys.stderr = io.TextIOWrapper(sys.stderr.buffer, encoding='utf-8')

from storage import Storage


def list_tasks(storage: Storage):
    """Показывает список всех задач"""
    tasks = storage.list_tasks()

    if not tasks:
        print("\nЗадач пока нет!")
        return

    print("\n=== Список задач ===")
    for task in tasks:
        print(task)


def add_task(storage: Storage):
    """Добавляет новую задачу"""
    print()
    title = input("Название задачи: ").strip()

    if not title:
        print("Название не может быть пустым!")
        return

    description = input("Описание задачи: ").strip()

    task = storage.add_task(title, description)
    storage.save()

    print(f"\n✓ Задача #{task.id} успешно создана!")


def edit_task(storage: Storage):
    """Редактирует существующую задачу"""
    print()
    try:
        task_id = int(input("Введите ID задачи: ").strip())
    except ValueError:
        print("Некорректный ID!")
        return

    task = storage.get_task(task_id)
    if not task:
        print("Задача не найдена!")
        return

    print(f"Текущее название: {task.title}")
    title = input("Новое название (Enter - оставить без изменений): ").strip()

    print(f"Текущее описание: {task.description}")
    description = input("Новое описание (Enter - оставить без изменений): ").strip()

    task.update(title if title else None, description if description else None)
    storage.save()

    print("\n✓ Задача успешно обновлена!")


def change_status(storage: Storage):
    """Изменяет статус задачи"""
    print()
    try:
        task_id = int(input("Введите ID задачи: ").strip())
    except ValueError:
        print("Некорректный ID!")
        return

    task = storage.get_task(task_id)
    if not task:
        print("Задача не найдена!")
        return

    print("\nДоступные статусы:")
    print("1. todo - К выполнению")
    print("2. in_progress - В процессе")
    print("3. done - Выполнено")

    status_input = input("\nВыберите статус: ").strip()

    status_map = {
        "1": "todo",
        "2": "in_progress",
        "3": "done",
        "todo": "todo",
        "in_progress": "in_progress",
        "done": "done"
    }

    status = status_map.get(status_input)
    if not status:
        print("Некорректный статус!")
        return

    try:
        task.update_status(status)
        storage.save()
        print("\n✓ Статус задачи обновлен!")
    except ValueError as e:
        print(f"Ошибка: {e}")


def delete_task(storage: Storage):
    """Удаляет задачу"""
    print()
    try:
        task_id = int(input("Введите ID задачи: ").strip())
    except ValueError:
        print("Некорректный ID!")
        return

    if storage.delete_task(task_id):
        storage.save()
        print("\n✓ Задача успешно удалена!")
    else:
        print("Задача не найдена!")


def filter_tasks(storage: Storage):
    """Фильтрует задачи по статусу"""
    print("\nДоступные статусы:")
    print("1. todo - К выполнению")
    print("2. in_progress - В процессе")
    print("3. done - Выполнено")

    status_input = input("\nВыберите статус для фильтрации: ").strip()

    status_map = {
        "1": "todo",
        "2": "in_progress",
        "3": "done",
        "todo": "todo",
        "in_progress": "in_progress",
        "done": "done"
    }

    status = status_map.get(status_input)
    if not status:
        print("Некорректный статус!")
        return

    tasks = storage.filter_tasks_by_status(status)

    if not tasks:
        print(f"\nЗадач со статусом '{status}' не найдено!")
        return

    print(f"\n=== Задачи со статусом '{status}' ===")
    for task in tasks:
        print(task)


def search_tasks(storage: Storage):
    """Ищет задачи по тексту"""
    print()
    query = input("Введите текст для поиска: ").strip()

    if not query:
        print("Запрос не может быть пустым!")
        return

    tasks = storage.search_tasks(query)

    if not tasks:
        print(f"\nЗадачи, содержащие '{query}', не найдены!")
        return

    print(f"\n=== Результаты поиска для '{query}' ===")
    for task in tasks:
        print(task)


def sort_tasks(storage: Storage):
    """Сортирует задачи"""
    print("\nДоступные варианты сортировки:")
    print("1. id - по ID")
    print("2. created - по дате создания (сначала старые)")
    print("3. updated - по дате обновления (сначала новые)")
    print("4. status - по статусу")

    sort_input = input("\nВыберите вариант сортировки: ").strip()

    sort_map = {
        "1": "id",
        "2": "created",
        "3": "updated",
        "4": "status",
        "id": "id",
        "created": "created",
        "updated": "updated",
        "status": "status"
    }

    sort_by = sort_map.get(sort_input)
    if not sort_by:
        print("Некорректный вариант сортировки!")
        return

    tasks = storage.sort_tasks(sort_by)

    if not tasks:
        print("\nЗадач пока нет!")
        return

    sort_names = {
        "id": "ID",
        "created": "дате создания",
        "updated": "дате обновления",
        "status": "статусу"
    }

    print(f"\n=== Задачи, отсортированные по {sort_names[sort_by]} ===")
    for task in tasks:
        print(task)


def main():
    """Основная функция программы"""
    storage = Storage()
    storage.load()

    print("=== Менеджер Задач ===")

    commands = {
        "list": list_tasks,
        "1": list_tasks,
        "add": add_task,
        "2": add_task,
        "edit": edit_task,
        "3": edit_task,
        "status": change_status,
        "4": change_status,
        "delete": delete_task,
        "5": delete_task,
        "filter": filter_tasks,
        "6": filter_tasks,
        "search": search_tasks,
        "7": search_tasks,
        "sort": sort_tasks,
        "8": sort_tasks,
    }

    while True:
        print("\nКоманды:")
        print("1. list - показать все задачи")
        print("2. add - добавить задачу")
        print("3. edit - редактировать задачу")
        print("4. status - изменить статус задачи")
        print("5. delete - удалить задачу")
        print("6. filter - фильтровать задачи по статусу")
        print("7. search - поиск задач")
        print("8. sort - сортировать задачи")
        print("9. exit - выход")

        command = input("\nВведите команду: ").strip()

        if command in ["exit", "9"]:
            print("До свидания!")
            break

        action = commands.get(command)
        if action:
            action(storage)
        else:
            print("Неизвестная команда!")


if __name__ == "__main__":
    main()
