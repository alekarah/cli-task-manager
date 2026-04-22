// Статус задачи
export type Status = 'todo' | 'in_progress' | 'done'

// Приоритет задачи
export type Priority = 'low' | 'medium' | 'high'

// Структура задачи — хранится в localStorage как JSON
export interface Task {
  id: number
  title: string
  description: string
  status: Status
  priority: Priority
  deadline: string | null  // ISO-строка или null если дедлайн не задан
  tags: string[]
  createdAt: string        // ISO-строка
  updatedAt: string        // ISO-строка
}

// Тип строки вывода в терминале — определяет цвет
export type OutputType = 'input' | 'output' | 'error' | 'success' | 'info' | 'task' | 'header'

// Одна строка в буфере терминала
export interface TerminalLine {
  id: number
  type: OutputType
  content: string
}
