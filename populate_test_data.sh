#!/bin/bash

# Скрипт для заполнения тестовыми данными обеих версий менеджера задач

echo "=== Заполнение тестовыми данными ==="

# Функция для создания JSON файла с задачами
create_tasks_json() {
    local file_path=$1

    cat > "$file_path" << 'EOF'
{
  "next_id": 11,
  "tasks": [
    {
      "id": 1,
      "title": "Изучить основы языка",
      "description": "Пройти базовый курс по синтаксису и основным конструкциям",
      "status": "done",
      "created_at": "2024-01-15T10:00:00Z",
      "updated_at": "2024-01-20T15:30:00Z"
    },
    {
      "id": 2,
      "title": "Написать REST API",
      "description": "Создать простой REST API для управления задачами",
      "status": "in_progress",
      "created_at": "2024-01-18T09:00:00Z",
      "updated_at": "2024-12-22T10:00:00Z"
    },
    {
      "id": 3,
      "title": "Добавить тесты",
      "description": "Написать unit-тесты для основных функций",
      "status": "todo",
      "created_at": "2024-01-20T14:00:00Z",
      "updated_at": "2024-01-20T14:00:00Z"
    },
    {
      "id": 4,
      "title": "Оптимизировать производительность",
      "description": "Провести профилирование и улучшить узкие места",
      "status": "todo",
      "created_at": "2024-02-01T11:00:00Z",
      "updated_at": "2024-02-01T11:00:00Z"
    },
    {
      "id": 5,
      "title": "Настроить CI/CD",
      "description": "Настроить автоматическую сборку и развертывание",
      "status": "in_progress",
      "created_at": "2024-02-05T08:00:00Z",
      "updated_at": "2024-12-21T16:45:00Z"
    },
    {
      "id": 6,
      "title": "Написать документацию",
      "description": "Создать подробную документацию для API и пользователей",
      "status": "todo",
      "created_at": "2024-02-10T13:00:00Z",
      "updated_at": "2024-02-10T13:00:00Z"
    },
    {
      "id": 7,
      "title": "Рефакторинг кода",
      "description": "Улучшить структуру проекта и читаемость кода",
      "status": "done",
      "created_at": "2024-01-25T10:00:00Z",
      "updated_at": "2024-02-15T12:00:00Z"
    },
    {
      "id": 8,
      "title": "Добавить логирование",
      "description": "Внедрить систему логирования для отслеживания ошибок",
      "status": "in_progress",
      "created_at": "2024-02-12T09:30:00Z",
      "updated_at": "2024-12-22T09:15:00Z"
    },
    {
      "id": 9,
      "title": "Провести code review",
      "description": "Просмотреть код и внести замечания по улучшению",
      "status": "done",
      "created_at": "2024-01-22T14:00:00Z",
      "updated_at": "2024-01-28T16:00:00Z"
    },
    {
      "id": 10,
      "title": "Обновить зависимости",
      "description": "Проверить и обновить все библиотеки до последних версий",
      "status": "todo",
      "created_at": "2024-02-08T15:00:00Z",
      "updated_at": "2024-02-08T15:00:00Z"
    }
  ]
}
EOF

    echo "✓ Создан файл: $file_path"
}

# Создание задач для Go версии
echo ""
echo "Создание тестовых данных для Go версии..."
create_tasks_json "Go/tasks.json"

# Создание задач для Python версии
echo ""
echo "Создание тестовых данных для Python версии..."
create_tasks_json "Python/tasks.json"

echo ""
echo "=== Готово! ==="
echo ""
echo "Тестовые данные созданы для обеих версий."
echo "Теперь вы можете запустить программы и протестировать сортировку:"
echo ""
echo "  Go версия:     cd Go && go run ."
echo "  Python версия: cd Python && python main.py"
echo ""
echo "Доступно 10 задач с разными статусами и датами для тестирования:"
echo "  - 3 задачи со статусом 'done'"
echo "  - 3 задачи со статусом 'in_progress'"
echo "  - 4 задачи со статусом 'todo'"
echo ""
