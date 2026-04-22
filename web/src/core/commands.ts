import type { Task, Status, Priority, TerminalLine, OutputType } from '../types'
import * as storage from './storage'

// Глобальный счётчик ID строк — обеспечивает уникальность ключей React
let lineId = 0

// Хелпер для создания строки терминала с нужным типом (цветом)
function line(type: OutputType, content: string): TerminalLine {
  return { id: lineId++, type, content }
}

const STATUS_EMOJI: Record<Status, string> = {
  todo: '📋',
  in_progress: '⚙️',
  done: '✅',
}
const PRIORITY_EMOJI: Record<Priority, string> = {
  low: '🟢',
  medium: '🟡',
  high: '🔴',
}

// Форматирует одну задачу в набор строк терминала с псевдографикой
function formatTask(task: Task): TerminalLine[] {
  const now = new Date()
  const deadlineInfo = task.deadline
    ? (() => {
        const d = new Date(task.deadline)
        const overdue = d < now && task.status !== 'done'
        const str = d.toLocaleDateString('ru-RU', { day: '2-digit', month: '2-digit', year: 'numeric', hour: '2-digit', minute: '2-digit' })
        return overdue ? `  deadline: ${str} ⏰ ПРОСРОЧЕН` : `  deadline: ${str}`
      })()
    : null

  const lines: TerminalLine[] = [
    line('task', `┌─ #${task.id} ${STATUS_EMOJI[task.status]} ${PRIORITY_EMOJI[task.priority]}  ${task.title}`),
    line('output', `│  ${task.description}`),
    line('output', `│  status: ${task.status}  priority: ${task.priority}`),
  ]
  if (deadlineInfo) lines.push(line(task.deadline && new Date(task.deadline) < now && task.status !== 'done' ? 'error' : 'output', `│${deadlineInfo}`))
  if (task.tags.length) lines.push(line('info', `│  tags: ${task.tags.map(t => '#' + t).join(' ')}`))
  lines.push(line('output', `└─ created: ${new Date(task.createdAt).toLocaleString('ru-RU')}  updated: ${new Date(task.updatedAt).toLocaleString('ru-RU')}`))
  return lines
}

// Оборачивает список задач в строки вывода с опциональным заголовком
function taskListOutput(tasks: Task[], header?: string): TerminalLine[] {
  if (!tasks.length) return [line('info', header ? `${header}: пусто` : 'Задач нет')]
  const out: TerminalLine[] = []
  if (header) out.push(line('header', header))
  tasks.forEach(t => out.push(...formatTask(t)))
  return out
}

// Список всех команд — используется для Tab-автодополнения в CommandInput
export const COMMANDS = ['add', 'list', 'delete', 'update', 'filter', 'search', 'sort', 'tags', 'help', 'clear']

// Парсит строку ввода и возвращает массив строк вывода или сигнал 'clear'
export function runCommand(input: string): TerminalLine[] | 'clear' {
  const trimmed = input.trim()
  if (!trimmed) return []

  const [cmd, ...args] = trimmed.split(/\s+/)

  switch (cmd.toLowerCase()) {

    case 'help': return [
      line('header', '┌─ Доступные команды ──────────────────────────────'),
      line('info',   '│  list                       — все задачи'),
      line('info',   '│  add "название" "описание"  — добавить задачу'),
      line('info',   '│  delete <id>                — удалить задачу'),
      line('info',   '│  update <id> <поле> <знач>  — обновить поле'),
      line('info',   '│    поля: status, priority, deadline, tags, title, desc'),
      line('info',   '│    статусы: todo, in_progress, done'),
      line('info',   '│    приоритеты: low, medium, high'),
      line('info',   '│  filter status <статус>     — фильтр по статусу'),
      line('info',   '│  filter tag <тег>           — фильтр по тегу'),
      line('info',   '│  search <текст>             — поиск'),
      line('info',   '│  sort <критерий>            — сортировка'),
      line('info',   '│    критерии: id, created, updated, status, priority'),
      line('info',   '│  tags                       — все теги'),
      line('info',   '│  clear                      — очистить экран'),
      line('header', '└──────────────────────────────────────────────────'),
    ]

    case 'clear': return 'clear'

    case 'list': return taskListOutput(storage.getTasks(), `Задачи (${storage.getTasks().length})`)

    case 'add': {
      const match = trimmed.match(/^add\s+"([^"]+)"\s+"([^"]*)"/)
      if (!match) return [line('error', 'Использование: add "название" "описание"')]
      const task = storage.addTask(match[1], match[2])
      return [line('success', `✓ Задача #${task.id} "${task.title}" создана`)]
    }

    case 'delete': {
      const id = parseInt(args[0])
      if (isNaN(id)) return [line('error', 'Использование: delete <id>')]
      const ok = storage.deleteTask(id)
      return ok
        ? [line('success', `✓ Задача #${id} удалена`)]
        : [line('error', `Задача #${id} не найдена`)]
    }

    case 'update': {
      const id = parseInt(args[0])
      if (isNaN(id) || args.length < 3) return [line('error', 'Использование: update <id> <поле> <значение>')]
      const field = args[1]
      const value = args.slice(2).join(' ')

      const task = storage.getTask(id)
      if (!task) return [line('error', `Задача #${id} не найдена`)]

      const validStatuses: Status[] = ['todo', 'in_progress', 'done']
      const validPriorities: Priority[] = ['low', 'medium', 'high']

      switch (field) {
        case 'status':
          if (!validStatuses.includes(value as Status))
            return [line('error', `Некорректный статус. Допустимые: ${validStatuses.join(', ')}`)]
          storage.updateTask(id, { status: value as Status })
          break
        case 'priority':
          if (!validPriorities.includes(value as Priority))
            return [line('error', `Некорректный приоритет. Допустимые: ${validPriorities.join(', ')}`)]
          storage.updateTask(id, { priority: value as Priority })
          break
        case 'deadline':
          if (value === 'none' || value === 'null') {
            storage.updateTask(id, { deadline: null })
          } else {
            const parts = value.split('.')
            if (parts.length >= 3) {
              const [d, m, y] = parts
              const iso = new Date(`${y}-${m}-${d}`).toISOString()
              storage.updateTask(id, { deadline: iso })
            } else {
              return [line('error', 'Формат дедлайна: ДД.ММ.ГГГГ или none')]
            }
          }
          break
        case 'tags':
          storage.updateTask(id, { tags: value.split(',').map(t => t.trim().toLowerCase()).filter(Boolean) })
          break
        case 'title':
          storage.updateTask(id, { title: value })
          break
        case 'desc':
          storage.updateTask(id, { description: value })
          break
        default:
          return [line('error', `Неизвестное поле "${field}". Допустимые: status, priority, deadline, tags, title, desc`)]
      }
      return [line('success', `✓ Задача #${id} обновлена`)]
    }

    case 'filter': {
      if (args.length < 2) return [line('error', 'Использование: filter status <статус> | filter tag <тег>')]
      const [type, value] = args
      if (type === 'status') {
        const tasks = storage.filterByStatus(value as Status)
        return taskListOutput(tasks, `Статус: ${value} (${tasks.length})`)
      }
      if (type === 'tag') {
        const tasks = storage.filterByTag(value)
        return taskListOutput(tasks, `Тег: #${value} (${tasks.length})`)
      }
      return [line('error', 'Тип фильтра: status или tag')]
    }

    case 'search': {
      if (!args.length) return [line('error', 'Использование: search <текст>')]
      const query = args.join(' ')
      const tasks = storage.searchTasks(query)
      return taskListOutput(tasks, `Поиск: "${query}" (${tasks.length})`)
    }

    case 'sort': {
      const by = args[0] || 'id'
      const valid = ['id', 'created', 'updated', 'status', 'priority']
      if (!valid.includes(by)) return [line('error', `Критерии сортировки: ${valid.join(', ')}`)]
      const tasks = storage.sortTasks(by)
      return taskListOutput(tasks, `Сортировка: ${by} (${tasks.length})`)
    }

    case 'tags': {
      const tags = storage.getAllTags()
      if (!tags.length) return [line('info', 'Тегов нет')]
      return [line('info', `Теги: ${tags.map(t => '#' + t).join('  ')}`)]
    }

    default:
      return [line('error', `Неизвестная команда: "${cmd}". Введите help для списка команд.`)]
  }
}
