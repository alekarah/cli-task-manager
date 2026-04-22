import { useState, useEffect, useRef } from 'react'
import type { TerminalLine as TLine } from '../types'
import TerminalLine from './TerminalLine'
import CommandInput from './CommandInput'
import { runCommand } from '../core/commands'
import styles from './Terminal.module.css'

// Начинаем с 1000, чтобы не пересекаться с ID из commands.ts (там от 0)
let nextId = 1000

// Приветственные строки, отображаемые при первом запуске
const BOOT_LINES: TLine[] = [
  { id: nextId++, type: 'header',  content: '┌────────────────────────────────────────────────┐' },
  { id: nextId++, type: 'header',  content: '│            TASK MANAGER  v1.0.0                │' },
  { id: nextId++, type: 'header',  content: '└────────────────────────────────────────────────┘' },
  { id: nextId++, type: 'info',    content: 'Добро пожаловать! Данные хранятся в localStorage.' },
  { id: nextId++, type: 'output',  content: 'Введите help для списка команд.' },
  { id: nextId++, type: 'output',  content: '' },
]

export default function Terminal() {
  const [lines, setLines] = useState<TLine[]>(BOOT_LINES)
  const [history, setHistory] = useState<string[]>([])
  const bottomRef = useRef<HTMLDivElement>(null)
  const containerRef = useRef<HTMLDivElement>(null)

  // Автоскролл вниз при появлении новых строк
  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: 'smooth' })
  }, [lines])

  // Клик по любому месту терминала возвращает фокус на поле ввода
  function focusInput() {
    const input = containerRef.current?.querySelector('input')
    input?.focus()
  }

  // Обрабатывает отправку команды: добавляет эхо ввода и результат в буфер
  function handleSubmit(cmd: string) {
    const trimmed = cmd.trim()
    const echo: TLine = { id: nextId++, type: 'input', content: `user@task-manager:~$ ${trimmed}` }

    if (!trimmed) {
      setLines(prev => [...prev, echo])
      return
    }

    setHistory(prev => [...prev, trimmed])

    const result = runCommand(trimmed)
    if (result === 'clear') {
      setLines([])
      return
    }

    setLines(prev => [...prev, echo, ...result])
  }

  return (
    <div className={styles.terminal} ref={containerRef} onClick={focusInput}>
      <div className={styles.titleBar}>
        <span className={styles.dot} style={{ background: '#ff5f57' }} />
        <span className={styles.dot} style={{ background: '#febc2e' }} />
        <span className={styles.dot} style={{ background: '#28c840' }} />
        <span className={styles.title}>task-manager — bash</span>
      </div>
      <div className={styles.body}>
        {lines.map(l => <TerminalLine key={l.id} line={l} />)}
        <CommandInput onSubmit={handleSubmit} history={history} />
        <div ref={bottomRef} />
      </div>
    </div>
  )
}
