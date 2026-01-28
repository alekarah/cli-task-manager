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


def update_task(storage: Storage):
    """Универсальная функция обновления задачи"""
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

    print(task)

    print("\nЧто вы хотите обновить?")
    print("1. Название и описание")
    print("2. Статус")
    print("3. Приоритет")
    print("4. Дедлайн")
    print("5. Теги")
    print("6. Всё сразу")

    option = input("\nВыберите опцию: ").strip()

    if option == "1":
        print(f"\nТекущее название: {task.title}")
        title = input("Новое название (Enter - оставить): ").strip()
        print(f"Текущее описание: {task.description}")
        description = input("Новое описание (Enter - оставить): ").strip()
        task.update(title if title else None, description if description else None)

    elif option == "2":
        print(f"\nТекущий статус: {task.status}")
        print("Доступные статусы:")
        print("1. todo\n2. in_progress\n3. done")
        status_input = input("\nВыберите статус: ").strip()
        status_map = {"1": "todo", "2": "in_progress", "3": "done",
                      "todo": "todo", "in_progress": "in_progress", "done": "done"}
        status = status_map.get(status_input)
        if status:
            task.update_status(status)
        else:
            print("Некорректный статус!")
            return

    elif option == "3":
        print(f"\nТекущий приоритет: {task.priority}")
        print("Доступные приоритеты:")
        print("1. low - Низкий 🟢\n2. medium - Средний 🟡\n3. high - Высокий 🔴")
        priority_input = input("\nВыберите приоритет: ").strip()
        priority_map = {"1": "low", "2": "medium", "3": "high",
                        "low": "low", "medium": "medium", "high": "high"}
        priority = priority_map.get(priority_input)
        if priority:
            task.update_priority(priority)
        else:
            print("Некорректный приоритет!")
            return

    elif option == "4":
        # Дедлайн
        from datetime import datetime
        if task.deadline:
            print(f"\nТекущий дедлайн: {task.deadline.strftime('%d.%m.%Y %H:%M')}")
        else:
            print("\nДедлайн не установлен")
        print("Введите новый дедлайн в формате ДД.ММ.ГГГГ ЧЧ:ММ")
        print("Или просто ДД.ММ.ГГГГ (время будет 23:59)")
        print("Или оставьте пустым для удаления дедлайна")
        deadline_input = input("Дедлайн: ").strip()

        if not deadline_input:
            task.update_deadline(None)
        else:
            try:
                # Пытаемся распарсить с временем
                if " " in deadline_input:
                    deadline = datetime.strptime(deadline_input, "%d.%m.%Y %H:%M")
                else:
                    # Если только дата, устанавливаем время 23:59
                    deadline = datetime.strptime(deadline_input, "%d.%m.%Y")
                    deadline = deadline.replace(hour=23, minute=59)
                task.update_deadline(deadline)
            except ValueError:
                print("Некорректный формат даты!")
                return

    elif option == "5":
        # Теги
        if task.tags:
            print(f"\nТекущие теги: {', '.join(['#' + t for t in task.tags])}")
        else:
            print("\nТеги не установлены")
        print("Введите теги через запятую (например: работа, срочно, проект)")
        print("Или оставьте пустым для удаления всех тегов")
        tags_input = input("Теги: ").strip()

        if not tags_input:
            task.set_tags([])
        else:
            tags = [t.strip() for t in tags_input.split(",")]
            task.set_tags(tags)

    elif option == "6":
        # Название и описание
        print(f"\nТекущее название: {task.title}")
        title = input("Новое название (Enter - оставить): ").strip()
        print(f"Текущее описание: {task.description}")
        description = input("Новое описание (Enter - оставить): ").strip()
        task.update(title if title else None, description if description else None)

        # Статус
        print(f"\nТекущий статус: {task.status}")
        print("Доступные статусы: 1. todo  2. in_progress  3. done")
        status_input = input("Выберите статус (Enter - оставить): ").strip()
        if status_input:
            status_map = {"1": "todo", "2": "in_progress", "3": "done",
                          "todo": "todo", "in_progress": "in_progress", "done": "done"}
            status = status_map.get(status_input)
            if status:
                task.update_status(status)

        # Приоритет
        print(f"\nТекущий приоритет: {task.priority}")
        print("Доступные приоритеты: 1. low 🟢  2. medium 🟡  3. high 🔴")
        priority_input = input("Выберите приоритет (Enter - оставить): ").strip()
        if priority_input:
            priority_map = {"1": "low", "2": "medium", "3": "high",
                            "low": "low", "medium": "medium", "high": "high"}
            priority = priority_map.get(priority_input)
            if priority:
                task.update_priority(priority)

        # Дедлайн
        from datetime import datetime
        if task.deadline:
            print(f"\nТекущий дедлайн: {task.deadline.strftime('%d.%m.%Y %H:%M')}")
        else:
            print("\nДедлайн не установлен")
        print("Формат: ДД.ММ.ГГГГ ЧЧ:ММ или ДД.ММ.ГГГГ")
        deadline_input = input("Дедлайн (Enter - оставить): ").strip()
        if deadline_input:
            if deadline_input.lower() == "удалить":
                task.update_deadline(None)
            else:
                try:
                    if " " in deadline_input:
                        deadline = datetime.strptime(deadline_input, "%d.%m.%Y %H:%M")
                    else:
                        deadline = datetime.strptime(deadline_input, "%d.%m.%Y")
                        deadline = deadline.replace(hour=23, minute=59)
                    task.update_deadline(deadline)
                except ValueError:
                    print("Некорректный формат даты! Дедлайн не обновлен.")

        # Теги
        if task.tags:
            print(f"\nТекущие теги: {', '.join(['#' + t for t in task.tags])}")
        else:
            print("\nТеги не установлены")
        print("Введите теги через запятую (Enter - оставить, 'удалить' - очистить)")
        tags_input = input("Теги: ").strip()
        if tags_input:
            if tags_input.lower() == "удалить":
                task.set_tags([])
            else:
                tags = [t.strip() for t in tags_input.split(",")]
                task.set_tags(tags)

    else:
        print("Некорректная опция!")
        return

    storage.save()
    print("\n✓ Задача успешно обновлена!")


def filter_tasks(storage: Storage):
    """Фильтрует задачи по статусу или тегу"""
    print("\nФильтровать по:")
    print("1. Статусу")
    print("2. Тегу")

    filter_type = input("\nВыберите тип фильтрации: ").strip()

    if filter_type == "1":
        # Фильтрация по статусу
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

        tasks = storage.filter_tasks_by_status(status)

        if not tasks:
            print(f"\nЗадач со статусом '{status}' не найдено!")
            return

        print(f"\n=== Задачи со статусом '{status}' ===")
        for task in tasks:
            print(task)

    elif filter_type == "2":
        # Фильтрация по тегу
        all_tags = storage.get_all_tags()
        if all_tags:
            print(f"\nДоступные теги: {', '.join(['#' + t for t in all_tags])}")
        else:
            print("\nТеги не найдены в задачах!")
            return

        tag_input = input("Введите тег для фильтрации: ").strip()
        if not tag_input:
            print("Тег не может быть пустым!")
            return

        # Убираем # если пользователь его ввел
        tag_input = tag_input.lstrip('#')

        tasks = storage.filter_tasks_by_tag(tag_input)

        if not tasks:
            print(f"\nЗадач с тегом '#{tag_input}' не найдено!")
            return

        print(f"\n=== Задачи с тегом '#{tag_input}' ===")
        for task in tasks:
            print(task)

    else:
        print("Некорректная опция!")


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
    print("5. priority - по приоритету (сначала высокий)")

    sort_input = input("\nВыберите вариант сортировки: ").strip()

    sort_map = {
        "1": "id",
        "2": "created",
        "3": "updated",
        "4": "status",
        "5": "priority",
        "id": "id",
        "created": "created",
        "updated": "updated",
        "status": "status",
        "priority": "priority"
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
        "status": "статусу",
        "priority": "приоритету"
    }

    print(f"\n=== Задачи, отсортированные по {sort_names[sort_by]} ===")
    for task in tasks:
        print(task)


def export_tasks(storage: Storage):
    """Экспортирует задачи в файл"""
    print("\nВыберите формат экспорта:")
    print("1. CSV")
    print("2. Markdown")

    format_input = input("\nВыберите формат: ").strip()

    if format_input == "1":
        filename = input("Введите имя файла (Enter - tasks.csv): ").strip()
        if not filename:
            filename = "tasks.csv"
        elif not filename.endswith('.csv'):
            filename += '.csv'

        if storage.export_to_csv(filename):
            print(f"\n✓ Задачи успешно экспортированы в {filename}")
        else:
            print("\nОшибка экспорта!")

    elif format_input == "2":
        filename = input("Введите имя файла (Enter - tasks.md): ").strip()
        if not filename:
            filename = "tasks.md"
        elif not filename.endswith('.md'):
            filename += '.md'

        if storage.export_to_markdown(filename):
            print(f"\n✓ Задачи успешно экспортированы в {filename}")
        else:
            print("\nОшибка экспорта!")

    else:
        print("Некорректный формат!")


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
        "update": update_task,
        "3": update_task,
        "delete": delete_task,
        "4": delete_task,
        "filter": filter_tasks,
        "5": filter_tasks,
        "search": search_tasks,
        "6": search_tasks,
        "sort": sort_tasks,
        "7": sort_tasks,
        "export": export_tasks,
        "8": export_tasks,
    }

    while True:
        print("\nКоманды:")
        print("1. list - показать все задачи")
        print("2. add - добавить задачу")
        print("3. update - обновить задачу (название, статус, приоритет, дедлайн, теги)")
        print("4. delete - удалить задачу")
        print("5. filter - фильтровать задачи по статусу или тегам")
        print("6. search - поиск задач")
        print("7. sort - сортировать задачи")
        print("8. export - экспортировать задачи")
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
