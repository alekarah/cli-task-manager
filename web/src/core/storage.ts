import type { Task, Status, Priority } from '../types'

const STORAGE_KEY = 'task-manager'

interface StorageData {
  tasks: Task[]
  nextId: number
}

// Читает данные из localStorage, при отсутствии или ошибке возвращает пустое состояние
function load(): StorageData {
  const raw = localStorage.getItem(STORAGE_KEY)
  if (!raw) return { tasks: [], nextId: 1 }
  try {
    return JSON.parse(raw)
  } catch {
    return { tasks: [], nextId: 1 }
  }
}

// Сохраняет всё состояние в localStorage
function save(data: StorageData): void {
  localStorage.setItem(STORAGE_KEY, JSON.stringify(data))
}

// Возвращает все задачи
export function getTasks(): Task[] {
  return load().tasks
}

// Создаёт новую задачу со значениями по умолчанию и сохраняет
export function addTask(title: string, description: string): Task {
  const data = load()
  const now = new Date().toISOString()
  const task: Task = {
    id: data.nextId,
    title,
    description,
    status: 'todo',
    priority: 'medium',
    deadline: null,
    tags: [],
    createdAt: now,
    updatedAt: now,
  }
  data.tasks.push(task)
  data.nextId++
  save(data)
  return task
}

// Ищет задачу по ID, возвращает undefined если не найдена
export function getTask(id: number): Task | undefined {
  return load().tasks.find(t => t.id === id)
}

// Удаляет задачу по ID, возвращает false если задача не найдена
export function deleteTask(id: number): boolean {
  const data = load()
  const idx = data.tasks.findIndex(t => t.id === id)
  if (idx === -1) return false
  data.tasks.splice(idx, 1)
  save(data)
  return true
}

// Обновляет указанные поля задачи и автоматически проставляет updatedAt
export function updateTask(id: number, changes: Partial<Omit<Task, 'id' | 'createdAt'>>): Task | null {
  const data = load()
  const task = data.tasks.find(t => t.id === id)
  if (!task) return null
  Object.assign(task, changes, { updatedAt: new Date().toISOString() })
  save(data)
  return task
}

// Фильтрует задачи по статусу
export function filterByStatus(status: Status): Task[] {
  return load().tasks.filter(t => t.status === status)
}

// Фильтрует задачи по тегу (без учёта регистра)
export function filterByTag(tag: string): Task[] {
  const lower = tag.toLowerCase().trim()
  return load().tasks.filter(t => t.tags.includes(lower))
}

// Полнотекстовый поиск по названию и описанию (без учёта регистра)
export function searchTasks(query: string): Task[] {
  const q = query.toLowerCase()
  return load().tasks.filter(
    t => t.title.toLowerCase().includes(q) || t.description.toLowerCase().includes(q)
  )
}

// Возвращает копию задач, отсортированную по указанному критерию
export function sortTasks(by: string): Task[] {
  const tasks = [...load().tasks]
  const priorityOrder: Record<Priority, number> = { high: 1, medium: 2, low: 3 }
  switch (by) {
    case 'id':       return tasks.sort((a, b) => a.id - b.id)
    case 'created':  return tasks.sort((a, b) => a.createdAt.localeCompare(b.createdAt))
    case 'updated':  return tasks.sort((a, b) => b.updatedAt.localeCompare(a.updatedAt))
    case 'status':   return tasks.sort((a, b) => a.status.localeCompare(b.status))
    case 'priority': return tasks.sort((a, b) => priorityOrder[a.priority] - priorityOrder[b.priority])
    default:         return tasks
  }
}

// Возвращает отсортированный список всех уникальных тегов
export function getAllTags(): string[] {
  const tags = new Set<string>()
  load().tasks.forEach(t => t.tags.forEach(tag => tags.add(tag)))
  return [...tags].sort()
}
