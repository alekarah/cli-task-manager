import { useState, type KeyboardEvent, useRef, useEffect } from 'react'
import { COMMANDS } from '../core/commands'
import styles from './Terminal.module.css'

interface Props {
  onSubmit: (cmd: string) => void
  history: string[]
}

export default function CommandInput({ onSubmit, history }: Props) {
  const [value, setValue] = useState('')
  // -1 означает «не в режиме просмотра истории» (текущий пустой ввод)
  const [historyIdx, setHistoryIdx] = useState(-1)
  const inputRef = useRef<HTMLInputElement>(null)

  // Автофокус при монтировании компонента
  useEffect(() => {
    inputRef.current?.focus()
  }, [])

  function handleKeyDown(e: KeyboardEvent<HTMLInputElement>) {
    if (e.key === 'Enter') {
      onSubmit(value)
      setValue('')
      setHistoryIdx(-1)
      return
    }

    // Навигация по истории вверх — идём к более старым командам
    if (e.key === 'ArrowUp') {
      e.preventDefault()
      const idx = Math.min(historyIdx + 1, history.length - 1)
      setHistoryIdx(idx)
      setValue(history[history.length - 1 - idx] ?? '')
      return
    }

    // Навигация по истории вниз — возвращаемся к новым командам
    if (e.key === 'ArrowDown') {
      e.preventDefault()
      const idx = Math.max(historyIdx - 1, -1)
      setHistoryIdx(idx)
      setValue(idx === -1 ? '' : history[history.length - 1 - idx] ?? '')
      return
    }

    // Tab-автодополнение по первому совпадению с началом команды
    if (e.key === 'Tab') {
      e.preventDefault()
      const match = COMMANDS.find(c => c.startsWith(value.toLowerCase()))
      if (match) setValue(match)
      return
    }
  }

  return (
    <div className={styles.inputRow}>
      <span className={styles.prompt}>user@task-manager:~$</span>
      <input
        ref={inputRef}
        className={styles.input}
        value={value}
        onChange={e => setValue(e.target.value)}
        onKeyDown={handleKeyDown}
        spellCheck={false}
        autoComplete="off"
      />
    </div>
  )
}
