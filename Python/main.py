#!/usr/bin/env python3
# -*- coding: utf-8 -*-

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
    }

    while True:
        print("\nКоманды:")
        print("1. list - показать все задачи")
        print("2. add - добавить задачу")
        print("3. edit - редактировать задачу")
        print("4. status - изменить статус задачи")
        print("5. delete - удалить задачу")
        print("6. exit - выход")

        command = input("\nВведите команду: ").strip()

        if command in ["exit", "6"]:
            print("До свидания!")
            break

        action = commands.get(command)
        if action:
            action(storage)
        else:
            print("Неизвестная команда!")


if __name__ == "__main__":
    main()
